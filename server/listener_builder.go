package server

import (
	"github.com/janmbaco/go-infrastructure/configuration"
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"github.com/janmbaco/go-infrastructure/logs"
)

// ListenerBuilder defines a object responsible to builds listeners
type ListenerBuilder interface {
	SetBootstrapper(bootstrapperFunc BootstrapperFunc) ListenerBuilder
	SetGrpcDefinitions(setProtobufFunc GrpcDefinitionsFunc) ListenerBuilder
	GetListener() Listener
}

type listenerBuilder struct {
	configHandler       configuration.ConfigHandler
	logger              logs.Logger
	errorCatcher        errors.ErrorCatcher
	errorThrower        errors.ErrorThrower
	errorDefer          errors.ErrorDefer
	bootstrapperFunc    BootstrapperFunc
	grpcDefinitionsFunc GrpcDefinitionsFunc
	isGrpcServer        bool
}

// NewListenerBuilder returns a ListenerBuilder
func NewListenerBuilder(configHandler configuration.ConfigHandler, logger logs.Logger, errorCatcher errors.ErrorCatcher, errorThrower errors.ErrorThrower) ListenerBuilder {
	errorschecker.CheckNilParameter(map[string]interface{}{"configHandler": configHandler, "logger": logger, "errorCatcher": errorCatcher, "errorThrower": errorThrower})
	return &listenerBuilder{configHandler: configHandler, logger: logger, errorCatcher: errorCatcher, errorThrower: errorThrower, errorDefer: errors.NewErrorDefer(errorThrower, &listenerBuilderErrorPipe{})}
}

// SetBootstrapper sets the function that serves like bootstraper to the listener start
func (lb *listenerBuilder) SetBootstrapper(bootstrapperFunc BootstrapperFunc) ListenerBuilder {
	errorschecker.CheckNilParameter(map[string]interface{}{"bootstrapperFunc": bootstrapperFunc})
	lb.bootstrapperFunc = bootstrapperFunc
	return lb
}

// SetGrpcDefinitions sets the definitions of grpc function
func (lb *listenerBuilder) SetGrpcDefinitions(grpcDefinitionsFunc GrpcDefinitionsFunc) ListenerBuilder {
	errorschecker.CheckNilParameter(map[string]interface{}{"grpcDefinitionsFunc": grpcDefinitionsFunc})
	lb.grpcDefinitionsFunc = grpcDefinitionsFunc
	return lb
}

// GetListener gets the listener
func (lb *listenerBuilder) GetListener() Listener {
	defer lb.errorDefer.TryThrowError()
	if lb.bootstrapperFunc == nil {
		panic(newListenerBuilderError(NilBootstraperError, "bootsrapper function is not set"))
	}
	serverSetter := &ServerSetter{}
	lb.bootstrapperFunc(lb.configHandler.GetConfig(), serverSetter)
	if serverSetter.ServerType == GRpcSever && lb.grpcDefinitionsFunc == nil {
		panic(newListenerBuilderError(NilGrpcDefinitionsError, "grpc definitions function is not set"))
	}
	listener := newListener(lb.configHandler, lb.logger, lb.errorCatcher, lb.errorThrower, lb.bootstrapperFunc, lb.grpcDefinitionsFunc)
	lb.bootstrapperFunc = nil
	lb.grpcDefinitionsFunc = nil
	return listener
}
