package redux

import (
	"github.com/janmbaco/go-infrastructure/errorhandler"
	"github.com/janmbaco/go-infrastructure/events"
)

const onNewStae = "onNewState"

type StateManager interface {
	GetState() interface{}
	SetState(interface{})
	Subscribe(fn func())
}

type stateManager struct {
	publisher events.EventPublisher
	state     interface{}
}

func NewStateManager(publisher events.EventPublisher, state interface{}) StateManager {
	errorhandler.CheckNilParameter(map[string]interface{}{"publisher": publisher, "state": state})
	return &stateManager{publisher: publisher, state: state}
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
