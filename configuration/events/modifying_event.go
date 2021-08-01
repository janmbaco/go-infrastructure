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

type ModifyingEvent struct {
	EventArgs *ModifyingEventArgs
}

func (*ModifyingEvent) IsParallelPropagation() bool {
	return false
}

func (*ModifyingEvent) GetTypeOfFunc() reflect.Type {
	return reflect.TypeOf(func(args *ModifyingEventArgs) {})
}

func (m *ModifyingEvent) StopPropagation() bool {
	return m.EventArgs.Cancel
}

func (m *ModifyingEvent) GetEventArgs() interface{} {
	return m.EventArgs
}

func (*ModifyingEvent) HasEventArgs() bool {
	return true
}
