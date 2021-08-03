package events

import (
	"reflect"
)

// ModificationCanceledEventArgs are the args for the event ModificationCanceled
type ModificationCanceledEventArgs struct {
	CancelMessage string
}

// ModificationCanceledEvent is the event that happen when cancelation occurs
type ModificationCanceledEvent struct {
	EventArgs *ModificationCanceledEventArgs
}

// IsParallelPropagation indicates if the propagation of the event is in parallel
func (m *ModificationCanceledEvent) IsParallelPropagation() bool {
	return true
}

// GetTypeOfFunc gets the signature of the function of the event
func (m *ModificationCanceledEvent) GetTypeOfFunc() reflect.Type {
	return reflect.TypeOf(func(args *ModificationCanceledEventArgs) {})
}

// StopPropagation stops the propagation of a event
func (m *ModificationCanceledEvent) StopPropagation() bool {
	return false
}

// GetEventArgs gets the args of a event
func (m *ModificationCanceledEvent) GetEventArgs() interface{} {
	return m.EventArgs
}

// HasEventArgs indicates if a event has args
func (m *ModificationCanceledEvent) HasEventArgs() bool {
	return true
}
