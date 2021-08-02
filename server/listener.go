package server

import (
	"context"
	"crypto/tls"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"net/http"

	"github.com/janmbaco/go-infrastructure/configuration"
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/logs"
)

type (
	Listener interface {
		Start()
		Stop()
	}

	GrpcDefinitionsFunc func(grpcServer *grpc.Server)

	ServerType uint8

	State uint8

	ServerSetter struct {
		Name         string
		Addr         string
		Handler      http.Handler
		TLSConfig    *tls.Config
		TLSNextProto map[string]func(*http.Server, *tls.Conn, http.Handler)
		ServerType   ServerType
	}

	BootstrapperFunc func(config interface{}, serverSetter *ServerSetter)

	listener struct {
		configHandler       configuration.ConfigHandler
		bootstrapperFunc    BootstrapperFunc
		grpcDefinitionsFunc GrpcDefinitionsFunc
		reStart             bool
		serverSetter        *ServerSetter
		httpServer          *http.Server
		grpcServer          *grpc.Server
		logger              logs.Logger
		errorCatcher        errors.ErrorCatcher
		errorDefer          errors.ErrorDefer
		start               chan bool
		stop                chan bool
		finish              chan bool
	}
)

const (
	HttpServer ServerType = iota
	GRpcSever
)

func newListener(configHandler configuration.ConfigHandler, logger logs.Logger, errorCatcher errors.ErrorCatcher, errorThrower errors.ErrorThrower, bootstrapperFunc BootstrapperFunc, grpdDefinitionsFunc GrpcDefinitionsFunc) Listener {
	errorschecker.CheckNilParameter(map[string]interface{}{"configHandler": configHandler, "logger": logger, "errorCatcher": errorCatcher, "errorThrower": errorThrower, "bootstrapperFunc": bootstrapperFunc})
	listener := &listener{
		configHandler:       configHandler,
		logger:              logger,
		errorCatcher:        errorCatcher,
		serverSetter:        &ServerSetter{},
		bootstrapperFunc:    bootstrapperFunc,
		grpcDefinitionsFunc: grpdDefinitionsFunc,
		errorDefer:          errors.NewErrorDefer(errorThrower, &listenerErrorPipe{}),
		start:               make(chan bool, 1),
		stop:                make(chan bool, 1),
		finish:              make(chan bool, 1),
	}

	restoredConfig := listener.onRestoredConfig
	modifiedConfig := listener.onModifiedConfig

	configHandler.RestoredSubscribe(&restoredConfig)
	configHandler.ModifiedSubscribe(&modifiedConfig)
	return listener
}

func (l *listener) Start() {
	defer l.errorDefer.TryThrowError()
	go func() {
		for {
			select {
			case <-l.start:
				l.errorCatcher.TryCatchError(func() {
					l.bootstrapperFunc(l.configHandler.GetConfig(), l.serverSetter)

					l.initializeServer()

					l.logger.Infof("%v - Listen on %v", l.serverSetter.Name, l.serverSetter.Addr)

					switch l.serverSetter.ServerType {
					case HttpServer:
						if l.serverSetter.TLSConfig != nil {
							errorschecker.TryPanic(l.httpServer.ListenAndServeTLS("", ""))
						} else {
							errorschecker.TryPanic(l.httpServer.ListenAndServe())
						}
					case GRpcSever:
						lis, err := net.Listen("tcp", l.serverSetter.Addr)
						errorschecker.TryPanic(err)
						l.grpcDefinitionsFunc(l.grpcServer)
						errorschecker.TryPanic(l.grpcServer.Serve(lis))
					}
				}, func(err error) {
					if err.Error() != "http: Server closed" {
						l.logger.Errorf("%v - %v", l.serverSetter.Name, err.Error())
						if l.configHandler.CanRestore() {
							l.configHandler.Restore()
						} else {
							panic(err)
						}
					}
				})
			case <-l.stop:
				l.finish <- true
				break
			}

		}
	}()
	l.start <- true
	<-l.finish
}

func (l *listener) Stop() {
	defer l.errorDefer.TryThrowError()

	l.logger.Infof("%v - Server Stop", l.serverSetter.Name)
	l.stopServer()
	l.stop <- true

}

func (l *listener) stopServer() {
	switch l.serverSetter.ServerType {
	case HttpServer:
		errorschecker.TryPanic(l.httpServer.Shutdown(context.Background()))
	case GRpcSever:
		l.grpcServer.GracefulStop()
	}
}

func (l *listener) initializeServer() {
	if l.serverSetter.Addr == "" {
		panic(newListenerError(AddressNotConfigured, "address not configured"))
	}
	switch l.serverSetter.ServerType {
	case HttpServer:
		l.httpServer = &http.Server{
			ErrorLog: l.logger.GetErrorLogger(),
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
}

func (l *listener) restart() {
	defer l.errorDefer.TryThrowError()
	l.logger.Tracef("%v Restart Server", l.serverSetter.Name)
	l.reStart = true
	l.stopServer()
	l.start <- true
}

func (l *listener) onRestoredConfig() {
	l.logger.Tracef("%v - Restored config", l.serverSetter.Name)
	l.restart()
}

func (l *listener) onModifiedConfig() {
	l.logger.Tracef("%v - Modified config", l.serverSetter.Name)
	l.restart()
}
