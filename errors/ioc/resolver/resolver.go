package resolver

import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/errors"
)

func GetErrorCatcher(resolver dependencyinjection.Resolver) errors.ErrorCatcher {
	return resolver.Type(new(errors.ErrorCatcher), nil).(errors.ErrorCatcher)
}