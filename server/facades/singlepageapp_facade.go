package facades

import (

	"github.com/janmbaco/go-infrastructure/server"

	fileConfigResolver "github.com/janmbaco/go-infrastructure/configuration/fileconfig/ioc/resolver"
	diskResolver "github.com/janmbaco/go-infrastructure/disk/ioc/resolver"
	serverResolver "github.com/janmbaco/go-infrastructure/server/ioc/resolver"
)

func SinglePageAppStart(port string, staticPath string, index string) {

	// all servers need a configuration.
	// The configuration is monitored to
	// restart the server in case it changes
	type conf struct {
		Port       string `json:"port"`
		StaticPath string `json:"static_path"`
		Index      string `json:"index"`
	}

	serverResolver.GetListenerBuilder(
			fileConfigResolver.GetFileConfigHandler(
				"config.json",  
				&conf{
					Port:       port,
					StaticPath: staticPath,
					Index:      index,
				},
				diskResolver.GetFileChangedNotifier("config.json"),
			),
		).

		// the bootstraper function is performed
		//every time the configuration is modified
		//hence the data is retrieved from the configuration again.
		SetBootstrapper(func(config interface{}, serverSetter *server.ServerSetter) {
			serverSetter.Handler = server.NewSinglePageApp(config.(*conf).StaticPath,config.(*conf).Index) 
			serverSetter.Addr = config.(*conf).Port
		}).
		GetListener().
		Start()
}
