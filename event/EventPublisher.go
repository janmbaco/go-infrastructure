package event

import "github.com/janmbaco/go-infrastructure/errorhandler"

type EventPublisher struct {
	subscriptions map[string][]func()
}

func NewEventPublisher() *EventPublisher {
	return &EventPublisher{subscriptions: make(map[string][]func())}
}

func (this *EventPublisher) Subscribe(name string, subscribeFunc func()) {
	this.subscriptions[name] = append(this.subscriptions[name], subscribeFunc)
}

func (this *EventPublisher) Publish(name string) {
	for _, f := range this.subscriptions[name] {
		errorhandler.OnErrorContinue(f)
	}
}
