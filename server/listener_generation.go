package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// The coordinator owns the generation until both graceful drain and Serve return.
type servingGeneration struct {
	done     <-chan error
	shutdown func()
}

func (g *servingGeneration) drain() {
	g.shutdown()
	<-g.done
}

func (l *listener) open(config any) (_ *servingGeneration, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("panic initializing listener: %v", recovered)
		}
	}()
	setter := l.setter
	if err := l.bootstrapperFunc(config, &setter); err != nil {
		return nil, err
	}
	if setter.Addr == "" {
		return nil, newListenerError(AddressNotConfigured, "address not configured", nil)
	}
	var serve func(net.Listener) error
	var shutdown func()
	switch setter.ServerType {
	case HTTPServer:
		s := &http.Server{
			Addr:              setter.Addr,
			Handler:           setter.Handler,
			TLSConfig:         setter.TLSConfig,
			TLSNextProto:      setter.TLSNextProto,
			ErrorLog:          l.logger.GetErrorLogger(),
			ReadHeaderTimeout: 10 * time.Second,
		}
		serve = s.Serve
		if setter.TLSConfig != nil {
			serve = func(ln net.Listener) error { return s.ServeTLS(ln, "", "") }
		}
		shutdown = func() {
			if err := s.Shutdown(context.Background()); err != nil {
				l.logger.Errorf("%s - Shutdown: %v", setter.Name, err)
			}
		}
	case GRpcSever:
		if l.grpcDefinitionsFunc == nil {
			return nil, fmt.Errorf("gRPC definitions are not configured")
		}
		var opts []grpc.ServerOption
		if setter.TLSConfig != nil {
			opts = append(opts, grpc.Creds(credentials.NewTLS(setter.TLSConfig)))
		}
		s := grpc.NewServer(opts...)
		l.grpcDefinitionsFunc(s)
		serve = s.Serve
		shutdown = s.GracefulStop
	default:
		return nil, fmt.Errorf("unsupported server type %d", setter.ServerType)
	}
	l.logger.Infof("%v - Bind on %v", setter.Name, setter.Addr)
	ln, err := net.Listen("tcp", setter.Addr)
	if err != nil {
		return nil, err
	}
	l.setter = setter
	done := make(chan error, 1)
	go func() {
		var serveErr error
		defer func() {
			if recovered := recover(); recovered != nil {
				serveErr = fmt.Errorf("panic serving %s: %v", setter.Name, recovered)
			}
			// ServeTLS can reject the certificate before it takes ownership.
			_ = ln.Close()
			done <- serveErr
		}()
		serveErr = serve(ln)
	}()
	return &servingGeneration{done: done, shutdown: shutdown}, nil
}
