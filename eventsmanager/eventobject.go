package eventsmanager

import "reflect"

// EventObject is the definition of a object responsible to make an event
type EventObject interface {
	GetEventArgs() interface{}
	HasEventArgs() bool
	StopPropagation() bool
	IsParallelPropagation() bool
	GetTypeOfFunc() reflect.Type
}
