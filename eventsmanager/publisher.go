package eventsmanager

import (
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"reflect"
	"sync"
)

// Publisher defines an object responsible to publish events
type Publisher interface {
	Publish(event EventObject)
}

type (
	publisher struct {
		subscriptions   Subscriptions
		eventPublishers sync.Map
		errorCatcher    errors.ErrorCatcher
	}
	eventPublisher struct {
		isPublishing chan bool
		errorCatcher errors.ErrorCatcher
	}
)

// NewPublisher returns a Publisher
func NewPublisher(subscriptions Subscriptions, errorCatcher errors.ErrorCatcher) Publisher {
	errorschecker.CheckNilParameter(map[string]interface{}{"subscriptions": subscriptions, "errorCatcher": errorCatcher})
	return &publisher{subscriptions: subscriptions, errorCatcher: errorCatcher}
}

// Publish publishes a event
func (p *publisher) Publish(event EventObject) {
	errorschecker.CheckNilParameter(map[string]interface{}{"event": event})
	typ := reflect.Indirect(reflect.ValueOf(event)).Type()
	ePublisher, _ := p.eventPublishers.LoadOrStore(typ, &eventPublisher{isPublishing: make(chan bool, 1), errorCatcher: p.errorCatcher})
	ePublisher.(*eventPublisher).publish(event, p.subscriptions.GetAlls(event))

}

func (e *eventPublisher) publish(event EventObject, functions []reflect.Value) {
	e.isPublishing <- true
	wg := sync.WaitGroup{}
	for _, function := range functions {
		if event.StopPropagation() {
			break
		}
		wg.Add(1)
		e.errorCatcher.OnErrorContinue(func() {
			callback := func(function reflect.Value) {
				if event.HasEventArgs() {
					function.Call([]reflect.Value{
						reflect.ValueOf(event.GetEventArgs()),
					})
				} else {
					function.Call(make([]reflect.Value, 0))
				}
				wg.Done()
			}
			if event.IsParallelPropagation() {
				go callback(function)
			} else {
				callback(function)
			}
		})
	}
	wg.Wait()
	<-e.isPublishing
}
