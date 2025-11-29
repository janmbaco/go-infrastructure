package eventsmanager

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock event for testing
type mockEvent struct {
	data string
}

func (e mockEvent) GetEventArgs() mockEvent {
	return e
}

func (e mockEvent) StopPropagation() bool {
	return false
}

func (e mockEvent) IsParallelPropagation() bool {
	return false
}

// Simple mock publisher for testing
type simpleMockPublisher struct {
	published []mockEvent
}

func (m *simpleMockPublisher) Publish(event mockEvent) {
	m.published = append(m.published, event)
}

func TestNewEventManager_WhenCreated_ThenReturnsManager(t *testing.T) {
	// Arrange & Act
	em := NewEventManager()

	// Assert
	assert.NotNil(t, em)
	assert.NotNil(t, em.publishers)
	assert.Empty(t, em.publishers)
}

func TestRegister_WhenRegistered_ThenPublisherIsStored(t *testing.T) {
	// Arrange
	em := NewEventManager()
	publisher := &simpleMockPublisher{}

	// Act
	Register(em, publisher)

	// Assert
	typ := reflect.TypeOf(mockEvent{})
	assert.Contains(t, em.publishers, typ)
	assert.Equal(t, publisher, em.publishers[typ])
}

func TestPublish_WhenPublisherRegistered_ThenPublishes(t *testing.T) {
	// Arrange
	em := NewEventManager()
	publisher := &simpleMockPublisher{}
	Register(em, publisher)
	event := mockEvent{data: "test"}

	// Act
	Publish(em, event)

	// Assert
	assert.Len(t, publisher.published, 1)
	assert.Equal(t, event, publisher.published[0])
}

func TestPublish_WhenNoPublisherRegistered_ThenDoesNothing(t *testing.T) {
	// Arrange
	em := NewEventManager()
	event := mockEvent{data: "test"}

	// Act
	Publish(em, event)

	// Assert
	// No panic, nothing happens
}
