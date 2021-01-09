package events

import (
	"github.com/janmbaco/go-infrastructure/errorhandler"
	"reflect"
	"sync"
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
	pointer := reflect.ValueOf(subscribeFunc).Pointer()
	if _, exists := this.subscriptions[name]; exists {
		if _, isContained := this.subscriptions[name][pointer]; isContained {
			delete(this.subscriptions[name], pointer)
		}
	}
	<-this.isBusy
}

func (this *publisher) Publish(name string) {
	this.isBusy <- true
	wg := sync.WaitGroup{}
	for _, f := range this.subscriptions[name] {
		wg.Add(1)
		go func(fn func()) {
			errorhandler.OnErrorContinue(fn)
			wg.Done()
		}(*f)
	}
	wg.Wait()
	<-this.isBusy
}
