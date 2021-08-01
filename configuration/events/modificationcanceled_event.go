package events

import (
	"reflect"
)

// ModificationCanceledEventArgs are the args for the event ModificationCanceled
type ModificationCanceledEventArgs struct {
	CancelMessage string
}

type ModificationCanceledEvent struct {
	EventArgs *ModificationCanceledEventArgs
}

func (m *ModificationCanceledEvent) IsParallelPropagation() bool {
	return true
}

func (m *ModificationCanceledEvent) GetTypeOfFunc() reflect.Type {
	return reflect.TypeOf(func(args *ModificationCanceledEventArgs) {})
}

func (m *ModificationCanceledEvent) StopPropagation() bool {
	return false
}

func (m *ModificationCanceledEvent) GetEventArgs() interface{} {
	return m.EventArgs
}

func (m *ModificationCanceledEvent) HasEventArgs() bool {
	return true
}
