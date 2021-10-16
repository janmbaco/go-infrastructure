package dialectors

import (
	"fmt"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"github.com/janmbaco/go-infrastructure/persistence/orm_base"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type sqlServerDialectorGetter struct {
}

func NewSqlServerDialectorGetter() orm_base.DialectorGetter {
	return &sqlServerDialectorGetter{}
}

func (dialectorGetter *sqlServerDialectorGetter) Get(info *orm_base.DatabaseInfo) gorm.Dialector {
	errorschecker.CheckNilParameter(map[string]interface{}{"info": info})
	return sqlserver.Open(fmt.Sprintf("sqlserver://%v:%v@%v:%v?orm_base=%v", info.UserName, info.UserPassword, info.Host, info.Port, info.Name))
}
