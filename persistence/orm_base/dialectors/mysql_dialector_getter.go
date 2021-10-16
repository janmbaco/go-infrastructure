package dialectors

import (
	"fmt"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"github.com/janmbaco/go-infrastructure/persistence/orm_base"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mysqlDialectorGetter struct {
}

func NewMysqlDialectorGetter() orm_base.DialectorGetter {
	return &mysqlDialectorGetter{}
}

func (dialectorGetter *mysqlDialectorGetter) Get(info *orm_base.DatabaseInfo) gorm.Dialector {
	errorschecker.CheckNilParameter(map[string]interface{}{"info": info})
	return mysql.Open(fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local", info.UserName, info.UserPassword, info.Host, info.Port, info.Name))
}
