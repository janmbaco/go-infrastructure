package resolver

import (
	"reflect"

	"github.com/janmbaco/go-infrastructure/dependencyinjection/static"
	"github.com/janmbaco/go-infrastructure/persistence/orm_base"
	"gorm.io/gorm"

	_ "github.com/janmbaco/go-infrastructure/logs/ioc"
	_ "github.com/janmbaco/go-infrastructure/errors/ioc"
	_ "github.com/janmbaco/go-infrastructure/persistence/orm_base/dialectors/ioc"
	_ "github.com/janmbaco/go-infrastructure/persistence/orm_base/ioc"
)

func GetgormDB(info *orm_base.DatabaseInfo, config *gorm.Config, tables []interface{}) *gorm.DB {
 	return static.Container.Resolver().Type(new(*gorm.DB), map[string]interface{}{
			"info": info,
			"config": config,
			"tables": tables,
		}).(*gorm.DB)
}

func GetDataAccess(modelType reflect.Type) orm_base.DataAccess{
	return static.Container.Resolver().Type(new(orm_base.DataAccess), map[string]interface{}{"modelType": modelType}).(orm_base.DataAccess)
}