package events

import (
	"github.com/janmbaco/go-infrastructure/errorhandler"
	"reflect"
)

type Publisher interface {
	Subscribe(name string, subscribeFunc *func())
	UnSubscribe(name string, subscribeFunc *func())
	Publish(name string)
}

type publisher struct {
	subscriptions map[string]map[uintptr]*func()
	isBusy        chan bool
}

func NewPublisher() Publisher {
	return &publisher{
		subscriptions: make(map[string]map[uintptr]*func()),
		isBusy:        make(chan bool, 1),
	}
}

func (this *publisher) Subscribe(name string, subscribeFunc *func()) {
	errorhandler.CheckNilParameter(map[string]interface{}{"subscribeFunc": subscribeFunc})
	this.isBusy <- true
	if _, isContained := this.subscriptions[name]; !isContained {
		this.subscriptions[name] = make(map[uintptr]*func())
	}
	pointer := reflect.ValueOf(subscribeFunc).Pointer()
	if _, isContained := this.subscriptions[name][pointer]; !isContained {
		this.subscriptions[name][pointer] = subscribeFunc
	}
	<-this.isBusy
}

func (this *publisher) UnSubscribe(name string, subscribeFunc *func()) {
	errorhandler.CheckNilParameter(map[string]interface{}{"subscribeFunc": subscribeFunc})
	this.isBusy <- true
	subscriptions := make(map[uintptr]*func())
	for pointer, s := range this.subscriptions[name] {
		if pointer != reflect.ValueOf(subscribeFunc).Pointer() {
			subscriptions[pointer] = s
		}
	}
	this.subscriptions[name] = subscriptions
	<-this.isBusy
}

func (this *publisher) Publish(name string) {
	this.isBusy <- true
	subsciptions := this.subscriptions[name]
	<-this.isBusy
	for _, f := range subsciptions {
		errorhandler.OnErrorContinue(*f)
	}
}
