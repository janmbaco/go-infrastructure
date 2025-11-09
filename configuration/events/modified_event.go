package events

// ModifiedEvent is the event that happen when modified occurs
type ModifiedEvent struct{}

// GetEventArgs gets the args of a event
func (e ModifiedEvent) GetEventArgs() ModifiedEvent {
	return e
}

// StopPropagation stops the propagation of a event
func (e ModifiedEvent) StopPropagation() bool {
	return false
}

// IsParallelPropagation indicates if the propagation of the event is in parallel
func (e ModifiedEvent) IsParallelPropagation() bool {
	return true
}
