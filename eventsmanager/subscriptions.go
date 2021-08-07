package eventsmanager

import (
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"reflect"
	"sync"
)

// Subscriptions is the definition of a object responsible to store subscriptions for an event
type Subscriptions interface {
	Add(event EventObject, subscribeFunc interface{})
	Remove(event EventObject, subscribeFunc interface{})
	GetAlls(event EventObject) []reflect.Value
}

type subscriptions struct {
	events     sync.Map
	errorDefer errors.ErrorDefer
}

// NewSubscriptions returns a Subscriptions
func NewSubscriptions(thrower errors.ErrorThrower) Subscriptions {
	return &subscriptions{errorDefer: errors.NewErrorDefer(thrower, &subscriptionsErrorPipe{})}
}

// Add sets a subscription to an event
func (s *subscriptions) Add(event EventObject, subscribeFunc interface{}) {
	errorschecker.CheckNilParameter(map[string]interface{}{"event": event, "subscribeFunc": subscribeFunc})
	defer s.errorDefer.TryThrowError()
	functionValue := reflect.Indirect(reflect.ValueOf(subscribeFunc))
	functionType := reflect.Indirect(reflect.ValueOf(subscribeFunc)).Type()
	if functionType != event.GetTypeOfFunc() {
		panic(&SubscriptionsError{
			CustomError: errors.CustomError{
				Message:       "the function hasn´t the correct signature",
				InternalError: nil,
			},
			ErrorType: BadFunctionSignature,
		})
	}
	pointer := reflect.ValueOf(subscribeFunc).Pointer()

	var subscriptions interface{}
	var isContained bool

	typ := reflect.Indirect(reflect.ValueOf(event)).Type()
	if subscriptions, isContained = s.events.Load(typ); !isContained {
		subscriptions, _ = s.events.LoadOrStore(typ, &sync.Map{})
	}
	subscriptions.(*sync.Map).Store(pointer, functionValue)

}

// Remove deletes a subscription to an event
func (s *subscriptions) Remove(event EventObject, subscribeFunc interface{}) {
	errorschecker.CheckNilParameter(map[string]interface{}{"event": event, "subscribeFunc": subscribeFunc})
	defer s.errorDefer.TryThrowError()
	pointer := reflect.ValueOf(subscribeFunc).Pointer()
	typ := reflect.Indirect(reflect.ValueOf(event)).Type()
	showError := true
	if subscriptions, isContained := s.events.Load(typ); isContained {
		if _, isRegistered := subscriptions.(*sync.Map).LoadAndDelete(pointer); isRegistered {
			showError = false
		}
	}

	if showError {
		panic(&SubscriptionsError{
			CustomError: errors.CustomError{
				Message:       "this function is not registered",
				InternalError: nil,
			},
			ErrorType: FunctionNoSubscribed,
		})
	}
}

// GetAlls gets all subscriptions for an event
func (s *subscriptions) GetAlls(event EventObject) []reflect.Value {
	errorschecker.CheckNilParameter(map[string]interface{}{"event": event})
	result := make([]reflect.Value, 0)
	typ := reflect.Indirect(reflect.ValueOf(event)).Type()
	if subscriptions, isContained := s.events.Load(typ); isContained {
		subscriptions.(*sync.Map).Range(func(key, value interface{}) bool {
			result = append(result, value.(reflect.Value))
			return true
		})
	}
	return result
}