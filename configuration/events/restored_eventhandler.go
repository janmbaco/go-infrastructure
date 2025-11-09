package events

import (
	"github.com/janmbaco/go-infrastructure/eventsmanager"
)

// RestoredEventHandler handles the subscriptions to the event
type RestoredEventHandler struct {
	subscriptions eventsmanager.Subscriptions[RestoredEvent]
}

// NewRestoredEventHandler returns a  RestoredEventHandler
func NewRestoredEventHandler(subscriptions eventsmanager.Subscriptions[RestoredEvent]) *RestoredEventHandler {
	return &RestoredEventHandler{subscriptions: subscriptions}
}

// RestoredUnsubscribe sets new subscription to the event
func (r *RestoredEventHandler) RestoredUnsubscribe(subscription *func()) {
	if subscription != nil {
		fn := func(RestoredEvent) { (*subscription)() }
		r.subscriptions.Remove(fn)
	}
}

// RestoredSubscribe removes a subscriptions for the event
func (r *RestoredEventHandler) RestoredSubscribe(subscription *func()) {
	if subscription != nil {
		fn := func(RestoredEvent) { (*subscription)() }
		r.subscriptions.Add(fn)
	}
}
