package dialectors

import (
	"fmt"

	persistence "github.com/janmbaco/go-infrastructure/v2/persistence"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mysqlDialectorGetter struct {
}

func NewMysqlDialectorGetter() persistence.DialectorGetter {
	return &mysqlDialectorGetter{}
}

func (dialectorGetter *mysqlDialectorGetter) Get(info *persistence.DatabaseInfo) (gorm.Dialector, error) {
	return mysql.Open(fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local", info.UserName, info.UserPassword, info.Host, info.Port, info.Name)), nil
}
