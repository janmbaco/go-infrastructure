package dependencyinjection

import (
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"reflect"
)

// Register defines an object responsible to register the dependencies of a application
type Register interface {
	AsType(iface, provider interface{}, argNames map[uint]string)
	AsScope(iface, provider interface{}, argNames map[uint]string)
	AsSingleton(iface, provider interface{}, argNames map[uint]string)
	AsTenant(tenant string, iface, provider interface{}, argNames map[uint]string)
	AsSingletonTenant(tenant string, iface, provider interface{}, argNames map[uint]string)
	Bind(ifaceFrom, ifaceTo interface{})
}

type register struct {
	dependencies Dependencies
}

func newRegister(dependencies Dependencies) Register {
	return &register{dependencies: dependencies}
}

// AsType register that the dependecy goes to be provided by a provider and a args
func (r *register) AsType(iface, provider interface{}, argNames map[uint]string) {
	errorschecker.CheckNilParameter(map[string]interface{}{"iface": iface, "provider": provider})
	r.dependencies.Set(
		DependencyKey{Iface: reflect.Indirect(reflect.ValueOf(iface)).Type()},
		&dependencyObject{provider: provider, argNames: argNames},
	)
}

// AsSingleton register that the dependecy goes to be provided by a provider and a args like singleton
func (r *register) AsScope(iface, provider interface{}, argNames map[uint]string) {
	errorschecker.CheckNilParameter(map[string]interface{}{"iface": iface, "provider": provider})
	r.dependencies.Set(
		DependencyKey{Iface: reflect.Indirect(reflect.ValueOf(iface)).Type()},
		&dependencyObject{provider: provider, argNames: argNames, dependecyType: _ScopedType},
	)
}

// AsSingleton register that the dependecy goes to be provided by a provider and a args like singleton
func (r *register) AsSingleton(iface, provider interface{}, argNames map[uint]string) {
	errorschecker.CheckNilParameter(map[string]interface{}{"iface": iface, "provider": provider})
	r.dependencies.Set(
		DependencyKey{Iface: reflect.Indirect(reflect.ValueOf(iface)).Type()},
		&dependencyObject{provider: provider, argNames: argNames, dependecyType: _Singleton},
	)
}

// AsTenant register that the dependecy goes to be provided by a provider and a args with a tenant key
func (r *register) AsTenant(tenant string, iface, provider interface{}, argNames map[uint]string) {
	errorschecker.CheckNilParameter(map[string]interface{}{"iface": iface, "provider": provider})
	r.dependencies.Set(DependencyKey{
		Tenant: tenant,
		Iface:  reflect.Indirect(reflect.ValueOf(iface)).Type(),
	}, &dependencyObject{provider: provider, argNames: argNames})
}

// AsSingletonTenant register that the dependecy goes to be provided by a provider and a args with a tenant key as singleton
func (r *register) AsSingletonTenant(tenant string, iface, provider interface{}, argNames map[uint]string) {
	errorschecker.CheckNilParameter(map[string]interface{}{"iface": iface, "provider": provider})
	r.dependencies.Set(DependencyKey{
		Tenant: tenant,
		Iface:  reflect.Indirect(reflect.ValueOf(iface)).Type(),
	}, &dependencyObject{provider: provider, argNames: argNames, dependecyType: _Singleton})
}

// Bind registers a interface that is provided by a provider of another interface
func (r *register) Bind(ifaceFrom, ifaceTo interface{}) {
	errorschecker.CheckNilParameter(map[string]interface{}{"ifaceFrom": ifaceFrom, "ifaceTo": ifaceTo})
	r.dependencies.Bind(
		DependencyKey{Iface: reflect.Indirect(reflect.ValueOf(ifaceFrom)).Type()},
		DependencyKey{Iface: reflect.Indirect(reflect.ValueOf(ifaceTo)).Type()},
	)
}
