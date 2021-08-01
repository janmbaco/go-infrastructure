package events

import (
	"github.com/janmbaco/go-infrastructure/eventsmanager"
)

type ModificationCanceledEventHandler struct {
	subscriptions eventsmanager.Subscriptions
}

func NewModificationCanceledEventHandler(subscriptions eventsmanager.Subscriptions) *ModificationCanceledEventHandler {
	return &ModificationCanceledEventHandler{subscriptions: subscriptions}
}

func (m *ModificationCanceledEventHandler) ModificationCanceledSubscribe(subscription *func(eventArgs *ModificationCanceledEventArgs)) {
	m.subscriptions.Add(&ModificationCanceledEvent{}, subscription)
}

func (m *ModificationCanceledEventHandler) ModificationCanceledUnsubscribe(subscription *func(eventArgs *ModificationCanceledEventArgs)) {
	m.subscriptions.Remove(&ModificationCanceledEvent{}, subscription)
}
