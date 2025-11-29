package resolver

import (
	"reflect"

	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	persistence "github.com/janmbaco/go-infrastructure/v2/persistence"
	"github.com/janmbaco/go-infrastructure/v2/persistence/dataaccess"
	"gorm.io/gorm"
)

func GetGormDB(resolver dependencyinjection.Resolver, info *persistence.DatabaseInfo, config *gorm.Config, tables []interface{}) *gorm.DB {
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

func GetDataAccess(resolver dependencyinjection.Resolver, modelType reflect.Type) dataaccess.DataAccess {
	result, ok := resolver.Type((*dataaccess.DataAccess)(nil), map[string]interface{}{"modelType": modelType}).(dataaccess.DataAccess)
	if !ok {
		panic("failed to resolve dataaccess.DataAccess")
	}
	return result
}
