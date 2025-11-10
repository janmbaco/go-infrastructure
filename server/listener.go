package server

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"reflect"
	"time"

	"github.com/janmbaco/go-infrastructure/configuration"
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
		serverSetter         *ServerSetter
		httpServer           *http.Server
		grpcServer           *grpc.Server
		bootstrapperFunc     BootstrapperFunc
		grpcDefinitionsFunc  GrpcDefinitionsFunc
		configValidatorFunc  ConfigValidatorFunc
		configApplicatorFunc ConfigApplicatorFunc
		start                chan bool
		started              chan bool
		stop                 chan bool
		finish               chan ListenerError
		isBusy               chan bool
		stopped              bool
	}
)

const (
	HTTPServer ServerType = iota
	GRpcSever
)

func newListener(configHandler configuration.ConfigHandler, logger logs.Logger, errorCatcher errors.ErrorCatcher, bootstrapperFunc BootstrapperFunc, grpdDefinitionsFunc GrpcDefinitionsFunc, validationFunc ConfigValidatorFunc, applicationFunc ConfigApplicatorFunc) Listener {
	listener := &listener{
		configHandler:        configHandler,
		logger:               logger,
		errorCatcher:         errorCatcher,
		serverSetter:         &ServerSetter{},
		bootstrapperFunc:     bootstrapperFunc,
		grpcDefinitionsFunc:  grpdDefinitionsFunc,
		configValidatorFunc:  validationFunc,
		configApplicatorFunc: applicationFunc,
		start:                make(chan bool, 1),
		started:              make(chan bool, 1),
		stop:                 make(chan bool, 1),
		finish:               make(chan ListenerError, 1),
		isBusy:               make(chan bool, 1),
	}
	restoredConfig := listener.onRestoredConfig
	modifiedConfig := listener.onModifiedConfig

	configHandler.RestoredSubscribe(&restoredConfig)
	configHandler.ModifiedSubscribe(&modifiedConfig)
	return listener
}

func (l *listener) Start() chan ListenerError {
	l.stopped = false
	go l.startLoop()
	l.start <- true
	<-l.started
	return l.finish
}

func (l *listener) Stop() {
	l.isBusy <- true
	l.logger.Infof("%v - Server Stop", l.serverSetter.Name)
	_ = l.stopServer() //nolint:errcheck // server shutdown errors are not recoverable
	l.stop <- true
	<-l.isBusy
}

func (l *listener) startLoop() {
	defer func() {
		if r := recover(); r != nil {
			l.handleRecover(r, true)
		}
	}()
	for {
		select {
		case <-l.start:
			err := l.errorCatcher.TryCatchError(func() error {
				if err := l.bootstrapperFunc(l.configHandler.GetConfig(), l.serverSetter); err != nil {
					return err
				}
				if err := l.initializeServer(); err != nil {
					return err
				}
				l.logger.Infof("%v - Listen on %v", l.serverSetter.Name, l.serverSetter.Addr)

				switch l.serverSetter.ServerType {
				case HTTPServer:
					select {
					case l.started <- true:
					default:
					}
					if l.serverSetter.TLSConfig != nil {
						if err := l.httpServer.ListenAndServeTLS("", ""); err != nil {
							return err
						}
					} else {
						if err := l.httpServer.ListenAndServe(); err != nil {
							return err
						}
					}
				case GRpcSever:
					select {
					case l.started <- true:
					default:
					}
					lis, err := net.Listen("tcp", l.serverSetter.Addr)
					if err != nil {
						return err
					}
					l.grpcDefinitionsFunc(l.grpcServer)
					if err := l.grpcServer.Serve(lis); err != nil {
						return err
					}
				}
				return nil
			}, func(err error) {
				l.handleServerError(err)
			})
			if err != nil {
				l.handleServerError(err)
			}
		case <-l.stop:
			l.finish <- nil
			l.stopped = true
		}
		if l.stopped {
			break
		}
	}
}
func (l *listener) initializeServer() error {
	if l.serverSetter.Addr == "" {
		return newListenerError(AddressNotConfigured, "address not configured", nil)
	}
	switch l.serverSetter.ServerType {
	case HTTPServer:
		l.httpServer = &http.Server{
			ErrorLog:          l.logger.GetErrorLogger(),
			ReadHeaderTimeout: 10 * time.Second,
		}
		l.httpServer.Addr = l.serverSetter.Addr
		if l.serverSetter.Handler != nil {
			l.httpServer.Handler = l.serverSetter.Handler
		}
		l.httpServer.TLSConfig = l.serverSetter.TLSConfig
		if l.serverSetter.TLSNextProto != nil {
			l.httpServer.TLSNextProto = l.serverSetter.TLSNextProto
		}

	case GRpcSever:
		if l.serverSetter.TLSConfig != nil {
			opts := []grpc.ServerOption{grpc.Creds(credentials.NewTLS(l.serverSetter.TLSConfig))}
			l.grpcServer = grpc.NewServer(opts...)
		} else {
			l.grpcServer = grpc.NewServer()
		}
	}
	return nil
}

func (l *listener) stopServer() error {
	switch l.serverSetter.ServerType {
	case HTTPServer:
		if l.httpServer != nil {
			if err := l.httpServer.Shutdown(context.Background()); err != nil {
				return err
			}
		}
	case GRpcSever:
		if l.grpcServer != nil {
			l.grpcServer.GracefulStop()
		}
	}
	return nil
}

func (l *listener) restart() {
	l.isBusy <- true
	if !l.stopped {
		l.logger.Tracef("%v Restart Server", l.serverSetter.Name)
		_ = l.stopServer() //nolint:errcheck // server shutdown errors are not recoverable
		l.start <- true
	}
	<-l.isBusy
}

func (l *listener) onRestoredConfig() {
	l.logger.Tracef("%v - Restored config", l.serverSetter.Name)
	l.restart()
}

func (l *listener) onModifiedConfig() {
	l.logger.Tracef("%v - Modified config", l.serverSetter.Name)

	// Call validation function
	if l.configValidatorFunc != nil {
		valid, err := l.configValidatorFunc(l.configHandler.GetConfig())
		if err != nil {
			l.logger.Errorf("%v - Config validation error: %v", l.serverSetter.Name, err)
			return
		}
		if !valid {
			l.logger.Tracef("%v - Config validation cancelled", l.serverSetter.Name)
			return
		}
	}

	// Call application function
	needsRestart := false
	if l.configApplicatorFunc != nil {
		application := ConfigApplication{NeedsRestart: &needsRestart}
		err := l.configApplicatorFunc(l.configHandler.GetConfig(), &application)
		if err != nil {
			l.logger.Errorf("%v - Config application error: %v", l.serverSetter.Name, err)
			return
		}
	}

	// If application function is not set or needs restart, restart
	if l.configApplicatorFunc == nil || needsRestart {
		l.restart()
	} else {
		l.logger.Infof("%v - Config applied without restart", l.serverSetter.Name)
	}
}

func (l *listener) handleRecover(r interface{}, sendError bool) {
	var err error
	switch v := r.(type) {
	case error:
		err = v
	case string:
		err = newListenerError(UnexpectedError, v, nil)
	default:
		err = newListenerError(UnexpectedError, "panic in listener", nil)
	}
	l.logger.Errorf("%v - panic recovered: %v", l.serverSetter.Name, err)
	l.finalizeError(err, sendError)
}

func (l *listener) handleServerError(err error) {
	if err.Error() != "http: Server closed" {
		l.logger.Errorf("%v - %v", l.serverSetter.Name, err.Error())
		l.finalizeError(err, true)
	}
}

func (l *listener) finalizeError(err error, sendError bool) {
	if l.configHandler.CanRestore() {
		if restoreErr := l.configHandler.Restore(); restoreErr != nil {
			l.logger.Errorf("%v - Failed to restore config: %v", l.serverSetter.Name, restoreErr)
		}
	} else if sendError {
		if listenerErr, ok := l.pipeError(err).(ListenerError); ok {
			l.finish <- listenerErr
		} else {
			l.logger.Errorf("%v - Failed to convert error to ListenerError: %v", l.serverSetter.Name, err)
		}
		l.stopped = true
	}
}

func (l *listener) pipeError(err error) error {
	resultError := err

	if errType := reflect.Indirect(reflect.ValueOf(err)).Type(); !errType.Implements(reflect.TypeOf((*ListenerError)(nil)).Elem()) {
		resultError = newListenerError(UnexpectedError, err.Error(), err)
	}

	return resultError
}
