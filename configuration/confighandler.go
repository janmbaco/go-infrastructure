package configuration

type (
	// RestoredSubscriber defines an object responsible to subscribe function for the event Restored
	RestoredSubscriber interface {
		RestoredSubscribe(*func())
		RestoredUnsubscribe(*func())
	}
	// ModifiedSubscriber defines an object responsible to subscribe function for the event Modified
	ModifiedSubscriber interface {
		ModifiedSubscribe(*func())
		ModifiedUnsubscribe(*func())
	}
	// ConfigHandler defines a object that handles the configuration
	ConfigHandler interface {
		ModifiedSubscriber
		RestoredSubscriber
		GetConfig() interface{}
		SetConfig(interface{}) error
		Freeze()
		Unfreeze()
		CanRestore() bool
		Restore() error
		SetRefreshTime(Period) error
		ForceRefresh() error
	}
)
