package resolver

import (
	"reflect"

	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/persistence/orm_base"
	"gorm.io/gorm"
)

func GetgormDB(resolver dependencyinjection.Resolver, info *orm_base.DatabaseInfo, config *gorm.Config, tables []interface{}) *gorm.DB {
	result, ok := resolver.Type(new(*gorm.DB), map[string]interface{}{
		"info":   info,
		"config": config,
		"tables": tables,
	}).(*gorm.DB)
	if !ok {
		panic("failed to resolve *gorm.DB")
	}
	return result
}

func GetDataAccess(resolver dependencyinjection.Resolver, modelType reflect.Type) orm_base.DataAccess {
	result, ok := resolver.Type(new(orm_base.DataAccess), map[string]interface{}{"modelType": modelType}).(orm_base.DataAccess)
	if !ok {
		panic("failed to resolve orm_base.DataAccess")
	}
	return result
}
