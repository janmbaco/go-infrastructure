package resolver
import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection/static"
	"github.com/janmbaco/go-infrastructure/errors"

	_ "github.com/janmbaco/go-infrastructure/logs/ioc"
	_ "github.com/janmbaco/go-infrastructure/errors/ioc"
)

func GetErrorCatcher() errors.ErrorCatcher {
 	return static.Container.Resolver().Type(new(errors.ErrorCatcher), nil).(errors.ErrorCatcher)
}