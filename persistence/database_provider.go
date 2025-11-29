package persistence //nolint:revive // established package name, changing would break API

import (
	"gorm.io/gorm"
)

func NewDB(dialectorResolver DialectorResolver, info *DatabaseInfo, config *gorm.Config, tables []interface{}) (*gorm.DB, error) {
	dialector, err := dialectorResolver.Resolve(info)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, err
	}
	for _, table := range tables {
		if err := db.AutoMigrate(table); err != nil {
			return nil, err
		}
	}
	return db, nil
}
