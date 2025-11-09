package server

import (
	"reflect"

	"github.com/janmbaco/go-infrastructure/configuration"
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/logs"
)

// ListenerBuilder defines a object responsible to builds listeners
type ListenerBuilder interface {
	SetBootstrapper(bootstrapperFunc BootstrapperFunc) ListenerBuilder
	SetGrpcDefinitions(setProtobufFunc GrpcDefinitionsFunc) ListenerBuilder
	SetConfigValidatorFunc(configValidationFunc ConfigValidatorFunc) ListenerBuilder
	SetConfigApplicatorFunc(configApplicatorFunc ConfigApplicatorFunc) ListenerBuilder
	GetListener() (Listener, error)
}

type listenerBuilder struct {
	configHandler       configuration.ConfigHandler
	logger              logs.Logger
	errorCatcher        errors.ErrorCatcher
	bootstrapperFunc    BootstrapperFunc
	grpcDefinitionsFunc GrpcDefinitionsFunc
	validationFunc      ConfigValidatorFunc
	applicationFunc     ConfigApplicatorFunc
}

// NewListenerBuilder returns a ListenerBuilder
func NewListenerBuilder(configHandler configuration.ConfigHandler, logger logs.Logger, errorCatcher errors.ErrorCatcher) ListenerBuilder {
	return &listenerBuilder{configHandler: configHandler, logger: logger, errorCatcher: errorCatcher}
}

// SetBootstrapper sets the function that serves like bootstraper to the listener start
func (lb *listenerBuilder) SetBootstrapper(bootstrapperFunc BootstrapperFunc) ListenerBuilder {
	lb.bootstrapperFunc = bootstrapperFunc
	return lb
}

// SetGrpcDefinitions sets the definitions of grpc function
func (lb *listenerBuilder) SetGrpcDefinitions(grpcDefinitionsFunc GrpcDefinitionsFunc) ListenerBuilder {
	lb.grpcDefinitionsFunc = grpcDefinitionsFunc
	return lb
}

// SetConfigValidatorFunc sets the validation function for config changes
func (lb *listenerBuilder) SetConfigValidatorFunc(validationFunc ConfigValidatorFunc) ListenerBuilder {
	lb.validationFunc = validationFunc
	return lb
}

// SetConfigApplicatorFunc sets the application function for config changes
func (lb *listenerBuilder) SetConfigApplicatorFunc(applicationFunc ConfigApplicatorFunc) ListenerBuilder {
	lb.applicationFunc = applicationFunc
	return lb
}

// GetListener gets the listener
func (lb *listenerBuilder) GetListener() (Listener, error) {
	if lb.bootstrapperFunc == nil {
		return nil, lb.pipError(newListenerBuilderError(NilBootstraperError, "bootsrapper function is not set", nil))
	}
	serverSetter := &ServerSetter{}

	// Mute logs during validation to avoid duplicate output
	lb.logger.Mute()
	if err := lb.bootstrapperFunc(lb.configHandler.GetConfig(), serverSetter); err != nil {
		lb.logger.Unmute()
		return nil, lb.pipError(err)
	}
	lb.logger.Unmute()

	if serverSetter.ServerType == GRpcSever && lb.grpcDefinitionsFunc == nil {
		return nil, lb.pipError(newListenerBuilderError(NilGrpcDefinitionsError, "grpc definitions function is not set", nil))
	}
	listener := newListener(lb.configHandler, lb.logger, lb.errorCatcher, lb.bootstrapperFunc, lb.grpcDefinitionsFunc, lb.validationFunc, lb.applicationFunc)
	lb.bootstrapperFunc = nil
	lb.grpcDefinitionsFunc = nil
	lb.validationFunc = nil
	lb.applicationFunc = nil
	return listener, nil
}

func (lb *listenerBuilder) pipError(err error) error {
	resultError := err
	if errType := reflect.Indirect(reflect.ValueOf(err)).Type(); !errType.Implements(reflect.TypeOf((*ListenerBuilderError)(nil)).Elem()) {
		resultError = newListenerBuilderError(UnexpectedBuilderError, err.Error(), err)
	}
	return resultError
}
