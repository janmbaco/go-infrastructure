package ioc

import (
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	persistence "github.com/janmbaco/go-infrastructure/v2/persistence"
	"github.com/janmbaco/go-infrastructure/v2/persistence/dialectors"
)

// DialectorsModule implements Module for database dialectors
type DialectorsModule struct{}

// NewDialectorsModule creates a new dialectors module
func NewDialectorsModule() *DialectorsModule {
	return &DialectorsModule{}
}

// RegisterServices registers all database dialectors
func (m *DialectorsModule) RegisterServices(register dependencyinjection.Register) error {
	mysqlKey, _ := persistence.MySQL.ToString() //nolint:errcheck // ToString called with known constants that cannot fail
	dependencyinjection.RegisterSingletonTenant[persistence.DialectorGetter](
		register,
		mysqlKey,
		dialectors.NewMysqlDialectorGetter,
	)

	postgresKey, _ := persistence.Postgres.ToString() //nolint:errcheck // ToString called with known constants that cannot fail
	dependencyinjection.RegisterSingletonTenant[persistence.DialectorGetter](
		register,
		postgresKey,
		dialectors.NewPostgresDialectorGetter,
	)

	sqliteKey, _ := persistence.Sqlite.ToString() //nolint:errcheck // ToString called with known constants that cannot fail
	dependencyinjection.RegisterSingletonTenant[persistence.DialectorGetter](
		register,
		sqliteKey,
		dialectors.NewSqliteDialectorGetter,
	)

	sqlServerKey, _ := persistence.SQLServer.ToString() //nolint:errcheck // ToString called with known constants that cannot fail
	dependencyinjection.RegisterSingletonTenant[persistence.DialectorGetter](
		register,
		sqlServerKey,
		dialectors.NewSqlServerDialectorGetter,
	)

	return nil
}
