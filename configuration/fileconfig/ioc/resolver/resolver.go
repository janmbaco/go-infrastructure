package resolver
import (
	"github.com/janmbaco/go-infrastructure/configuration"
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
)

func GetFileConfigHandler(resolver dependencyinjection.Resolver, filePath string, defaults interface{}) configuration.ConfigHandler {
	return resolver.Type(
		new(configuration.ConfigHandler),
		map[string]interface{}{
			"filePath": filePath,
			"defaults": defaults,
		},
	).(configuration.ConfigHandler)
}
