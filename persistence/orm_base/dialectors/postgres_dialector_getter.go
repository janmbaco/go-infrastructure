package dialectors

import (
	"fmt"

	"github.com/janmbaco/go-infrastructure/persistence/orm_base"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDialectorGetter struct {
}

func NewPostgresDialectorGetter() orm_base.DialectorGetter {
	return &postgresDialectorGetter{}
}

func (dialectorGetter *postgresDialectorGetter) Get(info *orm_base.DatabaseInfo) (gorm.Dialector, error) {
	return postgres.Open(fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v", info.Host, info.Port, info.UserName, info.Name, info.UserPassword)), nil
}
