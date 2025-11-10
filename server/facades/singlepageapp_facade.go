package facades

import (
	"fmt"
	"os"

	"github.com/janmbaco/go-infrastructure/configuration/fileconfig/ioc"
	configResolver "github.com/janmbaco/go-infrastructure/configuration/fileconfig/ioc/resolver"
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	diskIoc "github.com/janmbaco/go-infrastructure/disk/ioc"
	errorsIoc "github.com/janmbaco/go-infrastructure/errors/ioc"
	eventsIoc "github.com/janmbaco/go-infrastructure/eventsmanager/ioc"
	logsIoc "github.com/janmbaco/go-infrastructure/logs/ioc"
	"github.com/janmbaco/go-infrastructure/server"
	serverIoc "github.com/janmbaco/go-infrastructure/server/ioc"
	serverResolver "github.com/janmbaco/go-infrastructure/server/ioc/resolver"
)

func SinglePageAppStart(port, staticPath, index string) {

	// all servers need a configuration.
	// The configuration is monitored to
	// restart the server in case it changes
	type conf struct {
		Port       string `json:"port"`
		StaticPath string `json:"static_path"`
		Index      string `json:"index"`
	}

	// Build container with required modules
	container := dependencyinjection.NewBuilder().
		AddModule(logsIoc.NewLogsModule()).
		AddModule(errorsIoc.NewErrorsModule()).
		AddModule(eventsIoc.NewEventsModule()).
		AddModule(diskIoc.NewDiskModule()).
		AddModule(ioc.NewConfigurationModule()).
		AddModule(serverIoc.NewServerModule()).
		MustBuild()

	resolver := container.Resolver()

	listener, err := serverResolver.GetListenerBuilder(
		resolver,
		configResolver.GetFileConfigHandler(
			resolver,
			os.Args[0]+".json",
			&conf{
				Port:       port,
				StaticPath: staticPath,
				Index:      index,
			},
		),
	).

		// the bootstraper function is performed
		// every time the configuration is modified
		// hence the data is retrieved from the configuration again.
		SetBootstrapper(func(config interface{}, serverSetter *server.ServerSetter) error {
			conf, ok := config.(*conf)
			if !ok {
				return fmt.Errorf("invalid config type")
			}
			serverSetter.Handler = server.NewSinglePageApp(conf.StaticPath, conf.Index)
			serverSetter.Addr = conf.Port
			return nil
		}).
		GetListener()

	if err != nil {
		panic(err)
	}

	<-listener.Start()
}
