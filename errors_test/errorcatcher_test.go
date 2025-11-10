package errors_test

import (
	"fmt"
	"testing"

	"github.com/janmbaco/go-infrastructure/v2/errors"
	"github.com/janmbaco/go-infrastructure/v2/logs"
)

// mockLogger is a test logger implementation
type mockLogger struct {
	lastError   error
	errorCalled bool
}

func (m *mockLogger) PrintError(level logs.LogLevel, err error)    {}
func (m *mockLogger) TryPrintError(level logs.LogLevel, err error) {}
func (m *mockLogger) TryTrace(err error)                           {}
func (m *mockLogger) TryInfo(err error)                            {}
func (m *mockLogger) TryWarning(err error)                         {}
func (m *mockLogger) TryError(err error) {
	m.errorCalled = true
	m.lastError = err
}
func (m *mockLogger) TryFatal(err error) {}

// Test_NewErrorCatcher_WhenLoggerIsNil_ThenReturnsValidCatcher validates that catcher works with nil logger
func Test_NewErrorCatcher_WhenLoggerIsNil_ThenReturnsValidCatcher(t *testing.T) {
	// Arrange & Act
	catcher := errors.NewErrorCatcher(nil)

	// Assert
	if catcher == nil {
		t.Fatal("Expected valid catcher, got nil")
	}
}

// Test_NewErrorCatcher_WhenLoggerIsValid_ThenReturnsValidCatcher validates that catcher is created with logger
func Test_NewErrorCatcher_WhenLoggerIsValid_ThenReturnsValidCatcher(t *testing.T) {
	// Arrange
	logger := &mockLogger{}

	// Act
	catcher := errors.NewErrorCatcher(logger)

	// Assert
	if catcher == nil {
		t.Fatal("Expected valid catcher, got nil")
	}
}

// Test_HandleError_WhenErrorIsNil_ThenReturnsNil validates that no error is returned for nil error
func Test_HandleError_WhenErrorIsNil_ThenReturnsNil(t *testing.T) {
	// Arrange
	catcher := errors.NewErrorCatcher(nil)
	var errorFnCalled bool
	errorFn := func(err error) { errorFnCalled = true }

	// Act
	err := catcher.HandleError(nil, errorFn)

	// Assert
	if err != nil {
		t.Fatalf("Expected nil error, got: %v", err)
	}
	if errorFnCalled {
		t.Fatal("Expected errorFn not to be called for nil error")
	}
}

// Test_HandleError_WhenErrorExists_ThenCallsErrorFn validates that errorFn is called for non-nil error
func Test_HandleError_WhenErrorExists_ThenCallsErrorFn(t *testing.T) {
	// Arrange
	catcher := errors.NewErrorCatcher(nil)
	testError := fmt.Errorf("test error")
	var capturedError error
	errorFn := func(err error) { capturedError = err }

	// Act
	returnedErr := catcher.HandleError(testError, errorFn)

	// Assert
	if returnedErr != testError {
		t.Fatalf("Expected returned error to be %v, got: %v", testError, returnedErr)
	}
	if capturedError != testError {
		t.Fatalf("Expected captured error to be %v, got: %v", testError, capturedError)
	}
}

// Test_HandleError_WhenErrorFnIsNil_ThenReturnsError validates that error is returned even with nil errorFn
func Test_HandleError_WhenErrorFnIsNil_ThenReturnsError(t *testing.T) {
	// Arrange
	catcher := errors.NewErrorCatcher(nil)
	testError := fmt.Errorf("test error")

	// Act
	err := catcher.HandleError(testError, nil)

	// Assert
	if err != testError {
		t.Fatalf("Expected error %v, got: %v", testError, err)
	}
}

// Test_HandleErrorWithFinally_WhenErrorExists_ThenCallsBothFunctions validates that both functions are called
func Test_HandleErrorWithFinally_WhenErrorExists_ThenCallsBothFunctions(t *testing.T) {
	// Arrange
	catcher := errors.NewErrorCatcher(nil)
	testError := fmt.Errorf("test error")
	var errorFnCalled bool
	var finallyFnCalled bool
	errorFn := func(err error) { errorFnCalled = true }
	finallyFn := func() { finallyFnCalled = true }

	// Act
	err := catcher.HandleErrorWithFinally(testError, errorFn, finallyFn)

	// Assert
	if err != testError {
		t.Fatalf("Expected error %v, got: %v", testError, err)
	}
	if !errorFnCalled {
		t.Fatal("Expected errorFn to be called")
	}
	if !finallyFnCalled {
		t.Fatal("Expected finallyFn to be called")
	}
}

// Test_HandleErrorWithFinally_WhenErrorIsNil_ThenCallsOnlyFinally validates that only finally is called for nil error
func Test_HandleErrorWithFinally_WhenErrorIsNil_ThenCallsOnlyFinally(t *testing.T) {
	// Arrange
	catcher := errors.NewErrorCatcher(nil)
	var errorFnCalled bool
	var finallyFnCalled bool
	errorFn := func(err error) { errorFnCalled = true }
	finallyFn := func() { finallyFnCalled = true }

	// Act
	err := catcher.HandleErrorWithFinally(nil, errorFn, finallyFn)

	// Assert
	if err != nil {
		t.Fatalf("Expected nil error, got: %v", err)
	}
	if errorFnCalled {
		t.Fatal("Expected errorFn not to be called for nil error")
	}
	if !finallyFnCalled {
		t.Fatal("Expected finallyFn to be called")
	}
}

// Test_TryCatchError_WhenTryFnIsNil_ThenReturnsError validates that error is returned for nil tryFn
func Test_TryCatchError_WhenTryFnIsNil_ThenReturnsError(t *testing.T) {
	// Arrange
	catcher := errors.NewErrorCatcher(nil)

	// Act
	err := catcher.TryCatchError(nil, nil)

	// Assert
	if err == nil {
		t.Fatal("Expected error for nil tryFn, got nil")
	}
	if err.Error() != "tryfn cannot be nil" {
		t.Fatalf("Expected 'tryfn cannot be nil', got: %v", err)
	}
}

// Test_TryCatchError_WhenTryFnReturnsError_ThenCallsErrorFn validates that errorFn is called when tryFn returns error
func Test_TryCatchError_WhenTryFnReturnsError_ThenCallsErrorFn(t *testing.T) {
	// Arrange
	catcher := errors.NewErrorCatcher(nil)
	testError := fmt.Errorf("test error")
	tryFn := func() error { return testError }
	var capturedError error
	errorFn := func(err error) { capturedError = err }

	// Act
	returnedErr := catcher.TryCatchError(tryFn, errorFn)

	// Assert
	if returnedErr != testError {
		t.Fatalf("Expected returned error to be %v, got: %v", testError, returnedErr)
	}
	if capturedError != testError {
		t.Fatalf("Expected captured error to be %v, got: %v", testError, capturedError)
	}
}

// Test_TryCatchError_WhenTryFnReturnsNil_ThenDoesNotCallErrorFn validates that errorFn is not called for nil error
func Test_TryCatchError_WhenTryFnReturnsNil_ThenDoesNotCallErrorFn(t *testing.T) {
	// Arrange
	catcher := errors.NewErrorCatcher(nil)
	tryFn := func() error { return nil }
	var errorFnCalled bool
	errorFn := func(err error) { errorFnCalled = true }

	// Act
	err := catcher.TryCatchError(tryFn, errorFn)

	// Assert
	if err != nil {
		t.Fatalf("Expected nil error, got: %v", err)
	}
	if errorFnCalled {
		t.Fatal("Expected errorFn not to be called for nil error")
	}
}

// Test_TryFinally_WhenTryFnIsNil_ThenReturnsError validates that error is returned for nil tryFn
func Test_TryFinally_WhenTryFnIsNil_ThenReturnsError(t *testing.T) {
	// Arrange
	catcher := errors.NewErrorCatcher(nil)

	// Act
	err := catcher.TryFinally(nil, nil)

	// Assert
	if err == nil {
		t.Fatal("Expected error for nil tryFn, got nil")
	}
	if err.Error() != "tryfn cannot be nil" {
		t.Fatalf("Expected 'tryfn cannot be nil', got: %v", err)
	}
}

// Test_TryFinally_WhenTryFnExecutes_ThenCallsFinallyFn validates that finallyFn is always called
func Test_TryFinally_WhenTryFnExecutes_ThenCallsFinallyFn(t *testing.T) {
	// Arrange
	catcher := errors.NewErrorCatcher(nil)
	var tryFnCalled bool
	var finallyFnCalled bool
	tryFn := func() error { tryFnCalled = true; return nil }
	finallyFn := func() { finallyFnCalled = true }

	// Act
	err := catcher.TryFinally(tryFn, finallyFn)

	// Assert
	if err != nil {
		t.Fatalf("Expected nil error, got: %v", err)
	}
	if !tryFnCalled {
		t.Fatal("Expected tryFn to be called")
	}
	if !finallyFnCalled {
		t.Fatal("Expected finallyFn to be called")
	}
}

// Test_TryFinally_WhenTryFnReturnsError_ThenStillCallsFinallyFn validates that finallyFn is called even on error
func Test_TryFinally_WhenTryFnReturnsError_ThenStillCallsFinallyFn(t *testing.T) {
	// Arrange
	catcher := errors.NewErrorCatcher(nil)
	testError := fmt.Errorf("test error")
	var finallyFnCalled bool
	tryFn := func() error { return testError }
	finallyFn := func() { finallyFnCalled = true }

	// Act
	err := catcher.TryFinally(tryFn, finallyFn)

	// Assert
	if err != testError {
		t.Fatalf("Expected error %v, got: %v", testError, err)
	}
	if !finallyFnCalled {
		t.Fatal("Expected finallyFn to be called even on error")
	}
}

// Test_TryCatchErrorAndFinally_WhenAllFunctionsProvided_ThenCallsInCorrectOrder validates execution order
func Test_TryCatchErrorAndFinally_WhenAllFunctionsProvided_ThenCallsInCorrectOrder(t *testing.T) {
	// Arrange
	catcher := errors.NewErrorCatcher(nil)
	testError := fmt.Errorf("test error")
	var callOrder []string
	tryFn := func() error { callOrder = append(callOrder, "try"); return testError }
	errorFn := func(err error) { callOrder = append(callOrder, "error") }
	finallyFn := func() { callOrder = append(callOrder, "finally") }

	// Act
	err := catcher.TryCatchErrorAndFinally(tryFn, errorFn, finallyFn)

	// Assert
	if err != testError {
		t.Fatalf("Expected error %v, got: %v", testError, err)
	}
	if len(callOrder) != 3 {
		t.Fatalf("Expected 3 calls, got: %d", len(callOrder))
	}
	// Note: finally is called via defer, so order is: try -> error -> finally (defer executes last)
	if callOrder[0] != "try" || callOrder[1] != "error" || callOrder[2] != "finally" {
		t.Fatalf("Expected order [try, error, finally], got: %v", callOrder)
	}
}

// Test_OnErrorContinue_WhenTryFnReturnsError_ThenLogsErrorAndReturnsNil validates that error is logged but nil returned
func Test_OnErrorContinue_WhenTryFnReturnsError_ThenLogsErrorAndReturnsNil(t *testing.T) {
	// Arrange
	logger := &mockLogger{}
	catcher := errors.NewErrorCatcher(logger)
	testError := fmt.Errorf("test error")
	tryFn := func() error { return testError }

	// Act
	err := catcher.OnErrorContinue(tryFn)

	// Assert
	if err != nil {
		t.Fatalf("Expected nil error, got: %v", err)
	}
	if !logger.errorCalled {
		t.Fatal("Expected logger.TryError to be called")
	}
	if logger.lastError != testError {
		t.Fatalf("Expected logged error to be %v, got: %v", testError, logger.lastError)
	}
}

// Test_OnErrorContinue_WhenTryFnIsNil_ThenReturnsError validates that error is returned for nil tryFn
func Test_OnErrorContinue_WhenTryFnIsNil_ThenReturnsError(t *testing.T) {
	// Arrange
	catcher := errors.NewErrorCatcher(nil)

	// Act
	err := catcher.OnErrorContinue(nil)

	// Assert
	if err == nil {
		t.Fatal("Expected error for nil tryFn, got nil")
	}
	if err.Error() != "tryfn cannot be nil" {
		t.Fatalf("Expected 'tryfn cannot be nil', got: %v", err)
	}
}
