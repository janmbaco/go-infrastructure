package ioc

import (
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/v2/errors"
)

// ErrorsModule implements Module for error handling services
type ErrorsModule struct{}

// NewErrorsModule creates a new errors module
func NewErrorsModule() *ErrorsModule {
	return &ErrorsModule{}
}

// RegisterServices registers all error handling services
func (m *ErrorsModule) RegisterServices(register dependencyinjection.Register) error {
	// Register core services
	register.AsSingleton(new(errors.ErrorCatcher), errors.NewErrorCatcher, nil)

	return nil
}
