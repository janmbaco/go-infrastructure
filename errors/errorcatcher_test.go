package errors

import (
	"errors"
	"testing"

	"github.com/janmbaco/go-infrastructure/v2/logs"
	"github.com/stretchr/testify/assert"
)

// Simple mock logger for testing
type simpleMockErrorLogger struct {
	tryErrorCalled bool
	lastError      error
}

func (m *simpleMockErrorLogger) PrintError(level logs.LogLevel, err error)    {}
func (m *simpleMockErrorLogger) TryPrintError(level logs.LogLevel, err error) {}
func (m *simpleMockErrorLogger) TryTrace(err error)                           {}
func (m *simpleMockErrorLogger) TryInfo(err error)                            {}
func (m *simpleMockErrorLogger) TryWarning(err error)                         {}
func (m *simpleMockErrorLogger) TryError(err error) {
	m.tryErrorCalled = true
	m.lastError = err
}
func (m *simpleMockErrorLogger) TryFatal(err error) {}

func TestNewErrorCatcher_WhenLoggerNil_ThenUsesNoOpLogger(t *testing.T) {
	// Arrange & Act
	catcher := NewErrorCatcher(nil)

	// Assert
	assert.NotNil(t, catcher)
	assert.IsType(t, &errorCatcher{}, catcher)
}

func TestNewErrorCatcher_WhenLoggerProvided_ThenUsesProvidedLogger(t *testing.T) {
	// Arrange
	logger := &simpleMockErrorLogger{}

	// Act
	catcher := NewErrorCatcher(logger)

	// Assert
	assert.NotNil(t, catcher)
	assert.IsType(t, &errorCatcher{}, catcher)
}

func TestErrorCatcher_HandleError_WhenErrNil_ThenReturnsNil(t *testing.T) {
	// Arrange
	catcher := NewErrorCatcher(nil)
	called := false
	errorfn := func(error) { called = true }

	// Act
	result := catcher.HandleError(nil, errorfn)

	// Assert
	assert.Nil(t, result)
	assert.False(t, called)
}

func TestErrorCatcher_HandleError_WhenErrNotNil_ThenCallsErrorfnAndReturnsErr(t *testing.T) {
	// Arrange
	catcher := NewErrorCatcher(nil)
	testErr := errors.New("test")
	called := false
	errorfn := func(e error) {
		called = true
		assert.Equal(t, testErr, e)
	}

	// Act
	result := catcher.HandleError(testErr, errorfn)

	// Assert
	assert.Equal(t, testErr, result)
	assert.True(t, called)
}

func TestErrorCatcher_HandleError_WhenErrorfnNil_ThenReturnsErr(t *testing.T) {
	// Arrange
	catcher := NewErrorCatcher(nil)
	testErr := errors.New("test")

	// Act
	result := catcher.HandleError(testErr, nil)

	// Assert
	assert.Equal(t, testErr, result)
}

func TestErrorCatcher_HandleErrorWithFinally_WhenCalled_ThenExecutesFinallyAndHandlesError(t *testing.T) {
	// Arrange
	catcher := NewErrorCatcher(nil)
	testErr := errors.New("test")
	finallyCalled := false
	errorCalled := false
	finallyfn := func() { finallyCalled = true }
	errorfn := func(error) { errorCalled = true }

	// Act
	result := catcher.HandleErrorWithFinally(testErr, errorfn, finallyfn)

	// Assert
	assert.Equal(t, testErr, result)
	assert.True(t, finallyCalled)
	assert.True(t, errorCalled)
}

func TestErrorCatcher_TryCatchError_WhenTryfnNil_ThenReturnsError(t *testing.T) {
	// Arrange
	catcher := NewErrorCatcher(nil)

	// Act
	result := catcher.TryCatchError(nil, nil)

	// Assert
	assert.Error(t, result)
	assert.Contains(t, result.Error(), "tryfn cannot be nil")
}

func TestErrorCatcher_TryCatchError_WhenTryfnReturnsNil_ThenReturnsNil(t *testing.T) {
	// Arrange
	catcher := NewErrorCatcher(nil)
	tryfn := func() error { return nil }
	called := false
	errorfn := func(error) { called = true }

	// Act
	result := catcher.TryCatchError(tryfn, errorfn)

	// Assert
	assert.Nil(t, result)
	assert.False(t, called)
}

func TestErrorCatcher_TryCatchError_WhenTryfnReturnsError_ThenCallsErrorfnAndReturnsError(t *testing.T) {
	// Arrange
	catcher := NewErrorCatcher(nil)
	testErr := errors.New("test")
	tryfn := func() error { return testErr }
	called := false
	errorfn := func(e error) {
		called = true
		assert.Equal(t, testErr, e)
	}

	// Act
	result := catcher.TryCatchError(tryfn, errorfn)

	// Assert
	assert.Equal(t, testErr, result)
	assert.True(t, called)
}

func TestErrorCatcher_TryFinally_WhenTryfnNil_ThenReturnsError(t *testing.T) {
	// Arrange
	catcher := NewErrorCatcher(nil)

	// Act
	result := catcher.TryFinally(nil, nil)

	// Assert
	assert.Error(t, result)
	assert.Contains(t, result.Error(), "tryfn cannot be nil")
}

func TestErrorCatcher_TryFinally_WhenCalled_ThenExecutesFinally(t *testing.T) {
	// Arrange
	catcher := NewErrorCatcher(nil)
	tryfn := func() error { return nil }
	finallyCalled := false
	finallyfn := func() { finallyCalled = true }

	// Act
	result := catcher.TryFinally(tryfn, finallyfn)

	// Assert
	assert.Nil(t, result)
	assert.True(t, finallyCalled)
}

func TestErrorCatcher_OnErrorContinue_WhenTryfnNil_ThenReturnsError(t *testing.T) {
	// Arrange
	catcher := NewErrorCatcher(nil)

	// Act
	result := catcher.OnErrorContinue(nil)

	// Assert
	assert.Error(t, result)
	assert.Contains(t, result.Error(), "tryfn cannot be nil")
}

func TestErrorCatcher_OnErrorContinue_WhenTryfnReturnsError_ThenLogsAndReturnsNil(t *testing.T) {
	// Arrange
	logger := &simpleMockErrorLogger{}
	catcher := NewErrorCatcher(logger)
	testErr := errors.New("test")
	tryfn := func() error { return testErr }

	// Act
	result := catcher.OnErrorContinue(tryfn)

	// Assert
	assert.Nil(t, result)
	assert.True(t, logger.tryErrorCalled)
	assert.Equal(t, testErr, logger.lastError)
}

func TestErrorCatcher_CatchError_WhenCalled_ThenDelegatesToHandleError(t *testing.T) {
	// Arrange
	catcher := NewErrorCatcher(nil)
	testErr := errors.New("test")
	called := false
	errorfn := func(error) { called = true }

	// Act
	result := catcher.CatchError(testErr, errorfn)

	// Assert
	assert.Equal(t, testErr, result)
	assert.True(t, called)
}

func TestErrorCatcher_CatchErrorAndFinally_WhenCalled_ThenDelegatesToHandleErrorWithFinally(t *testing.T) {
	// Arrange
	catcher := NewErrorCatcher(nil)
	testErr := errors.New("test")
	finallyCalled := false
	errorCalled := false
	finallyfn := func() { finallyCalled = true }
	errorfn := func(error) { errorCalled = true }

	// Act
	result := catcher.CatchErrorAndFinally(testErr, errorfn, finallyfn)

	// Assert
	assert.Equal(t, testErr, result)
	assert.True(t, finallyCalled)
	assert.True(t, errorCalled)
}
