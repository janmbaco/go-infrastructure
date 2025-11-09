package dependencyinjection
import (
	"context"
	"reflect"
)

// Resolver defines an object responsible to resolver the dependencies of a application
type Resolver interface {
	Type(iface interface{}, params map[string]interface{}) interface{}
	Tenant(tenantName string, iface interface{}, params map[string]interface{}) interface{}
	
	// Context-aware methods
	TypeCtx(ctx context.Context, iface interface{}, params map[string]interface{}) interface{}
	TenantCtx(ctx context.Context, tenantName string, iface interface{}, params map[string]interface{}) interface{}
}

type resolver struct {
	dependencies Dependencies
}

func newResolver(dependencies Dependencies) Resolver {
	return &resolver{dependencies: dependencies}
}

// Type gets a dependency by the interface and params
func (r *resolver) Type(iface interface{}, params map[string]interface{}) interface{} {
	return r.dependencies.Get(DependencyKey{Iface:reflect.Indirect(reflect.ValueOf(iface)).Type()}).Create(params, r.dependencies, make(map[DependencyObject]interface{}))
}

// Tenant gets a dependency by the interface, the tenant key and paramas
func (r *resolver) Tenant(tenant string, iface interface{}, params map[string]interface{}) interface{} {
	return r.dependencies.Get(DependencyKey{
		Tenant: tenant,
		Iface: reflect.Indirect(reflect.ValueOf(iface)).Type(),
	}).Create(params, r.dependencies, make(map[DependencyObject]interface{}))
}

// TypeCtx resolves with context (delegates to Type for now)
func (r *resolver) TypeCtx(ctx context.Context, iface interface{}, params map[string]interface{}) interface{} {
	return r.Type(iface, params)
}

// TenantCtx resolves with context (delegates to Tenant for now)
func (r *resolver) TenantCtx(ctx context.Context, tenant string, iface interface{}, params map[string]interface{}) interface{} {
	return r.Tenant(tenant, iface, params)
}
