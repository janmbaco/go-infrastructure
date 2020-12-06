package config

import (
	"github.com/janmbaco/go-infrastructure/events"
)

const onModifiedConfigEvent = "onModifiedConfigEvent"

type ConfigSubscriber struct {
	eventPublisher events.Publisher
}

func (this *ConfigSubscriber) OnModifiedConfigSubscriber(subscribeFunc func()) {
	this.eventPublisher.Subscribe(onModifiedConfigEvent, subscribeFunc)
}

func (this *ConfigSubscriber) onModifiedConfigPublish() {
	this.eventPublisher.Publish(onModifiedConfigEvent)
}
