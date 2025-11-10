package resolver

import (
	"github.com/janmbaco/go-infrastructure/configuration"
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/server"
)

func GetListenerBuilder(resolver dependencyinjection.Resolver, configHandler configuration.ConfigHandler) server.ListenerBuilder {
	result, ok := resolver.Type(new(server.ListenerBuilder),
		map[string]interface{}{
			"configHandler": configHandler,
		}).(server.ListenerBuilder)
	if !ok {
		panic("failed to resolve ListenerBuilder")
	}
	return result
}
