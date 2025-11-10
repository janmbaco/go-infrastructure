package eventsmanager

import (
	"sync"

	"github.com/janmbaco/go-infrastructure/v2/logs"
)

// Publisher defines an object responsible to publish events
type Publisher[T EventObject[T]] interface {
	Publish(event T)
}

type publisher[T EventObject[T]] struct {
	subscriptions SubscriptionsGetter[T]
	logger        logs.Logger
}

// NewPublisher returns a Publisher
func NewPublisher[T EventObject[T]](subscriptions SubscriptionsGetter[T], logger logs.Logger) Publisher[T] {
	return &publisher[T]{subscriptions: subscriptions, logger: logger}
}

// Publish publishes an event
func (p *publisher[T]) Publish(event T) {
	p.publishEvent(event, p.subscriptions.GetAlls())
}

func (p *publisher[T]) publishEvent(event T, functions []func(T)) {
	if event.IsParallelPropagation() {
		p.publishInParallel(event, functions)
	} else {
		p.publishSequentially(event, functions)
	}
}

func (p *publisher[T]) publishInParallel(event T, functions []func(T)) {
	var wg sync.WaitGroup
	p.iterateAndExecute(event, functions, func(fn func(T)) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.executeWithPanicRecovery(fn, event, nil)
		}()
	})
	wg.Wait()
}

func (p *publisher[T]) publishSequentially(event T, functions []func(T)) {
	p.iterateAndExecute(event, functions, func(fn func(T)) {
		p.executeWithPanicRecovery(fn, event, nil)
	})
}

func (p *publisher[T]) iterateAndExecute(event T, functions []func(T), executor func(func(T))) {
	for _, function := range functions {
		if event.StopPropagation() {
			break
		}
		executor(function)
	}
}

func (p *publisher[T]) executeWithPanicRecovery(fn func(T), event T, done func()) {
	defer func() {
		if r := recover(); r != nil {
			if p.logger != nil {
				p.logger.Errorf("panic in event handler: %v", r)
			}
		}
		if done != nil {
			done()
		}
	}()
	fn(event)
}
