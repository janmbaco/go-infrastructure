package events

import (
	"github.com/janmbaco/go-infrastructure/eventsmanager"
)

// RestoredEventHandler handles the subscriptions to the event
type RestoredEventHandler struct {
	subscriptions eventsmanager.Subscriptions
}

// NewRestoredEventHandler returns a  RestoredEventHandler
func NewRestoredEventHandler(subscriptions eventsmanager.Subscriptions) *RestoredEventHandler {
	return &RestoredEventHandler{subscriptions: subscriptions}
}

// RestoredUnsubscribe sets new subscription to the event
func (r *RestoredEventHandler) RestoredUnsubscribe(subscription *func()) {
	r.subscriptions.Remove(&RestoredEvent{}, subscription)
}

// RestoredSubscribe removes a subscriptions for the event
func (r *RestoredEventHandler) RestoredSubscribe(subscription *func()) {
	r.subscriptions.Add(&RestoredEvent{}, subscription)
}
