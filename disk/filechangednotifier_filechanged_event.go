package disk

// FileChangedEvent represents a file change event
type FileChangedEvent struct{}

// GetEventArgs returns the event args
func (e FileChangedEvent) GetEventArgs() FileChangedEvent {
	return e
}

// StopPropagation indicates if propagation should stop
func (e FileChangedEvent) StopPropagation() bool {
	return false
}

// IsParallelPropagation indicates if should publish in parallel
func (e FileChangedEvent) IsParallelPropagation() bool {
	return false
}
