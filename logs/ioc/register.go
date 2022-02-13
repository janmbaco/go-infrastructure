package ioc
import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection/static"
	"github.com/janmbaco/go-infrastructure/logs"
)

func init(){
	static.Container.Register().AsSingleton(new(logs.Logger), logs.NewLogger, nil)
	static.Container.Register().Bind(new(logs.ErrorLogger), new(logs.Logger))
}