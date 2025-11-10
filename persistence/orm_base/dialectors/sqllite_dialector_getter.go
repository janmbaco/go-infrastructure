package dialectors

import (
	"os"
	"path/filepath"

	"github.com/janmbaco/go-infrastructure/v2/persistence/orm_base"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type sqliteDialectorGetter struct {
}

func NewSqliteDialectorGetter() orm_base.DialectorGetter {
	return &sqliteDialectorGetter{}
}

func (dialectorGetter *sqliteDialectorGetter) Get(info *orm_base.DatabaseInfo) (gorm.Dialector, error) {
	if err := os.MkdirAll(filepath.Dir(info.Host), 0o755); err != nil {
		return nil, err
	}
	return sqlite.Open(info.Host), nil
}
