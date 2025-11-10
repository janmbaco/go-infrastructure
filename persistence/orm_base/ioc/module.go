package ioc

import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/persistence/orm_base"
	"gorm.io/gorm"
)

// PersistenceModule implements Module for persistence services
type PersistenceModule struct{}

// NewPersistenceModule creates a new persistence module
func NewPersistenceModule() *PersistenceModule {
	return &PersistenceModule{}
}

// RegisterServices registers all persistence services
func (m *PersistenceModule) RegisterServices(register dependencyinjection.Register) error {
	register.AsSingleton(new(orm_base.DialectorResolver), orm_base.NewDialectorResolver, nil)
	register.AsSingleton(new(*gorm.DB), orm_base.NewDB, map[int]string{1: "info", 2: "config", 3: "tables"})
	register.AsType(new(orm_base.DataAccess), orm_base.NewDataAccess, map[int]string{2: "modelType"})

	return nil
}
