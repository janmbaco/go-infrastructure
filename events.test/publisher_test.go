package events_test

import (
	"github.com/janmbaco/go-infrastructure/events"
	"github.com/janmbaco/go-infrastructure/logs"
	"testing"
)

type subsciption struct {
	fn1 func()
}

func (s *subsciption) Initialize(name string) {
	s.fn1 = func() {
		logs.Log.Info(name + " onLaunchEvent_")
	}
}

func TestNewEventsPublisher(t *testing.T) {

	onLaunch := "onLaunch"
	sub1 := &subsciption{}
	sub1.Initialize("sub1")
	sub2 := &subsciption{}
	sub2.Initialize("sub2")
	publisher := events.NewPublisher()

	publisher.Subscribe(onLaunch, &sub1.fn1)
	publisher.Subscribe(onLaunch, &sub2.fn1)

	publisher.Publish(onLaunch)

	publisher.Subscribe(onLaunch, &sub1.fn1)
	publisher.Subscribe(onLaunch, &sub2.fn1)

	publisher.Publish(onLaunch)

	publisher.UnSubscribe(onLaunch, &sub1.fn1)

	publisher.Publish(onLaunch)

	publisher.UnSubscribe(onLaunch, &sub2.fn1)

}
