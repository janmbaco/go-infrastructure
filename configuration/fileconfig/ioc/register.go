package ioc
import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection/static"
	"github.com/janmbaco/go-infrastructure/configuration"
	"github.com/janmbaco/go-infrastructure/configuration/fileconfig"
)

func init(){
	static.Container.Register().AsSingleton(new(configuration.ConfigHandler), 
		fileconfig.NewFileConfigHandler,
				map[uint]string{
					0: "filePath", 
					1: "defaults",
					6: "fileChangedNotifier",
				})
}