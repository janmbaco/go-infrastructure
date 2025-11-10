package ioc

import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/eventsmanager"
)

// EventsModule implements Module for events manager services
type EventsModule struct{}

// NewEventsModule creates a new events module
func NewEventsModule() *EventsModule {
	return &EventsModule{}
}

// RegisterServices registers all events manager services
func (m *EventsModule) RegisterServices(register dependencyinjection.Register) error {
	register.AsSingleton(new(*eventsmanager.EventManager), eventsmanager.NewEventManager, nil)
	return nil
}
