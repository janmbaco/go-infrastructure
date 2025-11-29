package disk

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/janmbaco/go-infrastructure/v2/eventsmanager"
	"github.com/janmbaco/go-infrastructure/v2/logs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExistsPath_WhenPathExists_ThenReturnsTrue(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(filePath, []byte("test"), 0644)
	require.NoError(t, err)

	// Act
	exists := ExistsPath(filePath)

	// Assert
	assert.True(t, exists)
}

func TestExistsPath_WhenPathDoesNotExist_ThenReturnsFalse(t *testing.T) {
	// Arrange
	filePath := "/nonexistent/path"

	// Act
	exists := ExistsPath(filePath)

	// Assert
	assert.False(t, exists)
}

func TestCreateFile_WhenValidPath_ThenCreatesFileWithContent(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")
	content := []byte("hello world")

	// Act
	err := CreateFile(filePath, content)

	// Assert
	require.NoError(t, err)
	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, content, data)
}

func TestCreateFile_WhenInvalidPath_ThenReturnsError(t *testing.T) {
	// Arrange
	filePath := "/invalid/path/test.txt"
	content := []byte("test")

	// Act
	err := CreateFile(filePath, content)

	// Assert
	assert.Error(t, err)
}

func TestCopy_WhenValidFiles_ThenCopiesFile(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcPath := filepath.Join(tempDir, "src.txt")
	dstPath := filepath.Join(tempDir, "dst.txt")
	content := []byte("copy me")
	err := os.WriteFile(srcPath, content, 0644)
	require.NoError(t, err)

	// Act
	err = Copy(srcPath, dstPath)

	// Assert
	require.NoError(t, err)
	data, err := os.ReadFile(dstPath)
	require.NoError(t, err)
	assert.Equal(t, content, data)
}

func TestCopy_WhenSrcDoesNotExist_ThenReturnsError(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	srcPath := "/nonexistent"
	dstPath := filepath.Join(tempDir, "dst.txt")

	// Act
	err := Copy(srcPath, dstPath)

	// Assert
	assert.Error(t, err)
}

func TestCopy_WhenSrcIsDirectory_ThenReturnsError(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	require.NoError(t, err)
	dstPath := filepath.Join(tempDir, "dst.txt")

	// Act
	err = Copy(subDir, dstPath)

	// Assert
	assert.Error(t, err)
}

func TestDeleteFile_WhenFileExists_ThenDeletesFile(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(filePath, []byte("test"), 0644)
	require.NoError(t, err)

	// Act
	err = DeleteFile(filePath)

	// Assert
	require.NoError(t, err)
	exists := ExistsPath(filePath)
	assert.False(t, exists)
}

func TestDeleteFile_WhenFileDoesNotExist_ThenReturnsError(t *testing.T) {
	// Arrange
	filePath := "/nonexistent"

	// Act
	err := DeleteFile(filePath)

	// Assert
	assert.Error(t, err)
}

func TestNewFileChangedNotifier_WhenValidPath_ThenReturnsNotifier(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(filePath, []byte("test"), 0644)
	require.NoError(t, err)
	em := eventsmanager.NewEventManager()
	logger := logs.NewLogger()

	// Act
	notifier, err := NewFileChangedNotifier(filePath, em, logger)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, notifier)
	assert.IsType(t, &fileChangedNotifier{}, notifier)
}

func TestFileChangedNotifier_Subscribe_WhenFirstSubscribe_ThenStartsWatching(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(filePath, []byte("test"), 0644)
	require.NoError(t, err)
	em := eventsmanager.NewEventManager()
	logger := logs.NewLogger()
	notifier, _ := NewFileChangedNotifier(filePath, em, logger)
	fcn := notifier.(*fileChangedNotifier)

	// Act
	err = notifier.Subscribe(func() {})

	// Assert
	require.NoError(t, err)
	assert.True(t, fcn.isWatchingFile)
}
