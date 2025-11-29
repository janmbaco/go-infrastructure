package ioc

import (
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	persistence "github.com/janmbaco/go-infrastructure/v2/persistence"
	"github.com/janmbaco/go-infrastructure/v2/persistence/dataaccess"
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
	register.AsSingleton((*persistence.DialectorResolver)(nil), persistence.NewDialectorResolver, nil)
	register.AsSingleton(new(*gorm.DB), persistence.NewDB, map[int]string{1: "info", 2: "config", 3: "tables"})
	register.AsType((*dataaccess.DataAccess)(nil), dataaccess.NewDataAccess, map[int]string{2: "modelType"})

	return nil
}
