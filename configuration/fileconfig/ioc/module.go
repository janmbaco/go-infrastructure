package ioc

import (
	"github.com/janmbaco/go-infrastructure/v2/configuration"
	"github.com/janmbaco/go-infrastructure/v2/configuration/fileconfig"
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
)

// ConfigurationModule implements Module for configuration services
type ConfigurationModule struct{}

// NewConfigurationModule creates a new configuration module
func NewConfigurationModule() *ConfigurationModule {
	return &ConfigurationModule{}
}

// RegisterServices registers all configuration services
func (m *ConfigurationModule) RegisterServices(register dependencyinjection.Register) error {
	dependencyinjection.RegisterTypeWithParams[configuration.ConfigHandler](
		register,
		fileconfig.NewFileConfigHandler,
		map[int]string{0: "filePath", 1: "defaults", 4: "fileChangedNotifier"},
	)

	return nil
}
