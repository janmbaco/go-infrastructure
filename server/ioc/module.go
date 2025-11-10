package ioc

import (
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/v2/server"
)

// ServerModule implements Module for server services
type ServerModule struct{}

// NewServerModule creates a new server module
func NewServerModule() *ServerModule {
	return &ServerModule{}
}

// RegisterServices registers all server services
func (m *ServerModule) RegisterServices(register dependencyinjection.Register) error {
	dependencyinjection.RegisterSingletonWithParams[server.ListenerBuilder](
		register,
		server.NewListenerBuilder,
		map[int]string{0: "configHandler"},
	)

	return nil
}
