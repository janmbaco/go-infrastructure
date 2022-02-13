package ioc
import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection/static"
	"github.com/janmbaco/go-infrastructure/server"
)

func init(){
	static.Container.Register().AsSingleton(new(server.ListenerBuilder), server.NewListenerBuilder, map[uint]string{0: "configHandler"})
}