package persistence //nolint:revive // established package name, changing would break API

import (
	"gorm.io/gorm"
)

type DialectorGetter interface {
	Get(info *DatabaseInfo) (gorm.Dialector, error)
}
