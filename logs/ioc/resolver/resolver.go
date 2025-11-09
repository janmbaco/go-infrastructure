package resolver
import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/logs"
)

func GetLogger(resolver dependencyinjection.Resolver) logs.Logger {
	return resolver.Type(new(logs.Logger), nil).(logs.Logger)
}
