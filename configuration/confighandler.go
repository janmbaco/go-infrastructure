package configuration

import "github.com/janmbaco/go-infrastructure/configuration/events"

type (
	// ModificationCanceledSubscriber defines an object responsible to subscribe function for the event ModificationCanceled
	ModificationCanceledSubscriber interface {
		ModificationCanceledSubscribe(*func(eventArgs *events.ModificationCanceledEventArgs))
		ModificationCanceledUnsubscribe(*func(eventArgs *events.ModificationCanceledEventArgs))
	}
	// RestoredSubscriber defines an object responsible to subscribe function for the event Restored
	RestoredSubscriber interface {
		RestoredSubscribe(*func())
		RestoredUnsubscribe(*func())
	}
	// ModifyingSubscriber defines an object responsible to subscribe function for the event Modifying
	ModifyingSubscriber interface {
		ModifyingSubscribe(*func(eventArgs *events.ModifyingEventArgs))
		ModifyingUnsubscribe(*func(eventArgs *events.ModifyingEventArgs))
	}
	// ModifiedSubscriber defines an object responsible to subscribe function for the event Modified
	ModifiedSubscriber interface {
		ModifiedSubscribe(*func())
		ModifiedUnsubscribe(*func())
	}
	// ConfigHandler defines a object that handles the configuration
	ConfigHandler interface {
		ModifiedSubscriber
		ModifyingSubscriber
		ModificationCanceledSubscriber
		RestoredSubscriber
		GetConfig() interface{}
		Freeze()
		Unfreeze()
		CanRestore() bool
		Restore()
		SetRefreshTime(Period)
		ForceRefresh()
	}
)
