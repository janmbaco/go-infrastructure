package ioc

import (
	"gorm.io/gorm"

	"github.com/janmbaco/go-infrastructure/dependencyinjection/static"
	"github.com/janmbaco/go-infrastructure/persistence/orm_base"
)

func init(){
	static.Container.Register().AsSingleton(new(orm_base.DialectorResolver), orm_base.NewDialectorResolver, nil)
	static.Container.Register().AsSingleton(new(*gorm.DB), orm_base.NewDB, map[uint]string{1: "info", 2: "config", 3: "tables"})
	static.Container.Register().AsType(new(orm_base.DataAccess), orm_base.NewDataAccess, map[uint]string{2: "modelType"})
}