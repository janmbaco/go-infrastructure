package dependencyinjection

import (
	"context"
	"reflect"
)

// Register defines an object responsible to register the dependencies of a application
type Register interface {
	AsType(iface, provider interface{}, argNames map[int]string)
	AsScope(iface, provider interface{}, argNames map[int]string)
	AsSingleton(iface, provider interface{}, argNames map[int]string)
	AsTenant(tenant string, iface, provider interface{}, argNames map[int]string)
	AsSingletonTenant(tenant string, iface, provider interface{}, argNames map[int]string)
	Bind(ifaceFrom, ifaceTo interface{})

	// Context-aware methods
	AsTypeCtx(ctx context.Context, iface, provider interface{}, argNames map[int]string)
	AsScopeCtx(ctx context.Context, iface, provider interface{}, argNames map[int]string)
	AsSingletonCtx(ctx context.Context, iface, provider interface{}, argNames map[int]string)
	AsTenantCtx(ctx context.Context, tenant string, iface, provider interface{}, argNames map[int]string)
	AsSingletonTenantCtx(ctx context.Context, tenant string, iface, provider interface{}, argNames map[int]string)
}

type register struct {
	dependencies Dependencies
}

func newRegister(dependencies Dependencies) Register {
	return &register{dependencies: dependencies}
}

// AsType register that the dependecy goes to be provided by a provider and a args
func (r *register) AsType(iface, provider interface{}, argNames map[int]string) {
	r.dependencies.Set(
		DependencyKey{Iface: reflect.Indirect(reflect.ValueOf(iface)).Type()},
		&dependencyObject{provider: provider, argNames: argNames},
	)
}

// AsSingleton register that the dependecy goes to be provided by a provider and a args like singleton
func (r *register) AsScope(iface, provider interface{}, argNames map[int]string) {
	r.dependencies.Set(
		DependencyKey{Iface: reflect.Indirect(reflect.ValueOf(iface)).Type()},
		&dependencyObject{provider: provider, argNames: argNames, dependecyType: _ScopedType},
	)
}

// AsSingleton register that the dependecy goes to be provided by a provider and a args like singleton
func (r *register) AsSingleton(iface, provider interface{}, argNames map[int]string) {
	r.dependencies.Set(
		DependencyKey{Iface: reflect.Indirect(reflect.ValueOf(iface)).Type()},
		&dependencyObject{provider: provider, argNames: argNames, dependecyType: _Singleton},
	)
}

// AsTenant register that the dependecy goes to be provided by a provider and a args with a tenant key
func (r *register) AsTenant(tenant string, iface, provider interface{}, argNames map[int]string) {
	r.dependencies.Set(DependencyKey{
		Tenant: tenant,
		Iface:  reflect.Indirect(reflect.ValueOf(iface)).Type(),
	}, &dependencyObject{provider: provider, argNames: argNames})
}

// AsSingletonTenant register that the dependecy goes to be provided by a provider and a args with a tenant key as singleton
func (r *register) AsSingletonTenant(tenant string, iface, provider interface{}, argNames map[int]string) {
	r.dependencies.Set(DependencyKey{
		Tenant: tenant,
		Iface:  reflect.Indirect(reflect.ValueOf(iface)).Type(),
	}, &dependencyObject{provider: provider, argNames: argNames, dependecyType: _Singleton})
}

// Bind registers a interface that is provided by a provider of another interface
func (r *register) Bind(ifaceFrom, ifaceTo interface{}) {
	r.dependencies.Bind(
		DependencyKey{Iface: reflect.Indirect(reflect.ValueOf(ifaceFrom)).Type()},
		DependencyKey{Iface: reflect.Indirect(reflect.ValueOf(ifaceTo)).Type()},
	)
}

// AsTypeCtx register with context (delegates to AsType for now)
func (r *register) AsTypeCtx(ctx context.Context, iface, provider interface{}, argNames map[int]string) {
	r.AsType(iface, provider, argNames)
}

// AsScopeCtx register with context (delegates to AsScope for now)
func (r *register) AsScopeCtx(ctx context.Context, iface, provider interface{}, argNames map[int]string) {
	r.AsScope(iface, provider, argNames)
}

// AsSingletonCtx register with context (delegates to AsSingleton for now)
func (r *register) AsSingletonCtx(ctx context.Context, iface, provider interface{}, argNames map[int]string) {
	r.AsSingleton(iface, provider, argNames)
}

// AsTenantCtx register with context (delegates to AsTenant for now)
func (r *register) AsTenantCtx(ctx context.Context, tenant string, iface, provider interface{}, argNames map[int]string) {
	r.AsTenant(tenant, iface, provider, argNames)
}

// AsSingletonTenantCtx register with context (delegates to AsSingletonTenant for now)
func (r *register) AsSingletonTenantCtx(ctx context.Context, tenant string, iface, provider interface{}, argNames map[int]string) {
	r.AsSingletonTenant(tenant, iface, provider, argNames)
}
