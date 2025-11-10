package eventsmanager

// EventObject is the definition of an object responsible to make an event
type EventObject[T any] interface {
	GetEventArgs() T
	StopPropagation() bool
	IsParallelPropagation() bool
}
