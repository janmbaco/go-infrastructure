package dependencyinjection

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRegister_WhenCreated_ThenReturnsRegister(t *testing.T) {
	// Arrange
	deps := newDependencies()

	// Act
	reg := newRegister(deps)

	// Assert
	assert.NotNil(t, reg)
}

func TestRegister_AsType_WhenRegistered_ThenCanResolve(t *testing.T) {
	// Arrange
	container := NewContainer()
	reg := container.Register()
	resolver := container.Resolver()

	// Act
	reg.AsType(new(string), func() string { return "test" }, nil)

	// Assert
	result := resolver.Type(new(string), nil)
	assert.Equal(t, "test", result)
}

func TestRegister_AsScope_WhenRegistered_ThenCanResolve(t *testing.T) {
	// Arrange
	container := NewContainer()
	reg := container.Register()
	resolver := container.Resolver()

	// Act
	reg.AsScope(new(string), func() string { return "scoped" }, nil)

	// Assert
	result := resolver.Type(new(string), nil)
	assert.Equal(t, "scoped", result)
}

func TestRegister_AsSingleton_WhenRegistered_ThenReturnsSameInstance(t *testing.T) {
	// Arrange
	container := NewContainer()
	reg := container.Register()
	resolver := container.Resolver()

	// Act
	reg.AsSingleton(new(string), func() string { return "singleton" }, nil)

	// Assert
	first := resolver.Type(new(string), nil)
	second := resolver.Type(new(string), nil)
	assert.Equal(t, "singleton", first)
	assert.Equal(t, first, second)
}

func TestRegister_AsTenant_WhenRegistered_ThenCanResolve(t *testing.T) {
	// Arrange
	container := NewContainer()
	reg := container.Register()
	resolver := container.Resolver()

	// Act
	reg.AsTenant("tenant1", new(string), func() string { return "tenant" }, nil)

	// Assert
	result := resolver.Tenant("tenant1", new(string), nil)
	assert.Equal(t, "tenant", result)
}

func TestRegister_AsSingletonTenant_WhenRegistered_ThenReturnsSameInstance(t *testing.T) {
	// Arrange
	container := NewContainer()
	reg := container.Register()
	resolver := container.Resolver()

	// Act
	reg.AsSingletonTenant("tenant1", new(string), func() string { return "singleton-tenant" }, nil)

	// Assert
	first := resolver.Tenant("tenant1", new(string), nil)
	second := resolver.Tenant("tenant1", new(string), nil)
	assert.Equal(t, "singleton-tenant", first)
	assert.Equal(t, first, second)
}

func TestRegister_Bind_WhenBound_ThenResolvesToBoundType(t *testing.T) {
	// Arrange
	container := NewContainer()
	reg := container.Register()
	resolver := container.Resolver()

	// Act
	reg.AsType(new(int), func() int { return 42 }, nil)
	reg.Bind(new(string), new(int))

	// Assert
	result := resolver.Type(new(string), nil)
	assert.Equal(t, 42, result)
}

func TestRegister_AsTypeCtx_WhenCalled_ThenDelegatesToAsType(t *testing.T) {
	// Arrange
	container := NewContainer()
	reg := container.Register()
	resolver := container.Resolver()
	ctx := context.Background()

	// Act
	reg.AsTypeCtx(ctx, new(string), func() string { return "ctx" }, nil)

	// Assert
	result := resolver.Type(new(string), nil)
	assert.Equal(t, "ctx", result)
}

func TestRegister_AsScopeCtx_WhenCalled_ThenDelegatesToAsScope(t *testing.T) {
	// Arrange
	container := NewContainer()
	reg := container.Register()
	resolver := container.Resolver()
	ctx := context.Background()

	// Act
	reg.AsScopeCtx(ctx, new(string), func() string { return "scoped-ctx" }, nil)

	// Assert
	result := resolver.Type(new(string), nil)
	assert.Equal(t, "scoped-ctx", result)
}

func TestRegister_AsSingletonCtx_WhenCalled_ThenDelegatesToAsSingleton(t *testing.T) {
	// Arrange
	container := NewContainer()
	reg := container.Register()
	resolver := container.Resolver()
	ctx := context.Background()

	// Act
	reg.AsSingletonCtx(ctx, new(string), func() string { return "singleton-ctx" }, nil)

	// Assert
	first := resolver.Type(new(string), nil)
	second := resolver.Type(new(string), nil)
	assert.Equal(t, "singleton-ctx", first)
	assert.Equal(t, first, second)
}

func TestRegister_AsTenantCtx_WhenCalled_ThenDelegatesToAsTenant(t *testing.T) {
	// Arrange
	container := NewContainer()
	reg := container.Register()
	resolver := container.Resolver()
	ctx := context.Background()

	// Act
	reg.AsTenantCtx(ctx, "tenant1", new(string), func() string { return "tenant-ctx" }, nil)

	// Assert
	result := resolver.Tenant("tenant1", new(string), nil)
	assert.Equal(t, "tenant-ctx", result)
}

func TestRegister_AsSingletonTenantCtx_WhenCalled_ThenDelegatesToAsSingletonTenant(t *testing.T) {
	// Arrange
	container := NewContainer()
	reg := container.Register()
	resolver := container.Resolver()
	ctx := context.Background()

	// Act
	reg.AsSingletonTenantCtx(ctx, "tenant1", new(string), func() string { return "singleton-tenant-ctx" }, nil)

	// Assert
	first := resolver.Tenant("tenant1", new(string), nil)
	second := resolver.Tenant("tenant1", new(string), nil)
	assert.Equal(t, "singleton-tenant-ctx", first)
	assert.Equal(t, first, second)
}
