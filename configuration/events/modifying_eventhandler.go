package events

import (
	"github.com/janmbaco/go-infrastructure/eventsmanager"
)

type ModifyingEventHandler struct {
	subscriptions eventsmanager.Subscriptions
}

func NewModifyingEventHandler(subscriptions eventsmanager.Subscriptions) *ModifyingEventHandler {
	return &ModifyingEventHandler{subscriptions: subscriptions}
}

func (m *ModifyingEventHandler) ModifyingSubscribe(subscription *func(eventArgs *ModifyingEventArgs)) {
	m.subscriptions.Add(&ModifyingEvent{}, subscription)
}

func (m *ModifyingEventHandler) ModifyingUnsubscribe(subscription *func(eventArgs *ModifyingEventArgs)) {
	m.subscriptions.Remove(&ModifyingEvent{}, subscription)
}
