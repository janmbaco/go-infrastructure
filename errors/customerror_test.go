package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomizableError_Error_WhenCalled_ThenReturnsMessage(t *testing.T) {
	// Arrange
	msg := "test error"
	err := &CustomizableError{Message: msg}

	// Act
	result := err.Error()

	// Assert
	assert.Equal(t, msg, result)
}

func TestCustomizableError_GetMessage_WhenCalled_ThenReturnsMessage(t *testing.T) {
	// Arrange
	msg := "test message"
	err := &CustomizableError{Message: msg}

	// Act
	result := err.GetMessage()

	// Assert
	assert.Equal(t, msg, result)
}

func TestCustomizableError_GetInternalError_WhenCalled_ThenReturnsInternalError(t *testing.T) {
	// Arrange
	internal := errors.New("internal")
	err := &CustomizableError{InternalError: internal}

	// Act
	result := err.GetInternalError()

	// Assert
	assert.Equal(t, internal, result)
}

func TestCustomizableError_GetInternalError_WhenNil_ThenReturnsNil(t *testing.T) {
	// Arrange
	err := &CustomizableError{}

	// Act
	result := err.GetInternalError()

	// Assert
	assert.Nil(t, result)
}
