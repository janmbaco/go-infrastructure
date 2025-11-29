package persistence

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDbEngine_ToString_WhenSQLServer_ThenReturnsCorrectString(t *testing.T) {
	// Arrange
	engine := SQLServer

	// Act
	result, err := engine.ToString()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "SqlServerDB", result)
}

func TestDbEngine_ToString_WhenPostgres_ThenReturnsCorrectString(t *testing.T) {
	// Arrange
	engine := Postgres

	// Act
	result, err := engine.ToString()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "PostgresDB", result)
}

func TestDbEngine_ToString_WhenMySQL_ThenReturnsCorrectString(t *testing.T) {
	// Arrange
	engine := MySQL

	// Act
	result, err := engine.ToString()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "MySqlDB", result)
}

func TestDbEngine_ToString_WhenSqlite_ThenReturnsCorrectString(t *testing.T) {
	// Arrange
	engine := Sqlite

	// Act
	result, err := engine.ToString()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "SqliteDB", result)
}

func TestDbEngine_ToString_WhenUnknownEngine_ThenReturnsError(t *testing.T) {
	// Arrange
	engine := DbEngine(10) // Invalid engine (higher than Sqlite = 3)

	// Act
	result, err := engine.ToString()

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "unknown database engine")
}

func TestDatabaseInfo_Struct_WhenInitialized_ThenHasCorrectFields(t *testing.T) {
	// Arrange & Act
	info := &DatabaseInfo{
		Host:         "localhost",
		Port:         "5432",
		Name:         "testdb",
		UserName:     "testuser",
		UserPassword: "testpass",
		Engine:       Postgres,
	}

	// Assert
	assert.Equal(t, "localhost", info.Host)
	assert.Equal(t, "5432", info.Port)
	assert.Equal(t, "testdb", info.Name)
	assert.Equal(t, "testuser", info.UserName)
	assert.Equal(t, "testpass", info.UserPassword)
	assert.Equal(t, Postgres, info.Engine)
}

// Mock DialectorResolver for testing
type mockDialectorResolver struct {
	shouldError bool
}

func (m *mockDialectorResolver) Resolve(info *DatabaseInfo) (gorm.Dialector, error) {
	if m.shouldError {
		return nil, assert.AnError
	}
	// Return a mock dialector - in real tests this would be more complex
	return nil, nil
}

func TestNewDB_WhenDialectorResolverFails_ThenReturnsError(t *testing.T) {
	// Arrange
	resolver := &mockDialectorResolver{shouldError: true}
	info := &DatabaseInfo{Engine: Postgres}
	config := &gorm.Config{}
	tables := []interface{}{}

	// Act
	db, err := NewDB(resolver, info, config, tables)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestNewDB_WhenValidParameters_ThenCreatesDBSuccessfully(t *testing.T) {
	// This test would require a real database connection
	// For unit testing, we would need to mock gorm.Open
	// For now, we'll skip this as it requires integration setup
	t.Skip("Requires database connection - tested in integration tests")
}
