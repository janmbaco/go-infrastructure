package events

// RestoredEvent is the event that happen when restored occurs
type RestoredEvent struct{}

// GetEventArgs gets the args of a event
func (e RestoredEvent) GetEventArgs() RestoredEvent {
	return e
}

// StopPropagation stops the propagation of a event
func (e RestoredEvent) StopPropagation() bool {
	return false
}

// IsParallelPropagation indicates if the propagation of the event is in parallel
func (e RestoredEvent) IsParallelPropagation() bool {
	return true
}
