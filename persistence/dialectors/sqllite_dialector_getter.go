package dialectors

import (
	"os"
	"path/filepath"

	persistence "github.com/janmbaco/go-infrastructure/v2/persistence"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type sqliteDialectorGetter struct {
}

func NewSqliteDialectorGetter() persistence.DialectorGetter {
	return &sqliteDialectorGetter{}
}

func (dialectorGetter *sqliteDialectorGetter) Get(info *persistence.DatabaseInfo) (gorm.Dialector, error) {
	dbPath := info.Name
	if dbPath == "" {
		dbPath = info.Host
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, err
	}
	return sqlite.Open(dbPath), nil
}
