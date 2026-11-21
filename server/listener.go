package server

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"

	"github.com/janmbaco/go-infrastructure/v2/configuration"
	"github.com/janmbaco/go-infrastructure/v2/errors"
	"github.com/janmbaco/go-infrastructure/v2/logs"
	"google.golang.org/grpc"
)

type (
	Listener interface {
		Start() chan ListenerError
		Stop()
	}

	GrpcDefinitionsFunc func(grpcServer *grpc.Server)

	ServerType uint8 //nolint:revive // established API name, stuttering is acceptable

	State uint8

	ServerSetter struct { //nolint:revive // established API name, stuttering is acceptable
		// Pointers and interfaces first for memory alignment
		Handler      http.Handler
		TLSConfig    *tls.Config
		TLSNextProto map[string]func(*http.Server, *tls.Conn, http.Handler)

		// Strings
		Name string
		Addr string

		// Primitives
		ServerType ServerType
	}

	BootstrapperFunc func(config interface{}, serverSetter *ServerSetter) error

	ConfigValidatorFunc func(config interface{}) (bool, error)

	ConfigApplication struct {
		NeedsRestart *bool
	}

	ConfigApplicatorFunc func(config interface{}, configApplication *ConfigApplication) error

	listener struct {
		configHandler        configuration.ConfigHandler
		logger               logs.Logger
		errorCatcher         errors.ErrorCatcher
		bootstrapperFunc     BootstrapperFunc
		grpcDefinitionsFunc  GrpcDefinitionsFunc
		configValidatorFunc  ConfigValidatorFunc
		configApplicatorFunc ConfigApplicatorFunc
		mu                   sync.Mutex
		run                  *listenerRun
		// Updated only by the serving coordinator after a successful bind.
		setter ServerSetter
	}

	listenerRun struct {
		ready   chan struct{}
		changed chan struct{}
		stop    chan struct{}
		finish  chan ListenerError
	}
)

const (
	HTTPServer ServerType = iota
	GRpcSever
)

func newListener(configHandler configuration.ConfigHandler, logger logs.Logger, errorCatcher errors.ErrorCatcher, bootstrapperFunc BootstrapperFunc, grpcDefinitionsFunc GrpcDefinitionsFunc, validationFunc ConfigValidatorFunc, applicationFunc ConfigApplicatorFunc) Listener {
	l := &listener{
		configHandler:        configHandler,
		logger:               logger,
		errorCatcher:         errorCatcher,
		bootstrapperFunc:     bootstrapperFunc,
		grpcDefinitionsFunc:  grpcDefinitionsFunc,
		configValidatorFunc:  validationFunc,
		configApplicatorFunc: applicationFunc,
	}
	changed := l.onModifiedConfig
	configHandler.ModifiedSubscribe(&changed)
	return l
}

func (l *listener) Start() chan ListenerError {
	l.mu.Lock()
	r := l.run
	if r == nil {
		r = &listenerRun{
			ready:   make(chan struct{}),
			changed: make(chan struct{}, 1),
			stop:    make(chan struct{}),
			finish:  make(chan ListenerError, 1),
		}
		l.run = r
		go l.serve(r)
	}
	l.mu.Unlock()
	<-r.ready
	return r.finish
}

func (l *listener) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if r := l.run; r != nil {
		select {
		case <-r.stop:
		default:
			close(r.stop)
		}
	}
}

// Notification can run inside the request which will be drained by Shutdown.
// The sole coordinator reads the latest snapshot after taking this wake-up.
func (l *listener) onModifiedConfig() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if r := l.run; r != nil {
		select {
		case r.changed <- struct{}{}:
		default:
		}
	}
}

func (l *listener) serve(r *listenerRun) {
	var active *servingGeneration
	var result error
	var ready sync.Once
	notifyReady := func() { ready.Do(func() { close(r.ready) }) }
	defer func() {
		if recovered := recover(); recovered != nil {
			result = fmt.Errorf("panic in listener: %v", recovered)
		}
		if active != nil {
			active.drain()
		}
		l.mu.Lock()
		l.run = nil
		var terminal ListenerError
		if result != nil {
			if err, ok := result.(ListenerError); ok {
				terminal = err
			} else {
				terminal = newListenerError(UnexpectedError, result.Error(), result)
			}
		}
		r.finish <- terminal
		close(r.finish)
		notifyReady()
		l.mu.Unlock()
	}()

	initial := true
	var config listenerConfig
	var lastVersion uint64
	for {
		select {
		case <-r.stop:
			return
		default:
		}
		if !initial {
			var served <-chan error
			if active != nil {
				served = active.done
			}
			select {
			case <-r.stop:
				return
			case err := <-served:
				active = nil
				if err == nil {
					return
				}
				if result = l.recoverConfig(config, err); result != nil {
					return
				}
				continue
			case <-r.changed:
			}
		}
		next, err := readListenerConfig(l.configHandler)
		if err != nil {
			result = err
			return
		}
		if !initial && next.versioned && next.version == lastVersion {
			continue
		}
		lastVersion = next.version
		if !initial && !l.shouldBind(next.value, active != nil) {
			continue
		}
		if active != nil {
			active.drain()
			active = nil
		}
		select {
		case <-r.stop:
			return
		default:
		}
		config = next
		err = l.errorCatcher.TryCatchError(func() error {
			var openErr error
			active, openErr = l.open(config.value)
			return openErr
		}, nil)
		if err != nil {
			select {
			case <-r.stop:
				return
			default:
			}
			if result = l.recoverConfig(config, err); result != nil {
				return
			}
		}
		initial = false
		notifyReady()
	}
}

func (l *listener) shouldBind(config any, serving bool) bool {
	if l.configValidatorFunc != nil {
		valid, err := l.configValidatorFunc(config)
		if err != nil {
			l.logger.Errorf("Config validation error: %v", err)
			return false
		}
		if !valid {
			return false
		}
	}
	if l.configApplicatorFunc == nil {
		return true
	}
	restart := false
	if err := l.configApplicatorFunc(config, &ConfigApplication{NeedsRestart: &restart}); err != nil {
		l.logger.Errorf("Config application error: %v", err)
		return false
	}
	// Hot apply can keep an active server. A recovered run may have none.
	return restart || !serving
}

// A successful rollback keeps the run available for a new change or Stop.
func (l *listener) recoverConfig(config listenerConfig, failure error) error {
	l.logger.Errorf("Listener error: %v", failure)
	if versioned, ok := l.configHandler.(versionedListenerConfig); ok {
		// Compare the revision and restore under one lock. Another listener may
		// already have restored, or a writer may have published a newer value.
		if _, err := versioned.RestoreIfVersion(config.version); err != nil {
			l.logger.Errorf("Could not restore listener configuration: %v", err)
			return failure
		}
		return nil
	}
	if !l.configHandler.CanRestore() {
		return failure
	}
	return l.configHandler.Restore()
}

// The built-in file provider supplies a stable typed snapshot and can condition
// rollback on its revision. Other ConfigHandler implementations remain valid.
type versionedListenerConfig interface {
	ConfigSnapshot() (any, uint64, error)
	RestoreIfVersion(version uint64) (bool, error)
}

type listenerConfig struct {
	value     any
	version   uint64
	versioned bool
}

func readListenerConfig(handler configuration.ConfigHandler) (listenerConfig, error) {
	if versioned, ok := handler.(versionedListenerConfig); ok {
		value, version, err := versioned.ConfigSnapshot()
		return listenerConfig{value: value, version: version, versioned: true}, err
	}
	return listenerConfig{value: handler.GetConfig()}, nil
}
