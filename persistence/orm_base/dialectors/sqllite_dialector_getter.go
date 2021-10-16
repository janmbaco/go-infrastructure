package dialectors

import (
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"github.com/janmbaco/go-infrastructure/persistence/orm_base"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

type sqliteDialectorGetter struct {
}

func NewSqliteDialectorGetter() orm_base.DialectorGetter {
	return &sqliteDialectorGetter{}
}

func (dialectorGetter *sqliteDialectorGetter) Get(info *orm_base.DatabaseInfo) gorm.Dialector {
	errorschecker.CheckNilParameter(map[string]interface{}{"info": info})
	errorschecker.TryPanic(os.MkdirAll(filepath.Dir(info.Host), 0666))
	return sqlite.Open(info.Host)
}
