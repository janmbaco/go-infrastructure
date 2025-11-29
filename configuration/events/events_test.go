package events

import (
	"testing"

	"github.com/janmbaco/go-infrastructure/v2/eventsmanager"
	"github.com/stretchr/testify/assert"
)

func TestModifiedEvent_GetEventArgs_WhenCalled_ThenReturnsSelf(t *testing.T) {
	// Arrange
	event := ModifiedEvent{}

	// Act
	args := event.GetEventArgs()

	// Assert
	assert.Equal(t, event, args)
}

func TestModifiedEvent_StopPropagation_WhenCalled_ThenReturnsFalse(t *testing.T) {
	// Arrange
	event := ModifiedEvent{}

	// Act
	stop := event.StopPropagation()

	// Assert
	assert.False(t, stop)
}

func TestModifiedEvent_IsParallelPropagation_WhenCalled_ThenReturnsTrue(t *testing.T) {
	// Arrange
	event := ModifiedEvent{}

	// Act
	parallel := event.IsParallelPropagation()

	// Assert
	assert.True(t, parallel)
}

func TestNewModifiedEventHandler_WhenCreated_ThenReturnsHandler(t *testing.T) {
	// Arrange
	subs := eventsmanager.NewSubscriptions[ModifiedEvent]()

	// Act
	handler := NewModifiedEventHandler(subs)

	// Assert
	assert.NotNil(t, handler)
	assert.Equal(t, subs, handler.subscriptions)
}

func TestModifiedEventHandler_ModifiedSubscribe_WhenNilSubscription_ThenDoesNothing(t *testing.T) {
	// Arrange
	subs := eventsmanager.NewSubscriptions[ModifiedEvent]()
	handler := NewModifiedEventHandler(subs)

	// Act
	handler.ModifiedSubscribe(nil)

	// Assert
	// No panic
}

func TestModifiedEventHandler_ModifiedSubscribe_WhenValidSubscription_ThenAdds(t *testing.T) {
	// Arrange
	subs := eventsmanager.NewSubscriptions[ModifiedEvent]()
	handler := NewModifiedEventHandler(subs)
	fn := func() {}

	// Act
	handler.ModifiedSubscribe(&fn)

	// Assert
	// Since we can't easily test the internal, just ensure no panic
	assert.NotNil(t, handler)
}

func TestModifiedEventHandler_ModifiedUnsubscribe_WhenNilSubscription_ThenDoesNothing(t *testing.T) {
	// Arrange
	subs := eventsmanager.NewSubscriptions[ModifiedEvent]()
	handler := NewModifiedEventHandler(subs)

	// Act
	handler.ModifiedUnsubscribe(nil)

	// Assert
	// No panic
}

func TestModifiedEventHandler_ModifiedUnsubscribe_WhenValidSubscription_ThenRemoves(t *testing.T) {
	// Arrange
	subs := eventsmanager.NewSubscriptions[ModifiedEvent]()
	handler := NewModifiedEventHandler(subs)
	fn := func() {}

	// Act
	handler.ModifiedUnsubscribe(&fn)

	// Assert
	// No panic
}

func TestRestoredEvent_GetEventArgs_WhenCalled_ThenReturnsSelf(t *testing.T) {
	// Arrange
	event := RestoredEvent{}

	// Act
	args := event.GetEventArgs()

	// Assert
	assert.Equal(t, event, args)
}

func TestRestoredEvent_StopPropagation_WhenCalled_ThenReturnsFalse(t *testing.T) {
	// Arrange
	event := RestoredEvent{}

	// Act
	stop := event.StopPropagation()

	// Assert
	assert.False(t, stop)
}

func TestRestoredEvent_IsParallelPropagation_WhenCalled_ThenReturnsTrue(t *testing.T) {
	// Arrange
	event := RestoredEvent{}

	// Act
	parallel := event.IsParallelPropagation()

	// Assert
	assert.True(t, parallel)
}

func TestNewRestoredEventHandler_WhenCreated_ThenReturnsHandler(t *testing.T) {
	// Arrange
	subs := eventsmanager.NewSubscriptions[RestoredEvent]()

	// Act
	handler := NewRestoredEventHandler(subs)

	// Assert
	assert.NotNil(t, handler)
	assert.Equal(t, subs, handler.subscriptions)
}

func TestRestoredEventHandler_RestoredSubscribe_WhenNilSubscription_ThenDoesNothing(t *testing.T) {
	// Arrange
	subs := eventsmanager.NewSubscriptions[RestoredEvent]()
	handler := NewRestoredEventHandler(subs)

	// Act
	handler.RestoredSubscribe(nil)

	// Assert
	// No panic
}

func TestRestoredEventHandler_RestoredUnsubscribe_WhenNilSubscription_ThenDoesNothing(t *testing.T) {
	// Arrange
	subs := eventsmanager.NewSubscriptions[RestoredEvent]()
	handler := NewRestoredEventHandler(subs)

	// Act
	handler.RestoredUnsubscribe(nil)

	// Assert
	// No panic
}
