package resolver

import (
	
	"github.com/janmbaco/go-infrastructure/dependencyinjection/static"
	_ "github.com/janmbaco/go-infrastructure/logs/ioc"
	_ "github.com/janmbaco/go-infrastructure/errors/ioc"
	_ "github.com/janmbaco/go-infrastructure/eventsmanager/ioc"
	_ "github.com/janmbaco/go-infrastructure/configuration/fileconfig/ioc"

	"github.com/janmbaco/go-infrastructure/configuration"
	"github.com/janmbaco/go-infrastructure/disk"
)

func GetFileConfigHandler(filePath string, defaults interface{}, fileChangedNotifier disk.FileChangedNotifier) configuration.ConfigHandler {
 	return  static.Container.Resolver().Type(
			new(configuration.ConfigHandler),
			map[string]interface{}{
				"filePath": filePath,
				"defaults": defaults,
				"fileChangeNotifier": fileChangedNotifier,
			},
		).(configuration.ConfigHandler)
}