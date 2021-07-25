package events

import (
	"github.com/janmbaco/go-infrastructure/errorhandler"
	"reflect"
)

// Publisher defines an object that publishes events
type Publisher interface {
	Subscribe(name string, subscribeFunc *func())
	UnSubscribe(name string, subscribeFunc *func())
	Publish(name string)
}

type publisher struct {
	subscriptions   map[string]map[uintptr]*func()
	isPublishing    chan bool
	isSubscribing   chan bool
	isUnSubscribing chan bool
}

// NewPublisher returns a Publisher
func NewPublisher() Publisher {
	return &publisher{
		subscriptions:   make(map[string]map[uintptr]*func()),
		isPublishing:    make(chan bool, 1),
		isSubscribing:   make(chan bool, 1),
		isUnSubscribing: make(chan bool, 1),
	}
}

// Subscribe subscribes a function to a event
func (publisher *publisher) Subscribe(name string, subscribeFunc *func()) {
	errorhandler.CheckNilParameter(map[string]interface{}{"subscribeFunc": subscribeFunc})
	publisher.isSubscribing <- true
	errorhandler.TryFinally(func() {
		if _, isContained := publisher.subscriptions[name]; !isContained {
			publisher.subscriptions[name] = make(map[uintptr]*func())
		}
		pointer := reflect.ValueOf(subscribeFunc).Pointer()
		if _, isContained := publisher.subscriptions[name][pointer]; !isContained {
			publisher.subscriptions[name][pointer] = subscribeFunc
		} else {
			panic("This function is already subscribed.")
		}
	}, func() {
		<-publisher.isSubscribing
	})
}

// UnSubscribe unsubscribes a function for a event
func (publisher *publisher) UnSubscribe(name string, subscribeFunc *func()) {
	errorhandler.CheckNilParameter(map[string]interface{}{"subscribeFunc": subscribeFunc})
	publisher.isUnSubscribing <- true
	errorhandler.TryFinally(func() {
		pointer := reflect.ValueOf(subscribeFunc).Pointer()
		if _, exists := publisher.subscriptions[name]; exists {
			if _, isContained := publisher.subscriptions[name][pointer]; isContained {
				delete(publisher.subscriptions[name], pointer)
			} else {
				panic("This function is not subscribed.")
			}
		} else {
			panic("This event is not registered.")
		}
	}, func() {
		<-publisher.isUnSubscribing
	})
}

// Publish publishes a event
func (publisher *publisher) Publish(name string) {
	publisher.isPublishing <- true
	for _, f := range publisher.subscriptions[name] {
		errorhandler.OnErrorContinue(*f)
	}
	<-publisher.isPublishing
}
