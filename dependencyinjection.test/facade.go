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

func Registerfacade(register dependencyinjection.Register) {
	register.AsSingleton(new(logs.Logger), logs.NewLogger, nil)
	register.Bind(new(logs.ErrorLogger), new(logs.Logger))
	register.AsSingleton(new(errors.ErrorCatcher), errors.NewErrorCatcher, nil)
	register.AsSingleton(new(errors.ErrorManager), errors.NewErrorManager, nil)
	register.Bind(new(errors.ErrorCallbacks), new(errors.ErrorManager))
	register.AsSingleton(new(errors.ErrorThrower), errors.NewErrorThrower, nil)
	register.AsSingleton(new(configuration.ConfigHandler), fileconfig.NewFileConfigHandler, map[uint]string{0: "filePath", 1: "defaults"})
	register.AsSingleton(new(server.ListenerBuilder), server.NewListenerBuilder, nil)
	register.AsType(new(persistence.DataAccess), persistence.NewDataAccess, map[uint]string{0: "db", 1: "modelType"})
	register.AsSingleton(new(crypto.Cipher), crypto.NewCipher, map[uint]string{0: "key"})
}
