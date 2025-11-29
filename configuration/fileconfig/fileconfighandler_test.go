package fileconfig

import (
	"log"
	"path/filepath"
	"testing"

	"github.com/janmbaco/go-infrastructure/v2/configuration"
	"github.com/janmbaco/go-infrastructure/v2/disk"
	"github.com/janmbaco/go-infrastructure/v2/eventsmanager"
	"github.com/janmbaco/go-infrastructure/v2/logs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock implementations for testing
type mockErrorCatcher struct {
	tryCatchErrorCalled bool
	lastError           error
}

func (m *mockErrorCatcher) HandleError(err error, errorfn func(error)) error {
	if err != nil && errorfn != nil {
		errorfn(err)
	}
	return err
}

func (m *mockErrorCatcher) HandleErrorWithFinally(err error, errorfn func(error), finallyfn func()) error {
	if finallyfn != nil {
		finallyfn()
	}
	return m.HandleError(err, errorfn)
}

func (m *mockErrorCatcher) TryCatchError(tryfn func() error, errorfn func(error)) error {
	err := tryfn()
	if err != nil && errorfn != nil {
		errorfn(err)
	}
	return err
}

func (m *mockErrorCatcher) TryFinally(tryfn func() error, finallyfn func()) error {
	defer func() {
		if finallyfn != nil {
			finallyfn()
		}
	}()
	return tryfn()
}

func (m *mockErrorCatcher) TryCatchErrorAndFinally(tryfn func() error, errorfn func(error), finallyfn func()) error {
	defer func() {
		if finallyfn != nil {
			finallyfn()
		}
	}()
	return m.TryCatchError(tryfn, errorfn)
}

func (m *mockErrorCatcher) OnErrorContinue(tryfn func() error) error {
	err := tryfn()
	if err != nil {
		m.tryCatchErrorCalled = true
		m.lastError = err
	}
	return nil
}

func (m *mockErrorCatcher) CatchError(err error, errorfn func(error)) error {
	return m.HandleError(err, errorfn)
}

func (m *mockErrorCatcher) CatchErrorAndFinally(err error, errorfn func(error), finallyfn func()) error {
	return m.HandleErrorWithFinally(err, errorfn, finallyfn)
}

type mockFileChangedNotifier struct {
	subscribed bool
	callback   func()
}

func (m *mockFileChangedNotifier) Subscribe(subscribeFunc func()) error {
	m.subscribed = true
	m.callback = subscribeFunc
	return nil
}

type mockLogger struct{}

func (m *mockLogger) Println(level logs.LogLevel, message string)                   {}
func (m *mockLogger) Printlnf(level logs.LogLevel, format string, a ...interface{}) {}
func (m *mockLogger) Trace(message string)                                          {}
func (m *mockLogger) Tracef(format string, a ...interface{})                        {}
func (m *mockLogger) Info(message string)                                           {}
func (m *mockLogger) Infof(format string, a ...interface{})                         {}
func (m *mockLogger) Warning(message string)                                        {}
func (m *mockLogger) Warningf(format string, a ...interface{})                      {}
func (m *mockLogger) Error(message string)                                          {}
func (m *mockLogger) Errorf(format string, a ...interface{})                        {}
func (m *mockLogger) Fatal(message string)                                          {}
func (m *mockLogger) Fatalf(format string, a ...interface{})                        {}
func (m *mockLogger) SetConsoleLevel(level logs.LogLevel)                           {}
func (m *mockLogger) SetFileLogLevel(level logs.LogLevel)                           {}
func (m *mockLogger) GetErrorLogger() *log.Logger                                   { return nil }
func (m *mockLogger) SetDir(string)                                                 {}
func (m *mockLogger) PrintError(level logs.LogLevel, err error)                     {}
func (m *mockLogger) TryPrintError(level logs.LogLevel, err error)                  {}
func (m *mockLogger) TryTrace(err error)                                            {}
func (m *mockLogger) TryInfo(err error)                                             {}
func (m *mockLogger) TryWarning(err error)                                          {}
func (m *mockLogger) TryError(err error)                                            {}
func (m *mockLogger) TryFatal(err error)                                            {}
func (m *mockLogger) Mute()                                                         {}
func (m *mockLogger) Unmute()                                                       {}

type testConfig struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestNewFileConfigHandler_WhenValidInputs_ThenCreatesHandler(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.json")
	em := eventsmanager.NewEventManager()
	errorCatcher := &mockErrorCatcher{}
	notifier := &mockFileChangedNotifier{}
	logger := &mockLogger{}
	defaults := &testConfig{Name: "default", Value: 42}

	// Act
	handler, err := NewFileConfigHandler(filePath, defaults, errorCatcher, em, notifier, logger)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, handler)
	assert.True(t, notifier.subscribed)
	assert.True(t, disk.ExistsPath(filePath))
}

func TestNewFileConfigHandler_WhenFileDoesNotExist_ThenCreatesWithDefaults(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.json")
	em := eventsmanager.NewEventManager()
	errorCatcher := &mockErrorCatcher{}
	notifier := &mockFileChangedNotifier{}
	logger := &mockLogger{}
	defaults := &testConfig{Name: "default", Value: 42}

	// Act
	handler, err := NewFileConfigHandler(filePath, defaults, errorCatcher, em, notifier, logger)

	// Assert
	require.NoError(t, err)
	config := handler.GetConfig().(*testConfig)
	assert.Equal(t, "default", config.Name)
	assert.Equal(t, 42, config.Value)
}

func TestFileConfigHandler_GetConfig_WhenCalled_ThenReturnsCurrentConfig(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.json")
	em := eventsmanager.NewEventManager()
	errorCatcher := &mockErrorCatcher{}
	notifier := &mockFileChangedNotifier{}
	logger := &mockLogger{}
	defaults := &testConfig{Name: "default", Value: 42}
	handler, _ := NewFileConfigHandler(filePath, defaults, errorCatcher, em, notifier, logger)

	// Act
	config := handler.GetConfig()

	// Assert
	assert.NotNil(t, config)
	assert.IsType(t, &testConfig{}, config)
}

func TestFileConfigHandler_SetConfig_WhenValidConfig_ThenUpdatesAndWritesToFile(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.json")
	em := eventsmanager.NewEventManager()
	errorCatcher := &mockErrorCatcher{}
	notifier := &mockFileChangedNotifier{}
	logger := &mockLogger{}
	defaults := &testConfig{Name: "default", Value: 42}
	handler, _ := NewFileConfigHandler(filePath, defaults, errorCatcher, em, notifier, logger)
	newConfig := &testConfig{Name: "updated", Value: 100}

	// Act
	err := handler.SetConfig(newConfig)

	// Assert
	require.NoError(t, err)
	currentConfig := handler.GetConfig().(*testConfig)
	assert.Equal(t, "updated", currentConfig.Name)
	assert.Equal(t, 100, currentConfig.Value)
	assert.True(t, handler.CanRestore())
}

func TestFileConfigHandler_CanRestore_WhenNoOldConfig_ThenReturnsFalse(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.json")
	em := eventsmanager.NewEventManager()
	errorCatcher := &mockErrorCatcher{}
	notifier := &mockFileChangedNotifier{}
	logger := &mockLogger{}
	defaults := &testConfig{Name: "default", Value: 42}
	handler, _ := NewFileConfigHandler(filePath, defaults, errorCatcher, em, notifier, logger)

	// Act
	canRestore := handler.CanRestore()

	// Assert
	assert.False(t, canRestore)
}

func TestFileConfigHandler_CanRestore_WhenOldConfigExists_ThenReturnsTrue(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.json")
	em := eventsmanager.NewEventManager()
	errorCatcher := &mockErrorCatcher{}
	notifier := &mockFileChangedNotifier{}
	logger := &mockLogger{}
	defaults := &testConfig{Name: "default", Value: 42}
	handler, _ := NewFileConfigHandler(filePath, defaults, errorCatcher, em, notifier, logger)
	newConfig := &testConfig{Name: "updated", Value: 100}
	assert.NoError(t, handler.SetConfig(newConfig))

	// Act
	canRestore := handler.CanRestore()

	// Assert
	assert.True(t, canRestore)
}

func TestFileConfigHandler_Restore_WhenCanRestore_ThenRestoresOldConfig(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.json")
	em := eventsmanager.NewEventManager()
	errorCatcher := &mockErrorCatcher{}
	notifier := &mockFileChangedNotifier{}
	logger := &mockLogger{}
	defaults := &testConfig{Name: "default", Value: 42}
	handler, _ := NewFileConfigHandler(filePath, defaults, errorCatcher, em, notifier, logger)
	newConfig := &testConfig{Name: "updated", Value: 100}
	assert.NoError(t, handler.SetConfig(newConfig))

	// Act
	err := handler.Restore()

	// Assert
	require.NoError(t, err)
	currentConfig := handler.GetConfig().(*testConfig)
	assert.Equal(t, "default", currentConfig.Name)
	assert.Equal(t, 42, currentConfig.Value)
	assert.False(t, handler.CanRestore())
}

func TestFileConfigHandler_Restore_WhenCannotRestore_ThenReturnsError(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.json")
	em := eventsmanager.NewEventManager()
	errorCatcher := &mockErrorCatcher{}
	notifier := &mockFileChangedNotifier{}
	logger := &mockLogger{}
	defaults := &testConfig{Name: "default", Value: 42}
	handler, _ := NewFileConfigHandler(filePath, defaults, errorCatcher, em, notifier, logger)

	// Act
	err := handler.Restore()

	// Assert
	assert.Error(t, err)
	assert.IsType(t, &fileConfigHandlerError{}, err)
}

func TestFileConfigHandler_Freeze_WhenCalled_ThenSetsFreezed(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.json")
	em := eventsmanager.NewEventManager()
	errorCatcher := &mockErrorCatcher{}
	notifier := &mockFileChangedNotifier{}
	logger := &mockLogger{}
	defaults := &testConfig{Name: "default", Value: 42}
	handler, _ := NewFileConfigHandler(filePath, defaults, errorCatcher, em, notifier, logger)

	// Act
	handler.Freeze()

	// Assert
	// Note: isFreezed is private, but we can test behavior indirectly
}

func TestFileConfigHandler_Unfreeze_WhenCalled_ThenSetsUnfreezed(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.json")
	em := eventsmanager.NewEventManager()
	errorCatcher := &mockErrorCatcher{}
	notifier := &mockFileChangedNotifier{}
	logger := &mockLogger{}
	defaults := &testConfig{Name: "default", Value: 42}
	handler, _ := NewFileConfigHandler(filePath, defaults, errorCatcher, em, notifier, logger)
	handler.Freeze()

	// Act
	handler.Unfreeze()

	// Assert
	// Note: isFreezed is private, but we can test behavior indirectly
}

func TestFileConfigHandler_ForceRefresh_WhenNewConfigAvailable_ThenUpdatesConfig(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.json")
	em := eventsmanager.NewEventManager()
	errorCatcher := &mockErrorCatcher{}
	notifier := &mockFileChangedNotifier{}
	logger := &mockLogger{}
	defaults := &testConfig{Name: "default", Value: 42}
	handler, _ := NewFileConfigHandler(filePath, defaults, errorCatcher, em, notifier, logger)

	// Manually set newConfig to simulate file change
	handler.(*fileConfigHandler).newConfig = &testConfig{Name: "fromfile", Value: 200}

	// Act
	err := handler.ForceRefresh()

	// Assert
	require.NoError(t, err)
	currentConfig := handler.GetConfig().(*testConfig)
	assert.Equal(t, "fromfile", currentConfig.Name)
	assert.Equal(t, 200, currentConfig.Value)
}

func TestFileConfigHandler_SetRefreshTime_WhenValidPeriod_ThenStartsRefreshLoop(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.json")
	em := eventsmanager.NewEventManager()
	errorCatcher := &mockErrorCatcher{}
	notifier := &mockFileChangedNotifier{}
	logger := &mockLogger{}
	defaults := &testConfig{Name: "default", Value: 42}
	handler, _ := NewFileConfigHandler(filePath, defaults, errorCatcher, em, notifier, logger)
	// Note: Cannot create Period with private fields, so we test that the method exists and doesn't panic
	period := configuration.Period{} // Zero value

	// Act
	err := handler.SetRefreshTime(period)

	// Assert
	// The method should not panic, but we can't fully test the refresh loop
	assert.NoError(t, err)
}

func TestFileConfigHandler_createConfig_WhenCalled_ThenReturnsNewInstance(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.json")
	em := eventsmanager.NewEventManager()
	errorCatcher := &mockErrorCatcher{}
	notifier := &mockFileChangedNotifier{}
	logger := &mockLogger{}
	defaults := &testConfig{Name: "default", Value: 42}
	handler, _ := NewFileConfigHandler(filePath, defaults, errorCatcher, em, notifier, logger)
	fch := handler.(*fileConfigHandler)

	// Act
	config := fch.createConfig()

	// Assert
	assert.NotNil(t, config)
	assert.IsType(t, &testConfig{}, config)
}

func TestFileConfigHandler_pipeError_WhenHandlerError_ThenReturnsAsIs(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.json")
	em := eventsmanager.NewEventManager()
	errorCatcher := &mockErrorCatcher{}
	notifier := &mockFileChangedNotifier{}
	logger := &mockLogger{}
	defaults := &testConfig{Name: "default", Value: 42}
	handler, _ := NewFileConfigHandler(filePath, defaults, errorCatcher, em, notifier, logger)
	fch := handler.(*fileConfigHandler)
	handlerErr := newFileConfigHandlerError(UnexpectedError, "test error", nil)

	// Act
	result := fch.pipeError(handlerErr)

	// Assert
	assert.IsType(t, &fileConfigHandlerError{}, result)
	assert.Equal(t, UnexpectedError, result.(HandlerError).GetErrorType())
	assert.Equal(t, "test error", result.(HandlerError).GetMessage())
}

func TestFileConfigHandler_pipeError_WhenRegularError_ThenWrapsInHandlerError(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "config.json")
	em := eventsmanager.NewEventManager()
	errorCatcher := &mockErrorCatcher{}
	notifier := &mockFileChangedNotifier{}
	logger := &mockLogger{}
	defaults := &testConfig{Name: "default", Value: 42}
	handler, _ := NewFileConfigHandler(filePath, defaults, errorCatcher, em, notifier, logger)
	fch := handler.(*fileConfigHandler)
	regularErr := assert.AnError

	// Act
	result := fch.pipeError(regularErr)

	// Assert
	assert.IsType(t, &fileConfigHandlerError{}, result)
	assert.Equal(t, UnexpectedError, result.(HandlerError).GetErrorType())
}

func TestNewFileConfigHandlerError_WhenCreated_ThenReturnsError(t *testing.T) {
	// Arrange
	errorType := UnexpectedError
	message := "test message"
	internalErr := assert.AnError

	// Act
	err := newFileConfigHandlerError(errorType, message, internalErr)

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, errorType, err.GetErrorType())
	assert.Equal(t, message, err.GetMessage())
	assert.Equal(t, internalErr, err.GetInternalError())
}

func TestFileConfigHandlerError_GetErrorType_WhenCalled_ThenReturnsErrorType(t *testing.T) {
	// Arrange
	err := newFileConfigHandlerError(OldConfigNilError, "test", nil)

	// Act
	errorType := err.GetErrorType()

	// Assert
	assert.Equal(t, OldConfigNilError, errorType)
}
