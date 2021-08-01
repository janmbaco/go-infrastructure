package events

import (
	"github.com/janmbaco/go-infrastructure/eventsmanager"
)

type RestoredEventHandler struct {
	subscriptions eventsmanager.Subscriptions
}

func NewRestoredEventHandler(subscriptions eventsmanager.Subscriptions) *RestoredEventHandler {
	return &RestoredEventHandler{subscriptions: subscriptions}
}

func (r *RestoredEventHandler) RestoredUnsubscribe(subscription *func()) {
	r.subscriptions.Remove(&RestoredEvent{}, subscription)
}

func (r *RestoredEventHandler) RestoredSubscribe(subscription *func()) {
	r.subscriptions.Add(&RestoredEvent{}, subscription)
}
