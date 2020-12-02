package events

import "github.com/janmbaco/go-infrastructure/errorhandler"

type EventPublisher interface {
	Subscribe(name string, subscribeFunc func())
	Publish(name string)
}

type eventPublisher struct {
	subscriptions map[string][]func()
	isBusy        chan bool
}

func NewEventPublisher() EventPublisher {
	return &eventPublisher{
		subscriptions: make(map[string][]func()),
		isBusy:        make(chan bool, 1),
	}
}

func (this *eventPublisher) Subscribe(name string, subscribeFunc func()) {
	this.isBusy <- true
	this.subscriptions[name] = append(this.subscriptions[name], subscribeFunc)
	<-this.isBusy
}

func (this *eventPublisher) Publish(name string) {
	this.isBusy <- true
	subsciptions := this.subscriptions[name]
	<-this.isBusy
	for _, f := range subsciptions {
		errorhandler.OnErrorContinue(f)
	}
}
