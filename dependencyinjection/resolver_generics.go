package dependencyinjection
import (
	"context"
	"reflect"
)

// Resolve resolves a dependency with type safety using generics
func Resolve[T any](resolver Resolver) T {
	var instance T
	result := resolver.Type(&instance, nil)
	return result.(T)
}

// ResolveWithParams resolves a dependency with parameters using generics
func ResolveWithParams[T any](resolver Resolver, params map[string]interface{}) T {
	var instance T
	result := resolver.Type(&instance, params)
	return result.(T)
}

// ResolveTenant resolves a tenant dependency with type safety using generics
func ResolveTenant[T any](resolver Resolver, tenant string) T {
	var instance T
	result := resolver.Tenant(tenant, &instance, nil)
	return result.(T)
}

// ResolveTenantWithParams resolves a tenant dependency with parameters using generics
func ResolveTenantWithParams[T any](resolver Resolver, tenant string, params map[string]interface{}) T {
	var instance T
	result := resolver.Tenant(tenant, &instance, params)
	return result.(T)
}

// MustResolve resolves a dependency and panics if it fails
func MustResolve[T any](resolver Resolver) T {
	result := Resolve[T](resolver)
	if reflect.ValueOf(result).IsNil() {
		var instance T
		typeName := reflect.TypeOf(&instance).Elem().String()
		panic("failed to resolve " + typeName)
	}
	return result
}

// TryResolve attempts to resolve a dependency and returns ok=false if it fails
func TryResolve[T any](resolver Resolver) (result T, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()
	result = Resolve[T](resolver)
	ok = true
	return
}

// ResolveOrDefault resolves a dependency or returns a default value
func ResolveOrDefault[T any](resolver Resolver, defaultValue T) T {
	result, ok := TryResolve[T](resolver)
	if !ok {
		return defaultValue
	}
	return result
}

// Context-aware generic resolution functions

// ResolveCtx resolves a dependency with context support
func ResolveCtx[T any](ctx context.Context, resolver Resolver) T {
	var instance T
	result := resolver.TypeCtx(ctx, &instance, nil)
	return result.(T)
}

// ResolveWithParamsCtx resolves with context and parameters
func ResolveWithParamsCtx[T any](ctx context.Context, resolver Resolver, params map[string]interface{}) T {
	var instance T
	result := resolver.TypeCtx(ctx, &instance, params)
	return result.(T)
}

// ResolveTenantCtx resolves a tenant dependency with context
func ResolveTenantCtx[T any](ctx context.Context, resolver Resolver, tenant string) T {
	var instance T
	result := resolver.TenantCtx(ctx, tenant, &instance, nil)
	return result.(T)
}

// ResolveTenantWithParamsCtx resolves tenant dependency with context and parameters
func ResolveTenantWithParamsCtx[T any](ctx context.Context, resolver Resolver, tenant string, params map[string]interface{}) T {
	var instance T
	result := resolver.TenantCtx(ctx, tenant, &instance, params)
	return result.(T)
}

// MustResolveCtx resolves with context and panics if it fails
func MustResolveCtx[T any](ctx context.Context, resolver Resolver) T {
	result := ResolveCtx[T](ctx, resolver)
	if reflect.ValueOf(result).IsNil() {
		var instance T
		typeName := reflect.TypeOf(&instance).Elem().String()
		panic("failed to resolve " + typeName)
	}
	return result
}

// TryResolveCtx attempts to resolve with context
func TryResolveCtx[T any](ctx context.Context, resolver Resolver) (result T, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()
	result = ResolveCtx[T](ctx, resolver)
	ok = true
	return
}
