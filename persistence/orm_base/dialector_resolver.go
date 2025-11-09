package orm_base
import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	"gorm.io/gorm"
)

type DialectorResolver interface {
	Resolve(info *DatabaseInfo) (gorm.Dialector, error)
}

type dialectorResolver struct {
	resolver dependencyinjection.Resolver
}

func NewDialectorResolver(resolver dependencyinjection.Resolver) DialectorResolver {
	return &dialectorResolver{resolver: resolver}
}

func (dbResolver *dialectorResolver) Resolve(info *DatabaseInfo) (gorm.Dialector, error) {
	engineStr, err := info.Engine.ToString()
	if err != nil {
		return nil, err
	}
	
	getter := dependencyinjection.ResolveTenant[DialectorGetter](dbResolver.resolver, engineStr)
	return getter.Get(info)
}
