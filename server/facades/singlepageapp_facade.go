package facades

import (
	"github.com/janmbaco/go-infrastructure/configuration"
	"github.com/janmbaco/go-infrastructure/configuration/fileconfig"
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/logs"
	"github.com/janmbaco/go-infrastructure/server"
	"net/http"
)

func SinglePageAppStart(port string, staticPath string, index string) {
	container := dependencyinjection.NewContainer()

	container.Register().AsSingleton(new(logs.Logger), logs.NewLogger, nil)
	container.Register().Bind(new(logs.ErrorLogger), new(logs.Logger))
	container.Register().AsSingleton(new(errors.ErrorCatcher), errors.NewErrorCatcher, nil)
	container.Register().AsSingleton(new(errors.ErrorManager), errors.NewErrorManager, nil)
	container.Register().Bind(new(errors.ErrorCallbacks), new(errors.ErrorManager))
	container.Register().AsSingleton(new(errors.ErrorThrower), errors.NewErrorThrower, nil)
	container.Register().AsSingleton(new(configuration.ConfigHandler), fileconfig.NewFileConfigHandler, map[uint]string{0: "filePath", 1: "defaults"})
	container.Register().AsSingleton(new(server.ListenerBuilder), server.NewListenerBuilder, nil)
	container.Register().AsSingleton(new(http.Handler), server.NewSinglePageApp, map[uint]string{0: "staticPath", 1: "index"})

	// all servers need a configuration.
	// The configuration is monitored to
	// restart the server in case it changes
	type conf struct {
		Port       string `json:"port"`
		StaticPath string `json:"static_path"`
		Index      string `json:"index"`
	}
	<-container.Resolver().Type(
		new(server.ListenerBuilder),
		map[string]interface{}{
			"filePath": "config.json",
			"defaults": &conf{
				Port:       port,
				StaticPath: staticPath,
				Index:      index,
			},
		}).(server.ListenerBuilder).

		// the bootstraper function is performed
		//every time the configuration is modified
		//hence the data is retrieved from the configuration again.
		SetBootstrapper(func(config interface{}, serverSetter *server.ServerSetter) {
			serverSetter.Handler = container.Resolver().Type(
				new(http.Handler),
				map[string]interface{}{
					"staticPath": config.(*conf).StaticPath,
					"index":      config.(*conf).Index,
				},
			).(http.Handler)
			serverSetter.Addr = config.(*conf).Port
		},
		).
		GetListener().
		Start()
}
