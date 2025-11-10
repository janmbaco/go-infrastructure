package dialectors

import (
	"fmt"

	"github.com/janmbaco/go-infrastructure/persistence/orm_base"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type sqlServerDialectorGetter struct {
}

func NewSqlServerDialectorGetter() orm_base.DialectorGetter { //nolint:revive // established API name, changing would break API
	return &sqlServerDialectorGetter{}
}

func (dialectorGetter *sqlServerDialectorGetter) Get(info *orm_base.DatabaseInfo) (gorm.Dialector, error) {
	return sqlserver.Open(fmt.Sprintf("sqlserver://%v:%v@%v:%v?orm_base=%v", info.UserName, info.UserPassword, info.Host, info.Port, info.Name)), nil
}
