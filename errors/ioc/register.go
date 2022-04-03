package ioc
import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection/static"
	"github.com/janmbaco/go-infrastructure/errors"
)

func init(){
	static.Container.Register().AsSingleton(new(errors.ErrorCatcher), errors.NewErrorCatcher, nil)
	static.Container.Register().AsSingleton(new(errors.ErrorManager), errors.NewErrorManager, nil)
	static.Container.Register().Bind(new(errors.ErrorCallbacks), new(errors.ErrorManager))
	static.Container.Register().AsSingleton(new(errors.ErrorThrower), errors.NewErrorThrower, nil)
}