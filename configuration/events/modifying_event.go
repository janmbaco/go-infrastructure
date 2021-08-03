package events

import (
	"reflect"
)

// ModifyingEventArgs are the args for the event Modifying
type ModifyingEventArgs struct {
	Cancel        bool
	CancelMessage string
	Config        interface{}
}

// ModifyingEvent is the event that happen when modifying occurs
type ModifyingEvent struct {
	EventArgs *ModifyingEventArgs
}

// IsParallelPropagation indicates if the propagation of the event is in parallel
func (*ModifyingEvent) IsParallelPropagation() bool {
	return false
}

// GetTypeOfFunc gets the signature of the function of the event
func (*ModifyingEvent) GetTypeOfFunc() reflect.Type {
	return reflect.TypeOf(func(args *ModifyingEventArgs) {})
}

// StopPropagation stops the propagation of a event
func (m *ModifyingEvent) StopPropagation() bool {
	return m.EventArgs.Cancel
}

// GetEventArgs gets the args of a event
func (m *ModifyingEvent) GetEventArgs() interface{} {
	return m.EventArgs
}

// HasEventArgs indicates if a event has args
func (*ModifyingEvent) HasEventArgs() bool {
	return true
}
