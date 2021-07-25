package config

import (
	"github.com/janmbaco/go-infrastructure/errorhandler"
	"github.com/janmbaco/go-infrastructure/events"
)

const modifiedEvent = "modifiedEvent"
const modifyingEvent = "modifyingEvent"

// ConfigSubscriber defines an object that is capable of subscribing to changes to a configuration
type ConfigSubscriber interface {
	OnModifiedConfigSubscriber(subscribeFunc func())
	OnModifyingConfigSubscriber(subscribeFunc func(newConfig interface{}))
}

type configSubscriber struct {
	eventPublisher events.Publisher
	cancel         bool
	cancelMessage  string
	newConfig      interface{}
}

// OnModifyingConfigSubscriber occurs when the configuration is being modifyng externally
func (configSubscriber *configSubscriber) OnModifyingConfigSubscriber(subscribeFunc func(newConfig interface{})) {
	onModifying := func() {
		if !configSubscriber.cancel {
			errorhandler.TryCatchError(func() {
				subscribeFunc(configSubscriber.newConfig)
			}, func(err error) {
				configSubscriber.cancelMessage = err.Error()
				configSubscriber.cancel = true
			})
		}
	}
	configSubscriber.eventPublisher.Subscribe(modifyingEvent, &onModifying)
}

// OnModifiedConfigSubscriber occurs when the configuration is modified externally
func (configSubscriber *configSubscriber) OnModifiedConfigSubscriber(subscribeFunc func()) {
	onModified := func() {
		subscribeFunc()
	}
	configSubscriber.eventPublisher.Subscribe(modifiedEvent, &onModified)
}

func (configSubscriber *configSubscriber) onModifiedConfigPublish() {
	configSubscriber.eventPublisher.Publish(modifiedEvent)
}

func (configSubscriber *configSubscriber) onModifyingConfigPublish(newConfig interface{}) bool {
	configSubscriber.cancel = false
	configSubscriber.newConfig = newConfig
	configSubscriber.eventPublisher.Publish(modifyingEvent)
	return configSubscriber.cancel
}
