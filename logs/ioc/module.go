package ioc

import (
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/v2/logs"
)

// LogsModule implements Module for logging services
type LogsModule struct{}

// NewLogsModule creates a new logs module
func NewLogsModule() *LogsModule {
	return &LogsModule{}
}

// RegisterServices registers all logging services
func (m *LogsModule) RegisterServices(register dependencyinjection.Register) error {
	dependencyinjection.RegisterSingleton[logs.Logger](register, logs.NewLogger)

	// Bind ErrorLogger interface to Logger implementation
	register.Bind(new(logs.ErrorLogger), new(logs.Logger))

	return nil
}
