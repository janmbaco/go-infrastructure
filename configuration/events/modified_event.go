package events

import (
	"reflect"
)

// ModifiedEvent is the event that happen when modified occurs
type ModifiedEvent struct{}

// IsParallelPropagation indicates if the propagation of the event is in parallel
func (*ModifiedEvent) IsParallelPropagation() bool {
	return true
}

// GetTypeOfFunc gets the signature of the function of the event
func (*ModifiedEvent) GetTypeOfFunc() reflect.Type {
	return reflect.TypeOf(func() {})
}

// StopPropagation stops the propagation of a event
func (*ModifiedEvent) StopPropagation() bool {
	return false
}

// HasEventArgs indicates if a event has args
func (*ModifiedEvent) HasEventArgs() bool {
	return false
}

// GetEventArgs gets the args of a event
func (*ModifiedEvent) GetEventArgs() interface{} {
	return nil
}
