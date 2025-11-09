package ioc
import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/persistence/orm_base"
	"github.com/janmbaco/go-infrastructure/persistence/orm_base/dialectors"
)

// DialectorsModule implements Module for database dialectors
type DialectorsModule struct{}

// NewDialectorsModule creates a new dialectors module
func NewDialectorsModule() *DialectorsModule {
	return &DialectorsModule{}
}

// RegisterServices registers all database dialectors
func (m *DialectorsModule) RegisterServices(register dependencyinjection.Register) error {
	mysqlKey, _ := orm_base.MySql.ToString()
	dependencyinjection.RegisterSingletonTenant[orm_base.DialectorGetter](
		register,
		mysqlKey,
		func() orm_base.DialectorGetter {
			return dialectors.NewMysqlDialectorGetter()
		},
	)

	postgresKey, _ := orm_base.Postgres.ToString()
	dependencyinjection.RegisterSingletonTenant[orm_base.DialectorGetter](
		register,
		postgresKey,
		func() orm_base.DialectorGetter {
			return dialectors.NewPostgresDialectorGetter()
		},
	)

	sqliteKey, _ := orm_base.Sqlite.ToString()
	dependencyinjection.RegisterSingletonTenant[orm_base.DialectorGetter](
		register,
		sqliteKey,
		func() orm_base.DialectorGetter {
			return dialectors.NewSqliteDialectorGetter()
		},
	)

	sqlServerKey, _ := orm_base.SqlServer.ToString()
	dependencyinjection.RegisterSingletonTenant[orm_base.DialectorGetter](
		register,
		sqlServerKey,
		func() orm_base.DialectorGetter {
			return dialectors.NewSqlServerDialectorGetter()
		},
	)

	return nil
}
