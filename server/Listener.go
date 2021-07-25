package server

import (
	"context"
	"crypto/tls"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"net/http"

	"github.com/janmbaco/go-infrastructure/config"
	"github.com/janmbaco/go-infrastructure/errorhandler"
	"github.com/janmbaco/go-infrastructure/logs"
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
		configureFunc        ConfigureListenerFunc
		setProtobufFunc      SetProtobufFunc
		reStart              bool
		serverSetter         *ServerSetter
		httpServer           *http.Server
		grpcServer           *grpc.Server
		onModifiedConfigFunc func()
	}
)

const (
	HttpServer ServerType = iota
	GRpcSever
)

func NewListener(configHandler config.ConfigHandler, configureFunc ConfigureListenerFunc) *listener {
	listener := &listener{
		configureFunc: configureFunc,
		serverSetter:  &ServerSetter{},
	}
	listener.onModifiedConfigFunc = func() {
		listener.Restart()
	}
	configHandler.OnModifiedConfigSubscriber(listener.onModifiedConfigFunc)
	return listener
}

func (this *listener) SetProtobuf(setProtobufFunc SetProtobufFunc) *listener {
	this.setProtobufFunc = setProtobufFunc
	this.serverSetter.ServerType = GRpcSever
	return this
}

func (this *listener) Start() {
	for (this.httpServer == nil && this.grpcServer == nil) || this.reStart {

		this.reStart = false

		if this.configureFunc == nil {
			panic(errors.New("not configured server"))
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
		case GRpcSever:
			var lis net.Listener
			lis, err = net.Listen("tcp", this.serverSetter.Addr)
			errorhandler.TryPanic(err)
			if this.setProtobufFunc == nil {
				panic(errors.New("ProtoBuffer are not defined, in gRPCServer you must use SetProtobuf function before start listener"))
			}
			this.setProtobufFunc(this.grpcServer)
			err = this.grpcServer.Serve(lis)
		}
		logs.Log.TryWarning(err)
	}
}

func (this *listener) Stop() error {
	logs.Log.Info("Server Stop")
	var err error
	switch this.serverSetter.ServerType {
	case HttpServer:
		err = this.httpServer.Shutdown(context.Background())
	case GRpcSever:
		this.grpcServer.GracefulStop()
	}
	return err
}

func (this *listener) initializeServer() {
	if this.serverSetter.Addr == "" {
		panic(errors.New("address not configured"))
	}
	switch this.serverSetter.ServerType {
	case HttpServer:
		this.httpServer = &http.Server{
			ErrorLog: logs.Log.ErrorLogger,
		}
		this.httpServer.Addr = this.serverSetter.Addr
		if this.serverSetter.Handler != nil {
			this.httpServer.Handler = this.serverSetter.Handler
		}
		this.httpServer.TLSConfig = this.serverSetter.TLSConfig
		if this.serverSetter.TLSNextProto != nil {
			this.httpServer.TLSNextProto = this.serverSetter.TLSNextProto
		}

	case GRpcSever:
		if this.serverSetter.TLSConfig != nil {
			opts := []grpc.ServerOption{grpc.Creds(credentials.NewTLS(this.serverSetter.TLSConfig))}
			this.grpcServer = grpc.NewServer(opts...)
		} else {
			this.grpcServer = grpc.NewServer()
		}
	}
}

func (this *listener) Restart() {
	this.reStart = true
	errorhandler.TryPanic(this.Stop())
}
