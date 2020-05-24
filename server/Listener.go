package server

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/janmbaco/go-infrastructure/config"
	"github.com/janmbaco/go-infrastructure/errorhandler"
	"github.com/janmbaco/go-infrastructure/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"net/http"
)

type (
	SetProtobufFunc func(grpcServer *grpc.Server)

	ServerType uint8

	ServerSetter struct {
		Addr         string
		Handler      http.Handler
		TLSConfig    *tls.Config
		TLSNextProto map[string]func(*http.Server, *tls.Conn, http.Handler)
		ServerType   ServerType
	}

	ConfigureListenerFunc func(serverSetter *ServerSetter)

	listener struct {
		configureFunc   ConfigureListenerFunc
		setProtobufFunc SetProtobufFunc
		reStart         bool
		serverSetter    *ServerSetter
		httpServer      *http.Server
		grpcServer      *grpc.Server
	}
)

const (
	HttpServer ServerType = iota
	gRpcSever
)

func NewListener(configureFunc ConfigureListenerFunc, onConfigChangeSubscriber config.OnConfigChangeSubscriber) *listener {
	listener := &listener{
		configureFunc: configureFunc,
		serverSetter:  &ServerSetter{},
	}
	onConfigChangeSubscriber.Subscribe(func(config *config.ConfigBase) {
		listener.Restart()
	})
	return listener
}

func (this *listener) SetProtobuf(setProtobufFunc SetProtobufFunc) *listener {
	this.setProtobufFunc = setProtobufFunc
	this.serverSetter.ServerType = gRpcSever
	return this
}

func (this *listener) Start() {
	for (this.httpServer == nil && this.grpcServer == nil) || this.reStart {

		this.reStart = false

		if this.configureFunc == nil {
			errorhandler.TryPanic(errors.New("not configured server"))
		}
		this.configureFunc(this.serverSetter)

		this.initializeServer()

		logs.Log.Info("Listen on " + this.serverSetter.Addr)
		var err error
		switch this.serverSetter.ServerType {
		case HttpServer:
			if this.serverSetter.TLSConfig != nil {
				err = this.httpServer.ListenAndServeTLS("", "")
			} else {
				err = this.httpServer.ListenAndServe()
			}
		case gRpcSever:
			var lis net.Listener
			lis, err = net.Listen("tcp", this.serverSetter.Addr)
			errorhandler.TryPanic(err)
			if this.setProtobufFunc != nil {
				this.setProtobufFunc(this.grpcServer)
			}
			err = this.grpcServer.Serve(lis)
		}
		logs.Log.Warning(err.Error())
	}
}

func (this *listener) Stop() error {
	logs.Log.Info("Server Stop")
	var err error
	switch this.serverSetter.ServerType {
	case HttpServer:
		err = this.httpServer.Shutdown(context.Background())
	case gRpcSever:
		this.grpcServer.GracefulStop()
	}
	return err
}

func (this *listener) initializeServer() {
	if this.serverSetter.Addr == "" {
		errorhandler.TryPanic(errors.New("address not configured"))
	}
	switch this.serverSetter.ServerType {
	case HttpServer:
		if this.serverSetter.Handler == nil {
			errorhandler.TryPanic(errors.New("handler routes not configured"))
		}
		this.httpServer = &http.Server{
			ErrorLog: logs.Log.ErrorLogger,
		}
		this.httpServer.Addr = this.serverSetter.Addr
		this.httpServer.Handler = this.serverSetter.Handler
		this.httpServer.TLSConfig = this.serverSetter.TLSConfig
		if this.serverSetter.TLSNextProto != nil {
			this.httpServer.TLSNextProto = this.serverSetter.TLSNextProto
		}

	case gRpcSever:
		opts := []grpc.ServerOption{grpc.Creds(credentials.NewTLS(this.serverSetter.TLSConfig))}
		this.grpcServer = grpc.NewServer(opts...)
	}
}

func (this *listener) Restart() {
	this.reStart = true
	errorhandler.TryPanic(this.Stop())
}
