package orm_base

import (
	"github.com/janmbaco/go-infrastructure/dependencyinjection"
	"gorm.io/gorm"
)

type DialectorResolver interface {
	Resolve(info *DatabaseInfo) gorm.Dialector
}

type dialectorResolver struct {
	resolver dependencyinjection.Resolver
}

func NewDialectorResolver(resolver dependencyinjection.Resolver) DialectorResolver {
	return &dialectorResolver{resolver: resolver}
}

func (dbResolver *dialectorResolver) Resolve(info *DatabaseInfo) gorm.Dialector {
	return dbResolver.resolver.Tenant(info.Engine.ToString(), new(DialectorGetter), nil).(DialectorGetter).Get(info)
}
