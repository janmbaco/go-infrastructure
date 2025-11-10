package eventsmanager

import (
	"reflect"
)

// EventManager manages multiple event publishers
type EventManager struct {
	publishers map[reflect.Type]interface{}
}

// NewEventManager creates a new EventManager
func NewEventManager() *EventManager {
	return &EventManager{
		publishers: make(map[reflect.Type]interface{}),
	}
}

// Register registers a publisher for a specific event type
func Register[T EventObject[T]](em *EventManager, publisher Publisher[T]) {
	typ := reflect.TypeOf(*new(T))
	em.publishers[typ] = publisher
}

// Publish publishes an event using the registered publisher
func Publish[T EventObject[T]](em *EventManager, event T) {
	typ := reflect.TypeOf(event)
	if publisher, exists := em.publishers[typ]; exists {
		if pub, ok := publisher.(Publisher[T]); ok {
			pub.Publish(event)
		}
	}
}
