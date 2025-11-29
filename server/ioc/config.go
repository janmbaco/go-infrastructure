package ioc

import (
	"github.com/janmbaco/go-infrastructure/v2/configuration/fileconfig/ioc"
	cryptoIoc "github.com/janmbaco/go-infrastructure/v2/crypto/ioc"
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	diskIoc "github.com/janmbaco/go-infrastructure/v2/disk/ioc"
	errorsIoc "github.com/janmbaco/go-infrastructure/v2/errors/ioc"
	eventsIoc "github.com/janmbaco/go-infrastructure/v2/eventsmanager/ioc"
	logsIoc "github.com/janmbaco/go-infrastructure/v2/logs/ioc"
)

// ConfigureServerModules returns all base modules needed for server functionality
func ConfigureServerModules() []dependencyinjection.Module {
	return []dependencyinjection.Module{
		logsIoc.NewLogsModule(),
		errorsIoc.NewErrorsModule(),
		eventsIoc.NewEventsModule(),
		diskIoc.NewDiskModule(),
		ioc.NewConfigurationModule(),
		cryptoIoc.NewCryptoModule(),
		NewServerModule(),
	}
}
