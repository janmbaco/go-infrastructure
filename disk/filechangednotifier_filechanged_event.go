package disk

import "reflect"

// fileChangedEvent is the event that happen when changes in file occurs
type fileChangedEvent struct{}

// IsParallelPropagation indicates if the propagation of the event is in parallel
func (*fileChangedEvent) IsParallelPropagation() bool {
	return true
}

// GetTypeOfFunc gets the signature of the function of the event
func (*fileChangedEvent) GetTypeOfFunc() reflect.Type {
	return reflect.TypeOf(func() {})
}

// StopPropagation stops the propagation of a event
func (*fileChangedEvent) StopPropagation() bool {
	return false
}

// GetEventArgs gets the args of a event
func (*fileChangedEvent) GetEventArgs() interface{} {
	return nil
}

// HasEventArgs indicates if a event has args
func (*fileChangedEvent) HasEventArgs() bool {
	return false
}
