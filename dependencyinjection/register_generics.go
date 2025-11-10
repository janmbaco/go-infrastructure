package dependencyinjection

import (
	"context"
)

// RegisterType registers a type with type safety using generics
func RegisterType[T any](r Register, factory func() T) {
	var instance T
	r.AsType(&instance, factory, nil)
}

// RegisterScoped registers a scoped dependency with type safety using generics
func RegisterScoped[T any](r Register, factory func() T) {
	var instance T
	r.AsScope(&instance, factory, nil)
}

// RegisterSingleton registers a singleton with type safety using generics
func RegisterSingleton[T any](r Register, factory func() T) {
	var instance T
	r.AsSingleton(&instance, factory, nil)
}

// RegisterTenant registers a tenant dependency with type safety using generics
func RegisterTenant[T any](r Register, tenant string, factory func() T) {
	var instance T
	r.AsTenant(tenant, &instance, factory, nil)
}

// RegisterSingletonTenant registers a singleton tenant with type safety using generics
func RegisterSingletonTenant[T any](r Register, tenant string, factory func() T) {
	var instance T
	r.AsSingletonTenant(tenant, &instance, factory, nil)
}

// RegisterTypeWithParams registers a type with parameters using generics
func RegisterTypeWithParams[T any](r Register, factory interface{}, argNames map[int]string) {
	var instance T
	r.AsType(&instance, factory, argNames)
}

// RegisterScopedWithParams registers a scoped dependency with parameters using generics
func RegisterScopedWithParams[T any](r Register, factory interface{}, argNames map[int]string) {
	var instance T
	r.AsScope(&instance, factory, argNames)
}

// RegisterSingletonWithParams registers a singleton with parameters using generics
func RegisterSingletonWithParams[T any](r Register, factory interface{}, argNames map[int]string) {
	var instance T
	r.AsSingleton(&instance, factory, argNames)
}

// RegisterTenantWithParams registers a tenant dependency with parameters using generics
func RegisterTenantWithParams[T any](r Register, tenant string, factory interface{}, argNames map[int]string) {
	var instance T
	r.AsTenant(tenant, &instance, factory, argNames)
}

// RegisterSingletonTenantWithParams registers a singleton tenant with parameters using generics
func RegisterSingletonTenantWithParams[T any](r Register, tenant string, factory interface{}, argNames map[int]string) {
	var instance T
	r.AsSingletonTenant(tenant, &instance, factory, argNames)
}

// Context-aware generic registration functions

// RegisterTypeCtx registers a type with context support
func RegisterTypeCtx[T any](ctx context.Context, r Register, factory func(context.Context) (T, error)) {
	var instance T
	r.AsTypeCtx(ctx, &instance, factory, nil)
}

// RegisterScopedCtx registers a scoped dependency with context support
func RegisterScopedCtx[T any](ctx context.Context, r Register, factory func(context.Context) (T, error)) {
	var instance T
	r.AsScopeCtx(ctx, &instance, factory, nil)
}

// RegisterSingletonCtx registers a singleton with context support
func RegisterSingletonCtx[T any](ctx context.Context, r Register, factory func(context.Context) (T, error)) {
	var instance T
	r.AsSingletonCtx(ctx, &instance, factory, nil)
}

// RegisterTenantCtx registers a tenant dependency with context support
func RegisterTenantCtx[T any](ctx context.Context, r Register, tenant string, factory func(context.Context) (T, error)) {
	var instance T
	r.AsTenantCtx(ctx, tenant, &instance, factory, nil)
}

// RegisterSingletonTenantCtx registers a singleton tenant with context support
func RegisterSingletonTenantCtx[T any](ctx context.Context, r Register, tenant string, factory func(context.Context) (T, error)) {
	var instance T
	r.AsSingletonTenantCtx(ctx, tenant, &instance, factory, nil)
}
