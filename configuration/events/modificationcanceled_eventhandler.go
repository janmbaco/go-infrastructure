package events

import (
	"github.com/janmbaco/go-infrastructure/eventsmanager"
)

// ModificationCanceledEventHandler handles the subscriptions to the event
type ModificationCanceledEventHandler struct {
	subscriptions eventsmanager.Subscriptions
}

// NewModificationCanceledEventHandler returns a ModificationCanceledEventHandler
func NewModificationCanceledEventHandler(subscriptions eventsmanager.Subscriptions) *ModificationCanceledEventHandler {
	return &ModificationCanceledEventHandler{subscriptions: subscriptions}
}

// ModificationCanceledSubscribe sets new subscription to the event
func (m *ModificationCanceledEventHandler) ModificationCanceledSubscribe(subscription *func(eventArgs *ModificationCanceledEventArgs)) {
	m.subscriptions.Add(&ModificationCanceledEvent{}, subscription)
}

// ModificationCanceledUnsubscribe removes a subscriptions for the event
func (m *ModificationCanceledEventHandler) ModificationCanceledUnsubscribe(subscription *func(eventArgs *ModificationCanceledEventArgs)) {
	m.subscriptions.Remove(&ModificationCanceledEvent{}, subscription)
}
