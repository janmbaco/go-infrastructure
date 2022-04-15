package resolver

import (
	
	"github.com/janmbaco/go-infrastructure/dependencyinjection/static"
	_ "github.com/janmbaco/go-infrastructure/eventsmanager/ioc"
	_ "github.com/janmbaco/go-infrastructure/disk/ioc"

	"github.com/janmbaco/go-infrastructure/disk"
)

func GetFileChangedNotifier(filePath string) disk.FileChangedNotifier {
 	return  static.Container.Resolver().Type(new(disk.FileChangedNotifier), map[string]interface{}{ "filePath": filePath }).(disk.FileChangedNotifier)
}
