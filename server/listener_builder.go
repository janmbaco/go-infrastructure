package server

import (
	"github.com/janmbaco/go-infrastructure/configuration"
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/logs"
)

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

func NewListenerBuilder(configHandler configuration.ConfigHandler, logger logs.Logger, errorCatcher errors.ErrorCatcher, errorThrower errors.ErrorThrower) ListenerBuilder {
	return &listenerBuilder{configHandler: configHandler, logger: logger, errorCatcher: errorCatcher, errorThrower: errorThrower, errorDefer: errors.NewErrorDefer(errorThrower, &listenerBuilderErrorPipe{})}
}

func (listenerBuilder *listenerBuilder) SetBootstrapper(bootstrapperFunc BootstrapperFunc) ListenerBuilder {
	listenerBuilder.bootstrapperFunc = bootstrapperFunc
	return listenerBuilder
}

func (listenerBuilder *listenerBuilder) SetGrpcDefinitions(grpcDefinitionsFunc GrpcDefinitionsFunc) ListenerBuilder {
	listenerBuilder.grpcDefinitionsFunc = grpcDefinitionsFunc
	return listenerBuilder
}

func (listenerBuilder *listenerBuilder) GetListener() Listener {
	defer listenerBuilder.errorDefer.TryThrowError()
	if listenerBuilder.bootstrapperFunc == nil {
		panic(newListenerBuilderError(NilBootstraperError, "bootsrapper function is not set"))
	}
	serverSetter := &ServerSetter{}
	listenerBuilder.bootstrapperFunc(listenerBuilder.configHandler.GetConfig(), serverSetter)
	if serverSetter.ServerType == GRpcSever && listenerBuilder.grpcDefinitionsFunc == nil {
		panic(newListenerBuilderError(NilGrpcDefinitionsError, "grpc definitions function is not set"))
	}

	return newListener(listenerBuilder.configHandler, listenerBuilder.logger, listenerBuilder.errorCatcher, listenerBuilder.errorThrower, listenerBuilder.bootstrapperFunc, listenerBuilder.grpcDefinitionsFunc)
}
