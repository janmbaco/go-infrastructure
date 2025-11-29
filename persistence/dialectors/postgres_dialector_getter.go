package dialectors

import (
	"fmt"

	persistence "github.com/janmbaco/go-infrastructure/v2/persistence"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDialectorGetter struct {
}

func NewPostgresDialectorGetter() persistence.DialectorGetter {
	return &postgresDialectorGetter{}
}

func (dialectorGetter *postgresDialectorGetter) Get(info *persistence.DatabaseInfo) (gorm.Dialector, error) {
	return postgres.Open(fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v", info.Host, info.Port, info.UserName, info.Name, info.UserPassword)), nil
}
