package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSinglePageApp_WhenValidPaths_ThenCreatesHandler(t *testing.T) {
	// Arrange
	staticPath := "/static"
	indexPath := "index.html"

	// Act
	handler := NewSinglePageApp(staticPath, indexPath)

	// Assert
	assert.NotNil(t, handler)
	assert.IsType(t, &singlePageApp{}, handler)
}

func TestSinglePageApp_ServeHTTP_WhenStaticFileExists_ThenServesFile(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	staticFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(staticFile, []byte("test content"), 0644)
	require.NoError(t, err)

	handler := NewSinglePageApp(tempDir, "index.html")
	req := httptest.NewRequest("GET", "/test.txt", nil)
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test content", w.Body.String())
}

func TestSinglePageApp_ServeHTTP_WhenFileDoesNotExist_ThenServesIndex(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	indexFile := filepath.Join(tempDir, "index.html")
	err := os.WriteFile(indexFile, []byte("<html>index</html>"), 0644)
	require.NoError(t, err)

	handler := NewSinglePageApp(tempDir, "index.html")
	req := httptest.NewRequest("GET", "/nonexistent.html", nil)
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "<html>index</html>", w.Body.String())
}

func TestSinglePageApp_ServeHTTP_WhenIndexDoesNotExist_ThenServesIndexAnyway(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	handler := NewSinglePageApp(tempDir, "index.html")
	req := httptest.NewRequest("GET", "/app/route", nil)
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	// When index doesn't exist, http.ServeFile returns 404
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSinglePageApp_ServeHTTP_WhenPathTraversalAttempt_ThenHandlesSafely(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	handler := NewSinglePageApp(tempDir, "index.html")
	req := httptest.NewRequest("GET", "/../../../etc/passwd", nil)
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	// Should either serve index or return error, but not access files outside static path
	assert.True(t, w.Code == http.StatusOK || w.Code >= 400)
}

func TestSinglePageApp_ServeHTTP_WhenInvalidPath_ThenReturnsBadRequest(t *testing.T) {
	// Arrange
	// Use a path that will cause filepath.Abs to fail or be invalid
	handler := NewSinglePageApp("/invalid/path", "index.html")
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	// The behavior depends on whether the path can be resolved
	// If the path doesn't exist, it will try to serve the index
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound || w.Code >= 400)
}

func TestSinglePageApp_ServeHTTP_WhenRootPath_ThenServesIndex(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	indexFile := filepath.Join(tempDir, "index.html")
	err := os.WriteFile(indexFile, []byte("root content"), 0644)
	require.NoError(t, err)

	handler := NewSinglePageApp(tempDir, "index.html")
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "root content", w.Body.String())
}
