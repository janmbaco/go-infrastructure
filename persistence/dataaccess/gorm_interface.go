package dataaccess //nolint:revive // established package name, changing would break API

import (
	"gorm.io/gorm"
)

// GormDBInterface abstrae las operaciones de GORM que necesitamos para testing
type GormDBInterface interface {
	Model(value interface{}) GormDBInterface
	Create(value interface{}) GormDBInterface
	Where(query interface{}, args ...interface{}) GormDBInterface
	Preload(query string, args ...interface{}) GormDBInterface
	Find(dest interface{}, conds ...interface{}) GormDBInterface
	Updates(values interface{}) GormDBInterface
	Select(query interface{}, args ...interface{}) GormDBInterface
	Delete(value interface{}, conds ...interface{}) GormDBInterface
	GetError() error
	DB() interface{} // Expose underlying database for advanced operations
}

// gormDBWrapper envuelve *gorm.DB para implementar GormDBInterface
type gormDBWrapper struct {
	db *gorm.DB
}

func (w *gormDBWrapper) Model(value interface{}) GormDBInterface {
	return &gormDBWrapper{w.db.Model(value)}
}

func (w *gormDBWrapper) Create(value interface{}) GormDBInterface {
	return &gormDBWrapper{w.db.Create(value)}
}

func (w *gormDBWrapper) Where(query interface{}, args ...interface{}) GormDBInterface {
	return &gormDBWrapper{w.db.Where(query, args...)}
}

func (w *gormDBWrapper) Preload(query string, args ...interface{}) GormDBInterface {
	return &gormDBWrapper{w.db.Preload(query, args...)}
}

func (w *gormDBWrapper) Find(dest interface{}, conds ...interface{}) GormDBInterface {
	return &gormDBWrapper{w.db.Find(dest, conds...)}
}

func (w *gormDBWrapper) Updates(values interface{}) GormDBInterface {
	return &gormDBWrapper{w.db.Updates(values)}
}

func (w *gormDBWrapper) Select(query interface{}, args ...interface{}) GormDBInterface {
	return &gormDBWrapper{w.db.Select(query, args...)}
}

func (w *gormDBWrapper) Delete(value interface{}, conds ...interface{}) GormDBInterface {
	return &gormDBWrapper{w.db.Delete(value, conds...)}
}

func (w *gormDBWrapper) GetError() error {
	return w.db.Error
}

func (w *gormDBWrapper) DB() interface{} {
	return w.db
}
