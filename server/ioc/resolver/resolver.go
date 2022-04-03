package resolver

import (
	"github.com/janmbaco/go-infrastructure/configuration"
	"github.com/janmbaco/go-infrastructure/server"

	"github.com/janmbaco/go-infrastructure/dependencyinjection/static"
	
	_ "github.com/janmbaco/go-infrastructure/logs/ioc"
	_ "github.com/janmbaco/go-infrastructure/errors/ioc"
	_ "github.com/janmbaco/go-infrastructure/server/ioc"
)

func GetListenerBuilder(configHandler configuration.ConfigHandler) server.ListenerBuilder {
 	return static.Container.Resolver().Type(new(server.ListenerBuilder),
	  map[string]interface{}{
			"configHandler": configHandler,
		}).(server.ListenerBuilder)
}

