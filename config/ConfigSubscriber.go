package config

import (
	"github.com/janmbaco/go-infrastructure/event"
)

const onModifiedConfigEvent = "onModifiedConfigEvent"

type ConfigSubscriber struct {
	eventPublisher *event.EventPublisher
}

func (this *ConfigSubscriber) OnModifiedConfigSubscriber(subscribeFunc func()) {
	this.eventPublisher.Subscribe(onModifiedConfigEvent, subscribeFunc)
}

func (this *ConfigSubscriber) onModifiedConfigPublish() {
	this.eventPublisher.Publish(onModifiedConfigEvent)
}
