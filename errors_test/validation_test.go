package errors_test

import (
	"fmt"
	"testing"

	"github.com/janmbaco/go-infrastructure/errors"
)

// Test_ValidateNotNil_WhenAllParametersValid_ThenReturnsNil validates that no error is returned for valid parameters
func Test_ValidateNotNil_WhenAllParametersValid_ThenReturnsNil(t *testing.T) {
	// Arrange
	validString := "test"
	validInt := 42
	params := map[string]interface{}{
		"stringParam": validString,
		"intParam":    validInt,
	}

	// Act
	err := errors.ValidateNotNil(params)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

// Test_ValidateNotNil_WhenParameterIsNil_ThenReturnsError validates that error is returned for nil parameter
func Test_ValidateNotNil_WhenParameterIsNil_ThenReturnsError(t *testing.T) {
	// Arrange
	params := map[string]interface{}{
		"validParam": "test",
		"nilParam":   nil,
	}

	// Act
	err := errors.ValidateNotNil(params)

	// Assert
	if err == nil {
		t.Fatal("Expected error for nil parameter, got nil")
	}
	if !containsString(err.Error(), "nilParam") {
		t.Fatalf("Expected error message to contain 'nilParam', got: %v", err)
	}
}

// Test_ValidateNotNil_WhenNilPointer_ThenReturnsError validates that error is returned for nil pointer
func Test_ValidateNotNil_WhenNilPointer_ThenReturnsError(t *testing.T) {
	// Arrange
	var nilPtr *string
	params := map[string]interface{}{
		"ptrParam": nilPtr,
	}

	// Act
	err := errors.ValidateNotNil(params)

	// Assert
	if err == nil {
		t.Fatal("Expected error for nil pointer, got nil")
	}
	if !containsString(err.Error(), "ptrParam") {
		t.Fatalf("Expected error message to contain 'ptrParam', got: %v", err)
	}
}

// Test_ValidateNotNil_WhenNilSlice_ThenReturnsError validates that error is returned for nil slice
func Test_ValidateNotNil_WhenNilSlice_ThenReturnsError(t *testing.T) {
	// Arrange
	var nilSlice []string
	params := map[string]interface{}{
		"sliceParam": nilSlice,
	}

	// Act
	err := errors.ValidateNotNil(params)

	// Assert
	if err == nil {
		t.Fatal("Expected error for nil slice, got nil")
	}
	if !containsString(err.Error(), "sliceParam") {
		t.Fatalf("Expected error message to contain 'sliceParam', got: %v", err)
	}
}

// Test_ValidateNotNil_WhenNilMap_ThenReturnsError validates that error is returned for nil map
func Test_ValidateNotNil_WhenNilMap_ThenReturnsError(t *testing.T) {
	// Arrange
	var nilMap map[string]string
	params := map[string]interface{}{
		"mapParam": nilMap,
	}

	// Act
	err := errors.ValidateNotNil(params)

	// Assert
	if err == nil {
		t.Fatal("Expected error for nil map, got nil")
	}
	if !containsString(err.Error(), "mapParam") {
		t.Fatalf("Expected error message to contain 'mapParam', got: %v", err)
	}
}

// Test_ValidateNotNil_WhenMultipleNilParameters_ThenReturnsErrorWithAllNames validates that all nil parameters are reported
func Test_ValidateNotNil_WhenMultipleNilParameters_ThenReturnsErrorWithAllNames(t *testing.T) {
	// Arrange
	params := map[string]interface{}{
		"param1": nil,
		"param2": "valid",
		"param3": nil,
	}

	// Act
	err := errors.ValidateNotNil(params)

	// Assert
	if err == nil {
		t.Fatal("Expected error for nil parameters, got nil")
	}
	if !containsString(err.Error(), "param1") || !containsString(err.Error(), "param3") {
		t.Fatalf("Expected error message to contain both nil parameter names, got: %v", err)
	}
}

// Test_ValidateNotNil_WhenNilInterface_ThenReturnsError validates that error is returned for nil interface
func Test_ValidateNotNil_WhenNilInterface_ThenReturnsError(t *testing.T) {
	// Arrange
	var nilInterface error
	params := map[string]interface{}{
		"interfaceParam": nilInterface,
	}

	// Act
	err := errors.ValidateNotNil(params)

	// Assert
	if err == nil {
		t.Fatal("Expected error for nil interface, got nil")
	}
	if !containsString(err.Error(), "interfaceParam") {
		t.Fatalf("Expected error message to contain 'interfaceParam', got: %v", err)
	}
}

// Test_ValidateNotNil_WhenNilChannel_ThenReturnsError validates that error is returned for nil channel
func Test_ValidateNotNil_WhenNilChannel_ThenReturnsError(t *testing.T) {
	// Arrange
	var nilChan chan int
	params := map[string]interface{}{
		"chanParam": nilChan,
	}

	// Act
	err := errors.ValidateNotNil(params)

	// Assert
	if err == nil {
		t.Fatal("Expected error for nil channel, got nil")
	}
	if !containsString(err.Error(), "chanParam") {
		t.Fatalf("Expected error message to contain 'chanParam', got: %v", err)
	}
}

// Test_ValidateNotNil_WhenNilFunc_ThenReturnsError validates that error is returned for nil function
func Test_ValidateNotNil_WhenNilFunc_ThenReturnsError(t *testing.T) {
	// Arrange
	var nilFunc func()
	params := map[string]interface{}{
		"funcParam": nilFunc,
	}

	// Act
	err := errors.ValidateNotNil(params)

	// Assert
	if err == nil {
		t.Fatal("Expected error for nil function, got nil")
	}
	if !containsString(err.Error(), "funcParam") {
		t.Fatalf("Expected error message to contain 'funcParam', got: %v", err)
	}
}

// Helper function
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && len(substr) > 0 && fmt.Sprintf("%s", s)[0:len(s)] != "" && 
		fmt.Sprintf("%s", s) != "" && s[0:1] != "" && 
		findInString(s, substr)))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
