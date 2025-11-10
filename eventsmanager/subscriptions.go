package eventsmanager

import (
	"reflect"
)

type SubscriptionsGetter[T EventObject[T]] interface {
	GetAlls() []func(T)
}

// Subscriptions is the definition of an object responsible to store subscriptions for an event
type Subscriptions[T EventObject[T]] interface {
	SubscriptionsGetter[T]
	Add(subscribeFunc func(T)) error
	Remove(subscribeFunc func(T)) error
}

type subscriptions[T EventObject[T]] struct {
	subs map[uintptr]func(T)
}

// NewSubscriptions returns a Subscriptions
func NewSubscriptions[T EventObject[T]]() Subscriptions[T] {
	return &subscriptions[T]{subs: make(map[uintptr]func(T))}
}

// Add sets a subscription to an event
func (s *subscriptions[T]) Add(subscribeFunc func(T)) error {
	pointer := reflect.ValueOf(subscribeFunc).Pointer()
	s.subs[pointer] = subscribeFunc
	return nil
}

// Remove deletes a subscription to an event
func (s *subscriptions[T]) Remove(subscribeFunc func(T)) error {
	pointer := reflect.ValueOf(subscribeFunc).Pointer()
	if _, exists := s.subs[pointer]; exists {
		delete(s.subs, pointer)
		return nil
	}
	return newSubscriptionsError(FunctionNoSubscribed, "function not subscribed", nil)
}

// GetAlls gets all subscriptions for an event
func (s *subscriptions[T]) GetAlls() []func(T) {
	result := make([]func(T), 0, len(s.subs))
	for _, fn := range s.subs {
		result = append(result, fn)
	}
	return result
}
