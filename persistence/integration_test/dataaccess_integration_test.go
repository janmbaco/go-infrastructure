//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/v2/persistence"
	"github.com/janmbaco/go-infrastructure/v2/persistence/dataaccess"
	persistenceioc "github.com/janmbaco/go-infrastructure/v2/persistence/ioc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// TestUser modelo para testing de integración
type TestUser struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:100;not null"`
	Email     string `gorm:"size:100;unique;not null"`
	Age       int    `gorm:"not null"`
	Active    bool   `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TestProfile modelo relacionado para testing de asociaciones
type TestProfile struct {
	ID     uint     `gorm:"primaryKey"`
	UserID uint     `gorm:"not null"`
	Bio    string   `gorm:"size:500"`
	Avatar string   `gorm:"size:255"`
	User   TestUser `gorm:"foreignKey:UserID"`
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Engine   persistence.DbEngine
}

func getDatabaseConfigs() []DatabaseConfig {
	return []DatabaseConfig{
		{
			Host:     getEnvOrDefault("POSTGRES_HOST", "localhost"),
			Port:     getEnvOrDefault("POSTGRES_PORT", "5432"),
			User:     getEnvOrDefault("POSTGRES_USER", "testuser"),
			Password: getEnvOrDefault("POSTGRES_PASSWORD", "testpass"),
			DBName:   getEnvOrDefault("POSTGRES_DB", "testdb"),
			Engine:   persistence.Postgres,
		},
		{
			Host:     getEnvOrDefault("MYSQL_HOST", "localhost"),
			Port:     getEnvOrDefault("MYSQL_PORT", "3306"),
			User:     getEnvOrDefault("MYSQL_USER", "testuser"),
			Password: getEnvOrDefault("MYSQL_PASSWORD", "testpass"),
			DBName:   getEnvOrDefault("MYSQL_DB", "testdb"),
			Engine:   persistence.MySQL,
		},
		{
			Host:     getEnvOrDefault("SQLSERVER_HOST", "localhost"),
			Port:     getEnvOrDefault("SQLSERVER_PORT", "1433"),
			User:     getEnvOrDefault("SQLSERVER_USER", "sa"),
			Password: getEnvOrDefault("SQLSERVER_PASSWORD", "StrongPass123!"),
			DBName:   getEnvOrDefault("SQLSERVER_DB", "master"),
			Engine:   persistence.SQLServer,
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func createDatabaseConnection(config DatabaseConfig) (*gorm.DB, error) {
	// Create container with appropriate database module
	container := dependencyinjection.NewBuilder().
		AddModule(persistenceioc.ConfigureDatabaseModule(config.Host, config.Port, config.User, config.Password, config.DBName, config.Engine)).
		MustBuild()

	// Get database connection from container
	resolver := container.Resolver()
	db, ok := resolver.Type(new(*gorm.DB), nil).(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("failed to resolve database connection")
	}

	return db, nil
}

func waitForDatabase(config DatabaseConfig, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for database %v to be ready", config.Engine)
		case <-ticker.C:
			db, err := createDatabaseConnection(config)
			if err == nil {
				sqlDB, _ := db.DB()
				if err := sqlDB.Ping(); err == nil {
					sqlDB.Close()
					return nil
				}
				sqlDB.Close()
			}
		}
	}
}

func TestDataAccessIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	databaseConfigs := getDatabaseConfigs()

	for _, config := range databaseConfigs {
		t.Run(fmt.Sprintf("Database_%v", config.Engine), func(t *testing.T) {
			testDatabaseOperations(t, config)
		})
	}
}

func testDatabaseOperations(t *testing.T, config DatabaseConfig) {
	// Wait for database to be ready
	err := waitForDatabase(config, 60*time.Second)
	require.NoError(t, err, "Database should be ready")

	// Create database connection
	db, err := createDatabaseConnection(config)
	require.NoError(t, err, "Should connect to database")
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Auto-migrate tables
	err = db.AutoMigrate(&TestUser{}, &TestProfile{})
	require.NoError(t, err, "Should migrate tables")

	// Clean up tables before test
	db.Exec("DELETE FROM test_profiles")
	db.Exec("DELETE FROM test_users")

	// Test with generic functions
	testCRUDOperations(t, db)
	testAssociations(t, db)
	testErrorHandling(t, db)
}

func testCRUDOperations(t *testing.T, db *gorm.DB) {
	t.Run("CRUD_Operations", func(t *testing.T) {
		// Create DataAccess
		dataAccess := dataaccess.NewTypedDataAccess[TestUser](db)

		// Test INSERT
		user := &TestUser{
			Name:   "John Doe",
			Email:  "john@example.com",
			Age:    30,
			Active: true,
		}

		err := dataaccess.InsertRow(dataAccess, user)
		assert.NoError(t, err, "Should insert user")
		assert.NotZero(t, user.ID, "User ID should be set after insert")

		// Test SELECT by ID
		users, err := dataaccess.SelectRows(dataAccess, &TestUser{ID: user.ID})
		assert.NoError(t, err, "Should select user")
		assert.Len(t, users, 1, "Should return one user")
		assert.Equal(t, user.Name, users[0].Name, "User name should match")
		assert.Equal(t, user.Email, users[0].Email, "User email should match")

		// Test UPDATE
		user.Name = "John Smith"
		err = dataaccess.UpdateRow(dataAccess, &TestUser{ID: user.ID}, user)
		assert.NoError(t, err, "Should update user")

		// Verify UPDATE
		users, err = dataaccess.SelectRows(dataAccess, &TestUser{ID: user.ID})
		assert.NoError(t, err, "Should select updated user")
		assert.Equal(t, "John Smith", users[0].Name, "User name should be updated")

		// Test SELECT with conditions
		users, err = dataaccess.SelectRows(dataAccess, &TestUser{Email: "john@example.com"})
		assert.NoError(t, err, "Should select user by email")
		assert.Len(t, users, 1, "Should return one user")

		// Test DELETE
		err = dataaccess.DeleteRows(dataAccess, &TestUser{ID: user.ID})
		assert.NoError(t, err, "Should delete user")

		// Verify DELETE
		users, err = dataaccess.SelectRows(dataAccess, &TestUser{ID: user.ID})
		assert.NoError(t, err, "Should not error after delete")
		assert.Len(t, users, 0, "Should return no users after delete")
	})
}

func testAssociations(t *testing.T, db *gorm.DB) {
	t.Run("Associations", func(t *testing.T) {
		// Create user
		userDataAccess := dataaccess.NewTypedDataAccess[TestUser](db)
		user := &TestUser{
			Name:   "Jane Doe",
			Email:  "jane@example.com",
			Age:    25,
			Active: true,
		}

		err := dataaccess.InsertRow(userDataAccess, user)
		require.NoError(t, err, "Should insert user")

		// Create profile
		profileDataAccess := dataaccess.NewTypedDataAccess[TestProfile](db)
		profile := &TestProfile{
			UserID: user.ID,
			Bio:    "Software Developer",
			Avatar: "avatar.jpg",
		}

		err = dataaccess.InsertRow(profileDataAccess, profile)
		assert.NoError(t, err, "Should insert profile")

		// Test SELECT with preload
		profiles, err := dataaccess.SelectRows[TestProfile](profileDataAccess, &TestProfile{}, "User")
		assert.NoError(t, err, "Should select profiles with user preload")
		assert.Len(t, profiles, 1, "Should return one profile")
		assert.Equal(t, user.ID, profiles[0].UserID, "Profile should reference correct user")
		assert.Equal(t, user.Name, profiles[0].User.Name, "Preloaded user should have correct name")

		// Test DELETE with associations
		err = dataaccess.DeleteRows(profileDataAccess, &TestProfile{UserID: user.ID}, "User")
		assert.NoError(t, err, "Should delete profile with associations")

		// Verify profile is deleted
		profiles, err = dataaccess.SelectRows[TestProfile](profileDataAccess, &TestProfile{UserID: user.ID})
		assert.NoError(t, err, "Should not error after delete")
		assert.Len(t, profiles, 0, "Should return no profiles after delete")

		// Clean up user
		dataaccess.DeleteRows(userDataAccess, &TestUser{ID: user.ID})
	})
}

func testErrorHandling(t *testing.T, db *gorm.DB) {
	t.Run("Error_Handling", func(t *testing.T) {
		dataAccess := dataaccess.NewTypedDataAccess[TestUser](db)

		// Test duplicate email (unique constraint)
		user1 := &TestUser{
			Name:   "User 1",
			Email:  "duplicate@example.com",
			Age:    20,
			Active: true,
		}

		err := dataaccess.InsertRow(dataAccess, user1)
		assert.NoError(t, err, "Should insert first user")

		user2 := &TestUser{
			Name:   "User 2",
			Email:  "duplicate@example.com", // Same email
			Age:    21,
			Active: true,
		}

		err = dataaccess.InsertRow(dataAccess, user2)
		assert.Error(t, err, "Should error on duplicate email")

		// Clean up
		dataaccess.DeleteRows(dataAccess, &TestUser{Email: "duplicate@example.com"})
	})
}

func TestAdvancedQueriesWithDBMethod(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	databaseConfigs := getDatabaseConfigs()

	for _, config := range databaseConfigs {
		t.Run(fmt.Sprintf("AdvancedQueries_%v", config.Engine), func(t *testing.T) {
			testAdvancedQueries(t, config)
		})
	}
}

func testAdvancedQueries(t *testing.T, config DatabaseConfig) {
	// Wait for database to be ready
	err := waitForDatabase(config, 60*time.Second)
	require.NoError(t, err, "Database should be ready")

	// Create database connection
	db, err := createDatabaseConnection(config)
	require.NoError(t, err, "Should connect to database")
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Auto-migrate tables
	err = db.AutoMigrate(&TestUser{}, &TestProfile{})
	require.NoError(t, err, "Should migrate tables")

	// Clean up any existing data
	db.Unscoped().Delete(&TestUser{}, "1=1")
	db.Unscoped().Delete(&TestProfile{}, "1=1")

	dataAccess := dataaccess.NewTypedDataAccess[TestUser](db)

	t.Run("RawSQLQueries", func(t *testing.T) {
		// Insert test data
		users := []*TestUser{
			{Name: "Alice", Email: "alice@test.com", Age: 25, Active: true},
			{Name: "Bob", Email: "bob@test.com", Age: 30, Active: false},
			{Name: "Charlie", Email: "charlie@test.com", Age: 35, Active: true},
			{Name: "Diana", Email: "diana@test.com", Age: 28, Active: true},
		}

		for _, user := range users {
			err := dataaccess.InsertRow(dataAccess, user)
			require.NoError(t, err, "Should insert user")
		}

		// Test raw SQL query using DB() method
		gormDB := dataAccess.DB().(*gorm.DB)

		var stats struct {
			TotalUsers  int
			ActiveUsers int
			AvgAge      float64
			OldestUser  string
		}

		err := gormDB.Raw(`
			SELECT
				COUNT(*) as total_users,
				COUNT(CASE WHEN active = true THEN 1 END) as active_users,
				AVG(age) as avg_age,
				(SELECT name FROM test_users WHERE active = true ORDER BY age DESC LIMIT 1) as oldest_user
			FROM test_users
		`).Scan(&stats).Error

		assert.NoError(t, err, "Should execute raw SQL query")
		assert.Equal(t, 4, stats.TotalUsers, "Should count all users")
		assert.Equal(t, 3, stats.ActiveUsers, "Should count active users")
		assert.Equal(t, "Charlie", stats.OldestUser, "Should find oldest active user")
		assert.True(t, stats.AvgAge > 29 && stats.AvgAge < 30, "Should calculate average age")
	})

	t.Run("ComplexJoinsAndAggregations", func(t *testing.T) {
		// Create users with profiles
		userDataAccess := dataaccess.NewTypedDataAccess[TestUser](db)
		profileDataAccess := dataaccess.NewTypedDataAccess[TestProfile](db)

		users := []*TestUser{
			{Name: "John", Email: "john@test.com", Age: 40, Active: true},
			{Name: "Jane", Email: "jane@test.com", Age: 35, Active: true},
		}

		for _, user := range users {
			err := dataaccess.InsertRow(userDataAccess, user)
			require.NoError(t, err, "Should insert user")

			// Create profile for each user
			profile := &TestProfile{
				UserID: user.ID,
				Bio:    fmt.Sprintf("Bio for %s", user.Name),
				Avatar: fmt.Sprintf("%s.jpg", user.Name),
			}
			err = dataaccess.InsertRow(profileDataAccess, profile)
			require.NoError(t, err, "Should insert profile")
		}

		// Test complex join using DB() method
		gormDB := dataAccess.DB().(*gorm.DB)

		var results []struct {
			UserName   string
			UserEmail  string
			ProfileBio string
			HasProfile bool
		}

		err := gormDB.
			Table("test_users").
			Select("test_users.name as user_name, test_users.email as user_email, test_profiles.bio as profile_bio, CASE WHEN test_profiles.id IS NOT NULL THEN true ELSE false END as has_profile").
			Joins("LEFT JOIN test_profiles ON test_users.id = test_profiles.user_id").
			Where("test_users.active = ?", true).
			Order("test_users.name").
			Scan(&results).Error

		assert.NoError(t, err, "Should execute complex join query")
		assert.Len(t, results, 2, "Should return both users")
		assert.True(t, results[0].HasProfile, "First user should have profile")
		assert.True(t, results[1].HasProfile, "Second user should have profile")
		assert.Equal(t, "Jane", results[0].UserName, "Should order by name")
		assert.Equal(t, "John", results[1].UserName, "Should order by name")
	})

	t.Run("AdvancedFilteringAndGrouping", func(t *testing.T) {
		// Insert more test data for grouping
		moreUsers := []*TestUser{
			{Name: "Eve", Email: "eve@test.com", Age: 22, Active: true},
			{Name: "Frank", Email: "frank@test.com", Age: 45, Active: false},
			{Name: "Grace", Email: "grace@test.com", Age: 33, Active: true},
			{Name: "Henry", Email: "henry@test.com", Age: 41, Active: true},
		}

		for _, user := range moreUsers {
			err := dataaccess.InsertRow(dataAccess, user)
			require.NoError(t, err, "Should insert user")
		}

		// Test advanced filtering and grouping using DB() method
		gormDB := dataAccess.DB().(*gorm.DB)

		var ageGroups []struct {
			AgeGroup string
			Count    int
			AvgAge   float64
		}

		err := gormDB.
			Table("test_users").
			Select(`
				CASE
					WHEN age < 30 THEN 'Young'
					WHEN age BETWEEN 30 AND 40 THEN 'Middle'
					ELSE 'Senior'
				END as age_group,
				COUNT(*) as count,
				AVG(age) as avg_age
			`).
			Group("CASE WHEN age < 30 THEN 'Young' WHEN age BETWEEN 30 AND 40 THEN 'Middle' ELSE 'Senior' END").
			Order("avg_age").
			Scan(&ageGroups).Error

		assert.NoError(t, err, "Should execute grouping query")
		assert.True(t, len(ageGroups) >= 2, "Should have at least 2 age groups")

		// Verify the groups contain expected data
		totalUsers := 0
		for _, group := range ageGroups {
			totalUsers += group.Count
			assert.True(t, group.AvgAge > 0, "Should have valid average age")
		}
		assert.Equal(t, 6, totalUsers, "Should account for all users")
	})

	// Clean up
	db.Unscoped().Delete(&TestUser{}, "1=1")
	db.Unscoped().Delete(&TestProfile{}, "1=1")
}
