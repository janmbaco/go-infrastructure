package resolver

import (
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/v2/logs"
)

func GetLogger(resolver dependencyinjection.Resolver) logs.Logger {
	result := resolver.Type(new(logs.Logger), nil)
	if logger, ok := result.(logs.Logger); ok {
		return logger
	}
	panic("failed to resolve Logger")
}
