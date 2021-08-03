package events

import (
	"github.com/janmbaco/go-infrastructure/eventsmanager"
)

// ModifiedEventHandler handles the subscriptions to the event
type ModifiedEventHandler struct {
	subscriptions eventsmanager.Subscriptions
}

// NewModifiedEventHandler returns a ModifiedEventHandler
func NewModifiedEventHandler(subscriptions eventsmanager.Subscriptions) *ModifiedEventHandler {
	return &ModifiedEventHandler{subscriptions: subscriptions}
}

// ModifiedUnsubscribe sets new subscription to the event
func (m *ModifiedEventHandler) ModifiedUnsubscribe(subscription *func()) {
	m.subscriptions.Remove(&ModifiedEvent{}, subscription)
}

// ModifiedSubscribe removes a subscriptions for the event
func (m *ModifiedEventHandler) ModifiedSubscribe(subscription *func()) {
	m.subscriptions.Add(&ModifiedEvent{}, subscription)
}
