package configuration_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/janmbaco/go-infrastructure/v2/configuration"
)

type snapshotSettings struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
}

var errNameRequired = errors.New("name is required")

func validateSnapshot(settings snapshotSettings) error {
	if settings.Name == "" {
		return errNameRequired
	}
	return nil
}

func writeSnapshot(test *testing.T, content string) string {
	test.Helper()
	path := filepath.Join(test.TempDir(), "settings.json")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		test.Fatal(err)
	}
	return path
}

func TestReadSnapshotReturnsIndependentValidatedValues(test *testing.T) {
	// Arrange
	path := writeSnapshot(test, `{"name":"local","labels":{"mode":"original"}}`)
	first, err := configuration.ReadSnapshot(path, validateSnapshot)
	if err != nil {
		test.Fatal(err)
	}
	// Act
	first.Labels["mode"] = "changed"
	second, err := configuration.ReadSnapshot(path, validateSnapshot)
	// Assert
	if err != nil || second.Labels["mode"] != "original" {
		test.Fatal("snapshots share mutable configuration", second, err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != `{"name":"local","labels":{"mode":"original"}}` {
		test.Fatal("reader changed file", err)
	}
}

func TestReadSnapshotDoesNotCreateMissingFiles(test *testing.T) {
	// Arrange
	path := filepath.Join(test.TempDir(), "missing.json")
	// Act
	_, err := configuration.ReadSnapshot(path, validateSnapshot)
	// Assert
	if !errors.Is(err, fs.ErrNotExist) {
		test.Fatal("lost missing-file cause", err)
	}
	if _, err := os.Stat(path); !errors.Is(err, fs.ErrNotExist) {
		test.Fatal("reader created a file", err)
	}
}

func TestReadSnapshotPreservesValidationCause(test *testing.T) {
	// Arrange
	path := writeSnapshot(test, `{"name":""}`)
	// Act
	result, err := configuration.ReadSnapshot(path, validateSnapshot)
	// Assert
	if !errors.Is(err, errNameRequired) || result.Name != "" || result.Labels != nil {
		test.Fatal("invalid value escaped validation", result, err)
	}
}

func TestReadSnapshotRejectsUnknownTrailingOversizedAndMalformedDocuments(test *testing.T) {
	for name, content := range map[string]string{"unknown": `{"name":"local","extra":1}`, "trailing": `{"name":"local"} {}`, "malformed": `{"name":`, "oversized": strings.Repeat(" ", configuration.MaxSnapshotBytes+1)} {
		test.Run(name, func(test *testing.T) {
			// Arrange
			path := writeSnapshot(test, content)
			// Act
			_, err := configuration.ReadSnapshot(path, validateSnapshot)
			// Assert
			if err == nil {
				test.Fatal("invalid document accepted")
			}
		})
	}
}

func TestReadSnapshotRequiresValidation(test *testing.T) {
	// Arrange
	path := writeSnapshot(test, `{"name":"local"}`)
	// Act
	_, err := configuration.ReadSnapshot[snapshotSettings](path, nil)
	// Assert
	if err == nil {
		test.Fatal("missing validator accepted")
	}
}
