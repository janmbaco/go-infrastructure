package events

import (
	"github.com/janmbaco/go-infrastructure/eventsmanager"
)

// ModifyingEventHandler handles the subscriptions to the event
type ModifyingEventHandler struct {
	subscriptions eventsmanager.Subscriptions
}

// NewModifyingEventHandler returns a ModifyingEventHandler
func NewModifyingEventHandler(subscriptions eventsmanager.Subscriptions) *ModifyingEventHandler {
	return &ModifyingEventHandler{subscriptions: subscriptions}
}

// ModifyingSubscribe sets new subscription to the event
func (m *ModifyingEventHandler) ModifyingSubscribe(subscription *func(eventArgs *ModifyingEventArgs)) {
	m.subscriptions.Add(&ModifyingEvent{}, subscription)
}

// ModifyingUnsubscribe removes a subscriptions for the event
func (m *ModifyingEventHandler) ModifyingUnsubscribe(subscription *func(eventArgs *ModifyingEventArgs)) {
	m.subscriptions.Remove(&ModifyingEvent{}, subscription)
}
