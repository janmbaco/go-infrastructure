package eventsmanager

import (
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"reflect"
	"sync"
)

type SubscriptionsGetter interface {
	GetAlls(event EventObject) []reflect.Value
}

// Subscriptions is the definition of a object responsible to store subscriptions for an event
type Subscriptions interface {
	SubscriptionsGetter
	Add(event EventObject, subscribeFunc interface{})
	Remove(event EventObject, subscribeFunc interface{})
}

type subscriptions struct {
	events     sync.Map
	errorDefer errors.ErrorDefer
}

// NewSubscriptions returns a Subscriptions
func NewSubscriptions(errorDefer errors.ErrorDefer) Subscriptions {
	errorschecker.CheckNilParameter(map[string]interface{}{"errorDefer": errorDefer})
	return &subscriptions{errorDefer: errorDefer}
}

// Add sets a subscription to an event
func (s *subscriptions) Add(event EventObject, subscribeFunc interface{}) {
	defer s.errorDefer.TryThrowError(s.pipeError)
	errorschecker.CheckNilParameter(map[string]interface{}{"event": event, "subscribeFunc": subscribeFunc})
	functionValue := reflect.Indirect(reflect.ValueOf(subscribeFunc))
	functionType := reflect.Indirect(reflect.ValueOf(subscribeFunc)).Type()
	if functionType != event.GetTypeOfFunc() {
		panic(newSubscriptionsError(BadFunctionSignature, "the function hasn´t the correct signature", nil))
	}
	pointer := reflect.ValueOf(subscribeFunc).Pointer()

	var subscription interface{}
	var isContained bool

	typ := reflect.Indirect(reflect.ValueOf(event)).Type()
	if subscription, isContained = s.events.Load(typ); !isContained {
		subscription, _ = s.events.LoadOrStore(typ, &sync.Map{})
	}
	subscription.(*sync.Map).Store(pointer, functionValue)

}

// Remove deletes a subscription to an event
func (s *subscriptions) Remove(event EventObject, subscribeFunc interface{}) {
	defer s.errorDefer.TryThrowError(s.pipeError)
	errorschecker.CheckNilParameter(map[string]interface{}{"event": event, "subscribeFunc": subscribeFunc})
	pointer := reflect.ValueOf(subscribeFunc).Pointer()
	typ := reflect.Indirect(reflect.ValueOf(event)).Type()
	showError := true
	if subscriptions, isContained := s.events.Load(typ); isContained {
		if _, isRegistered := subscriptions.(*sync.Map).LoadAndDelete(pointer); isRegistered {
			showError = false
		}
	}

	if showError {
		panic(newSubscriptionsError(FunctionNoSubscribed, "this function is not registered", nil))
	}
}

// GetAlls gets all subscriptions for an event
func (s *subscriptions) GetAlls(event EventObject) []reflect.Value {
	defer s.errorDefer.TryThrowError(s.pipeError)
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

func (s *subscriptions) pipeError(err error) error {
	resultError := err
	if errType := reflect.Indirect(reflect.ValueOf(err)).Type(); !errType.Implements(reflect.TypeOf((*SubscriptionsError)(nil)).Elem()) {
		resultError = newSubscriptionsError(Unexpected, err.Error(), err)
	}
	return resultError
}