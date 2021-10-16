package orm_base

import (
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"gorm.io/gorm"
)

func NewDB(dialectorResolver DialectorResolver, info *DatabaseInfo, config *gorm.Config, tables []interface{}) *gorm.DB {
	errorschecker.CheckNilParameter(map[string]interface{}{"dialectorResolver": dialectorResolver, "info": info, "config": config})

	db, err := gorm.Open(dialectorResolver.Resolve(info), config)
	errorschecker.TryPanic(err)
	for _, table := range tables {
		errorschecker.TryPanic(db.AutoMigrate(table))
	}
	return db
}
