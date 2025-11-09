package orm_base
import (
	"gorm.io/gorm"
)

type DialectorGetter interface {
	Get(info *DatabaseInfo) (gorm.Dialector, error)
}
