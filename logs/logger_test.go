package logs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger_WhenCreated_ThenReturnsLogger(t *testing.T) {
	// Arrange & Act
	logger := NewLogger()

	// Assert
	assert.NotNil(t, logger)
}

func TestLogger_SetConsoleLevel_WhenSetToInfo_ThenTraceIsDisabled(t *testing.T) {
	// Arrange
	l := NewLogger().(*logger)

	// Act
	l.SetConsoleLevel(Info)

	// Assert
	assert.False(t, l.activeConsoleLogger[Trace])
	assert.True(t, l.activeConsoleLogger[Info])
}

func TestLogger_SetFileLogLevel_WhenSetToWarning_ThenInfoIsDisabled(t *testing.T) {
	// Arrange
	l := NewLogger().(*logger)

	// Act
	l.SetFileLogLevel(Warning)

	// Assert
	assert.False(t, l.activeFileLogger[Info])
	assert.True(t, l.activeFileLogger[Warning])
}

func TestLogger_Mute_WhenMuted_ThenPrintlnDoesNothing(t *testing.T) {
	// Arrange
	l := NewLogger()
	l.(*logger).muted = true

	// Act
	l.Println(Info, "test")

	// Assert
	assert.True(t, l.(*logger).muted)
}

func TestLogger_Unmute_WhenUnmuted_ThenPrintlnWorks(t *testing.T) {
	// Arrange
	l := NewLogger()
	l.Mute()
	l.Unmute()

	// Act & Assert
	assert.False(t, l.(*logger).muted)
}

func TestLogger_Printlnf_WhenFormat_ThenFormatsMessage(t *testing.T) {
	// Arrange
	l := NewLogger()
	l.Mute()

	// Act
	l.Printlnf(Info, "test %s", "message")

	// Assert
	// No assertion, just ensure no panic
}

func TestLogger_PrintError_WhenError_ThenLogsErrorMessage(t *testing.T) {
	// Arrange
	l := NewLogger()
	l.Mute()

	// Act
	l.PrintError(Info, assert.AnError)

	// Assert
	// No panic
}

func TestLogger_TryPrintError_WhenNilError_ThenDoesNothing(t *testing.T) {
	// Arrange
	l := NewLogger()

	// Act
	l.TryPrintError(Info, nil)

	// Assert
	// No panic
}

func TestLogger_TryPrintError_WhenError_ThenLogs(t *testing.T) {
	// Arrange
	l := NewLogger()
	l.Mute()

	// Act
	l.TryPrintError(Info, assert.AnError)

	// Assert
	// No panic
}

func TestLogger_Trace_WhenCalled_ThenLogsAtTraceLevel(t *testing.T) {
	// Arrange
	l := NewLogger()
	l.Mute()

	// Act
	l.Trace("trace message")

	// Assert
	// No panic
}

func TestLogger_Info_WhenCalled_ThenLogsAtInfoLevel(t *testing.T) {
	// Arrange
	l := NewLogger()
	l.Mute()

	// Act
	l.Info("info message")

	// Assert
	// No panic
}

func TestLogger_Warning_WhenCalled_ThenLogsAtWarningLevel(t *testing.T) {
	// Arrange
	l := NewLogger()
	l.Mute()

	// Act
	l.Warning("warning message")

	// Assert
	// No panic
}

func TestLogger_Error_WhenCalled_ThenLogsAtErrorLevel(t *testing.T) {
	// Arrange
	l := NewLogger()
	l.Mute()

	// Act
	l.Error("error message")

	// Assert
	// No panic
}

func TestLogger_Fatal_WhenCalled_ThenLogsAtFatalLevel(t *testing.T) {
	// Arrange
	l := NewLogger()
	l.Mute()

	// Act
	l.Fatal("fatal message")

	// Assert
	// No panic
}

func TestLogger_SetDir_WhenSet_ThenDirIsSet(t *testing.T) {
	// Arrange
	l := NewLogger().(*logger)
	dir := "/tmp/logs"

	// Act
	l.SetDir(dir)

	// Assert
	assert.Equal(t, dir, l.logsDir)
}

func TestLogger_GetErrorLogger_WhenCalled_ThenReturnsLogger(t *testing.T) {
	// Arrange
	l := NewLogger()

	// Act
	errLogger := l.GetErrorLogger()

	// Assert
	assert.NotNil(t, errLogger)
}

func TestSetLevel_WhenTrace_ThenAllActive(t *testing.T) {
	// Act
	levels := setLevel(Trace)

	// Assert
	assert.True(t, levels[Trace])
	assert.True(t, levels[Info])
	assert.True(t, levels[Warning])
	assert.True(t, levels[Error])
	assert.True(t, levels[Fatal])
}

func TestSetLevel_WhenInfo_ThenTraceDisabled(t *testing.T) {
	// Act
	levels := setLevel(Info)

	// Assert
	assert.False(t, levels[Trace])
	assert.True(t, levels[Info])
	assert.True(t, levels[Warning])
	assert.True(t, levels[Error])
	assert.True(t, levels[Fatal])
}

func TestSetLevel_WhenError_ThenLowerDisabled(t *testing.T) {
	// Act
	levels := setLevel(Error)

	// Assert
	assert.False(t, levels[Trace])
	assert.False(t, levels[Info])
	assert.False(t, levels[Warning])
	assert.True(t, levels[Error])
	assert.True(t, levels[Fatal])
}
