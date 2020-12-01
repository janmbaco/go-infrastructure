package event

import "github.com/janmbaco/go-infrastructure/errorhandler"

type EventPublisher struct {
	subscriptions map[string][]func()
	isBusy        chan bool
}

func NewEventPublisher() *EventPublisher {
	return &EventPublisher{
		subscriptions: make(map[string][]func()),
		isBusy:        make(chan bool, 1),
	}
}

func (this *EventPublisher) Subscribe(name string, subscribeFunc func()) {
	this.isBusy <- true
	this.subscriptions[name] = append(this.subscriptions[name], subscribeFunc)
	<-this.isBusy
}

func (this *EventPublisher) Publish(name string) {
	this.isBusy <- true
	subsciptions := this.subscriptions[name]
	<-this.isBusy
	for _, f := range subsciptions {
		errorhandler.OnErrorContinue(f)
	}
}
