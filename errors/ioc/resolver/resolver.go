package resolver

import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/errors"
)

func GetErrorCatcher(resolver dependencyinjection.Resolver) errors.ErrorCatcher {
	result := resolver.Type(new(errors.ErrorCatcher), nil)
	if errorCatcher, ok := result.(errors.ErrorCatcher); ok {
		return errorCatcher
	}
	panic("failed to resolve ErrorCatcher")
}
