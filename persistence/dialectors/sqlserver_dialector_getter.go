package dialectors

import (
	"fmt"

	persistence "github.com/janmbaco/go-infrastructure/v2/persistence"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type sqlServerDialectorGetter struct {
}

func NewSqlServerDialectorGetter() persistence.DialectorGetter { //nolint:revive // established API name, changing would break API
	return &sqlServerDialectorGetter{}
}

func (dialectorGetter *sqlServerDialectorGetter) Get(info *persistence.DatabaseInfo) (gorm.Dialector, error) {
	return sqlserver.Open(fmt.Sprintf("sqlserver://%v:%v@%v:%v?database=%v", info.UserName, info.UserPassword, info.Host, info.Port, info.Name)), nil
}
