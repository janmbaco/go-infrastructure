package ioc

import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection/static"
	"github.com/janmbaco/go-infrastructure/persistence/orm_base"
	"github.com/janmbaco/go-infrastructure/persistence/orm_base/dialectors"
)

func init(){
	static.Container.Register().AsSingletonTenant(orm_base.MySql.ToString(), new(orm_base.DialectorGetter), dialectors.NewMysqlDialectorGetter, nil)
	static.Container.Register().AsSingletonTenant(orm_base.Postgres.ToString(), new(orm_base.DialectorGetter), dialectors.NewPostgresDialectorGetter, nil)
	static.Container.Register().AsSingletonTenant(orm_base.Sqlite.ToString(), new(orm_base.DialectorGetter), dialectors.NewSqliteDialectorGetter, nil)
	static.Container.Register().AsSingletonTenant(orm_base.SqlServer.ToString(), new(orm_base.DialectorGetter), dialectors.NewSqlServerDialectorGetter, nil)
}