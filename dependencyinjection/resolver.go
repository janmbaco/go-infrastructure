package dependencyinjection

import (
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"reflect"
)

// Resolver defines an object responsible to resolver the dependencies of a application
type Resolver interface {
	Type(iface interface{}, params map[string]interface{}) interface{}
	Tenant(tenantName string, iface interface{}, params map[string]interface{}) interface{}
}

type resolver struct {
	dependencies Dependencies
}

func newResolver(dependencies Dependencies) Resolver {
	return &resolver{dependencies: dependencies}
}

// Type gets a dependency by the interface and params
func (r *resolver) Type(iface interface{}, params map[string]interface{}) interface{} {
	errorschecker.CheckNilParameter(map[string]interface{}{"iface": iface})
	return r.dependencies.Get(DependencyKey{Iface:reflect.Indirect(reflect.ValueOf(iface)).Type()}).Create(params, r.dependencies, make(map[DependencyObject]interface{}))
}

// Tenant gets a dependency by the interface, the tenant key and paramas
func (r *resolver) Tenant(tenant string, iface interface{}, params map[string]interface{}) interface{} {
	errorschecker.CheckNilParameter(map[string]interface{}{"iface": iface})

	return r.dependencies.Get(DependencyKey{
		Tenant: tenant,
		Iface: reflect.Indirect(reflect.ValueOf(iface)).Type(),
	}).Create(params, r.dependencies, make(map[DependencyObject]interface{}))
}
