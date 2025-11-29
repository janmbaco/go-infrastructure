package facades

import (
	"fmt"
	"os"

	configResolver "github.com/janmbaco/go-infrastructure/v2/configuration/fileconfig/ioc/resolver"
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/v2/server"
	serverIoc "github.com/janmbaco/go-infrastructure/v2/server/ioc"
	serverResolver "github.com/janmbaco/go-infrastructure/v2/server/ioc/resolver"
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
		AddModules(serverIoc.ConfigureServerModules()...).
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
