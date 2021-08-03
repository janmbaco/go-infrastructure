package events

import (
	"reflect"
)

// RestoredEvent is the event that happen when restored occurs
type RestoredEvent struct{}

// IsParallelPropagation indicates if the propagation of the event is in parallel
func (*RestoredEvent) IsParallelPropagation() bool {
	return true
}

// GetTypeOfFunc gets the signature of the function of the event
func (*RestoredEvent) GetTypeOfFunc() reflect.Type {
	return reflect.TypeOf(func() {})
}

// StopPropagation stops the propagation of a event
func (*RestoredEvent) StopPropagation() bool {
	return false
}

// HasEventArgs indicates if a event has args
func (*RestoredEvent) HasEventArgs() bool {
	return false
}

// GetEventArgs gets the args of a event
func (*RestoredEvent) GetEventArgs() interface{} {
	return nil
}
