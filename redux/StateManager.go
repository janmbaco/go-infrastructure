package redux

import "github.com/janmbaco/go-infrastructure/event"

const onNewStae = "onNewState"

type StateManager interface {
	GetState() interface{}
	SetState(interface{})
	Subscribe(fn func())
}

type stateManager struct {
	publisher *event.EventPublisher
	state     interface{}
}

func (s *stateManager) GetState() interface{} {
	return s.state
}
func (s *stateManager) SetState(newState interface{}) {
	s.state = newState
	s.publisher.Publish(onNewStae)
}
func (s *stateManager) Subscribe(fn func()) {
	s.publisher.Subscribe(onNewStae, fn)
}
