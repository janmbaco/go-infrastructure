package dependencyinjection_test

import (
	"github.com/janmbaco/go-infrastructure/configuration"
	"github.com/janmbaco/go-infrastructure/configuration/fileconfig"
	"github.com/janmbaco/go-infrastructure/crypto"
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/logs"
	"github.com/janmbaco/go-infrastructure/persistence"
	"github.com/janmbaco/go-infrastructure/server"
)

func Registerfacade(container *dependencyinjection.Container) {
	container.Register.AsSingleton(new(logs.Logger), logs.NewLogger, nil)
	container.Register.Bind(new(logs.ErrorLogger), new(logs.Logger))
	container.Register.AsSingleton(new(errors.ErrorCatcher), errors.NewErrorCatcher, nil)
	container.Register.AsSingleton(new(errors.ErrorManager), errors.NewErrorManager, nil)
	container.Register.Bind(new(errors.ErrorCallbacks), new(errors.ErrorManager))
	container.Register.AsSingleton(new(errors.ErrorThrower), errors.NewErrorThrower, nil)
	container.Register.AsSingleton(new(configuration.ConfigHandler), fileconfig.NewFileConfigHandler, map[uint]string{0: "filePath", 1: "defaults"})
	container.Register.AsSingleton(new(server.ListenerBuilder), server.NewListenerBuilder, nil)
	container.Register.AsType(new(persistence.DataAccess), persistence.NewDataAccess, map[uint]string{0: "db", 1: "modelType"})
	container.Register.AsSingleton(new(crypto.Cipher), crypto.NewCipher, map[uint]string{0: "key"})
}
