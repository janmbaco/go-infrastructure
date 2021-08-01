package events

import (
	"github.com/janmbaco/go-infrastructure/eventsmanager"
)

type ModifiedEventHandler struct {
	subscriptions eventsmanager.Subscriptions
}

func NewModifiedEventHandler(subscriptions eventsmanager.Subscriptions) *ModifiedEventHandler {
	return &ModifiedEventHandler{subscriptions: subscriptions}
}

func (m *ModifiedEventHandler) ModifiedUnsubscribe(subscription *func()) {
	m.subscriptions.Remove(&ModifiedEvent{}, subscription)
}

func (m *ModifiedEventHandler) ModifiedSubscribe(subscription *func()) {
	m.subscriptions.Add(&ModifiedEvent{}, subscription)
}
