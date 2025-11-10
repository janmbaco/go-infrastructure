package events

import (
	"github.com/janmbaco/go-infrastructure/eventsmanager"
)

// ModifiedEventHandler handles the subscriptions to the event
type ModifiedEventHandler struct {
	subscriptions eventsmanager.Subscriptions[ModifiedEvent]
}

// NewModifiedEventHandler returns a ModifiedEventHandler
func NewModifiedEventHandler(subscriptions eventsmanager.Subscriptions[ModifiedEvent]) *ModifiedEventHandler {
	return &ModifiedEventHandler{subscriptions: subscriptions}
}

// ModifiedUnsubscribe sets new subscription to the event
func (m *ModifiedEventHandler) ModifiedUnsubscribe(subscription *func()) {
	if subscription != nil {
		fn := func(ModifiedEvent) { (*subscription)() }
		_ = m.subscriptions.Remove(fn) //nolint:errcheck // event subscription errors are not actionable
	}
}

// ModifiedSubscribe removes a subscriptions for the event
func (m *ModifiedEventHandler) ModifiedSubscribe(subscription *func()) {
	if subscription != nil {
		fn := func(ModifiedEvent) { (*subscription)() }
		_ = m.subscriptions.Add(fn) //nolint:errcheck // event subscription errors are not actionable
	}
}
