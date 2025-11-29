package ioc

import (
	"fmt"

	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	persistence "github.com/janmbaco/go-infrastructure/v2/persistence"
	"github.com/janmbaco/go-infrastructure/v2/persistence/dialectors"
	"gorm.io/gorm"
)

// DatabaseConfig represents database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Engine   persistence.DbEngine
}

// PostgresModule implements Module for PostgreSQL database services
type PostgresModule struct {
	info *persistence.DatabaseInfo
}

// NewPostgresModule creates a new PostgreSQL module with database info
func NewPostgresModule(info *persistence.DatabaseInfo) *PostgresModule {
	return &PostgresModule{info: info}
}

// RegisterServices registers PostgreSQL database services
func (m *PostgresModule) RegisterServices(register dependencyinjection.Register) error {
	// Register database connection
	register.AsSingleton(new(*gorm.DB), func() (*gorm.DB, error) {
		dialectorGetter := dialectors.NewPostgresDialectorGetter()
		dialector, err := dialectorGetter.Get(m.info)
		if err != nil {
			return nil, err
		}
		return gorm.Open(dialector, &gorm.Config{})
	}, nil)

	return nil
}

// MysqlModule implements Module for MySQL database services
type MysqlModule struct {
	info *persistence.DatabaseInfo
}

// NewMysqlModule creates a new MySQL module with database info
func NewMysqlModule(info *persistence.DatabaseInfo) *MysqlModule {
	return &MysqlModule{info: info}
}

// RegisterServices registers MySQL database services
func (m *MysqlModule) RegisterServices(register dependencyinjection.Register) error {
	// Register database connection
	register.AsSingleton(new(*gorm.DB), func() (*gorm.DB, error) {
		dialectorGetter := dialectors.NewMysqlDialectorGetter()
		dialector, err := dialectorGetter.Get(m.info)
		if err != nil {
			return nil, err
		}
		return gorm.Open(dialector, &gorm.Config{})
	}, nil)

	return nil
}

// SqliteModule implements Module for SQLite database services
type SqliteModule struct {
	info *persistence.DatabaseInfo
}

// NewSqliteModule creates a new SQLite module with database info
func NewSqliteModule(info *persistence.DatabaseInfo) *SqliteModule {
	return &SqliteModule{info: info}
}

// RegisterServices registers SQLite database services
func (m *SqliteModule) RegisterServices(register dependencyinjection.Register) error {
	// Register database connection
	register.AsSingleton(new(*gorm.DB), func() (*gorm.DB, error) {
		dialectorGetter := dialectors.NewSqliteDialectorGetter()
		dialector, err := dialectorGetter.Get(m.info)
		if err != nil {
			return nil, err
		}
		return gorm.Open(dialector, &gorm.Config{})
	}, nil)

	return nil
}

// SqlServerModule implements Module for SQL Server database services
type SqlServerModule struct {
	info *persistence.DatabaseInfo
}

// NewSqlServerModule creates a new SQL Server module with database info
func NewSqlServerModule(info *persistence.DatabaseInfo) *SqlServerModule {
	return &SqlServerModule{info: info}
}

// RegisterServices registers SQL Server database services
func (m *SqlServerModule) RegisterServices(register dependencyinjection.Register) error {
	// Register database connection
	register.AsSingleton(new(*gorm.DB), func() (*gorm.DB, error) {
		dialectorGetter := dialectors.NewSqlServerDialectorGetter()
		dialector, err := dialectorGetter.Get(m.info)
		if err != nil {
			return nil, err
		}
		return gorm.Open(dialector, &gorm.Config{})
	}, nil)

	return nil
}

// ConfigureDatabaseModule creates the appropriate database module based on configuration
func ConfigureDatabaseModule(host, port, user, password, dbName string, engine persistence.DbEngine) dependencyinjection.Module {
	dbInfo := &persistence.DatabaseInfo{
		Host:         host,
		Port:         port,
		Name:         dbName,
		UserName:     user,
		UserPassword: password,
		Engine:       engine,
	}

	switch engine {
	case persistence.Postgres:
		return NewPostgresModule(dbInfo)
	case persistence.MySQL:
		return NewMysqlModule(dbInfo)
	case persistence.SQLServer:
		return NewSqlServerModule(dbInfo)
	case persistence.Sqlite:
		return NewSqliteModule(dbInfo)
	default:
		panic(fmt.Sprintf("unsupported database engine: %v", engine))
	}
}
