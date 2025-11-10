package ioc

import (
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/v2/persistence/orm_base"
	"github.com/janmbaco/go-infrastructure/v2/persistence/orm_base/dialectors"
)

// DialectorsModule implements Module for database dialectors
type DialectorsModule struct{}

// NewDialectorsModule creates a new dialectors module
func NewDialectorsModule() *DialectorsModule {
	return &DialectorsModule{}
}

// RegisterServices registers all database dialectors
func (m *DialectorsModule) RegisterServices(register dependencyinjection.Register) error {
	mysqlKey, _ := orm_base.MySQL.ToString() //nolint:errcheck // ToString called with known constants that cannot fail
	dependencyinjection.RegisterSingletonTenant[orm_base.DialectorGetter](
		register,
		mysqlKey,
		dialectors.NewMysqlDialectorGetter,
	)

	postgresKey, _ := orm_base.Postgres.ToString() //nolint:errcheck // ToString called with known constants that cannot fail
	dependencyinjection.RegisterSingletonTenant[orm_base.DialectorGetter](
		register,
		postgresKey,
		dialectors.NewPostgresDialectorGetter,
	)

	sqliteKey, _ := orm_base.Sqlite.ToString() //nolint:errcheck // ToString called with known constants that cannot fail
	dependencyinjection.RegisterSingletonTenant[orm_base.DialectorGetter](
		register,
		sqliteKey,
		dialectors.NewSqliteDialectorGetter,
	)

	sqlServerKey, _ := orm_base.SQLServer.ToString() //nolint:errcheck // ToString called with known constants that cannot fail
	dependencyinjection.RegisterSingletonTenant[orm_base.DialectorGetter](
		register,
		sqlServerKey,
		dialectors.NewSqlServerDialectorGetter,
	)

	return nil
}
