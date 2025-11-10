package resolver

import (
	"github.com/janmbaco/go-infrastructure/v2/configuration"
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
)

func GetFileConfigHandler(resolver dependencyinjection.Resolver, filePath string, defaults interface{}) configuration.ConfigHandler {
	result := resolver.Type(
		new(configuration.ConfigHandler),
		map[string]interface{}{
			"filePath": filePath,
			"defaults": defaults,
		},
	)
	if configHandler, ok := result.(configuration.ConfigHandler); ok {
		return configHandler
	}
	panic("failed to resolve ConfigHandler")
}
