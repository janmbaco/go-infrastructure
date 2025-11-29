package dependencyinjection

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResolver_WhenCreated_ThenReturnsResolver(t *testing.T) {
	// Arrange
	deps := newDependencies()

	// Act
	res := newResolver(deps)

	// Assert
	assert.NotNil(t, res)
}

func TestResolver_Type_WhenRegistered_ThenResolves(t *testing.T) {
	// Arrange
	container := NewContainer()
	reg := container.Register()
	resolver := container.Resolver()
	reg.AsType(new(string), func() string { return "resolved" }, nil)

	// Act
	result := resolver.Type(new(string), nil)

	// Assert
	assert.Equal(t, "resolved", result)
}

func TestResolver_Type_WhenNotRegistered_ThenPanics(t *testing.T) {
	// Arrange
	container := NewContainer()
	resolver := container.Resolver()

	// Act & Assert
	assert.Panics(t, func() {
		resolver.Type(new(string), nil)
	})
}

func TestResolver_Tenant_WhenRegistered_ThenResolves(t *testing.T) {
	// Arrange
	container := NewContainer()
	reg := container.Register()
	resolver := container.Resolver()
	reg.AsTenant("tenant1", new(string), func() string { return "tenant-resolved" }, nil)

	// Act
	result := resolver.Tenant("tenant1", new(string), nil)

	// Assert
	assert.Equal(t, "tenant-resolved", result)
}

func TestResolver_Tenant_WhenNotRegistered_ThenPanics(t *testing.T) {
	// Arrange
	container := NewContainer()
	resolver := container.Resolver()

	// Act & Assert
	assert.Panics(t, func() {
		resolver.Tenant("tenant1", new(string), nil)
	})
}

func TestResolver_TypeCtx_WhenCalled_ThenDelegatesToType(t *testing.T) {
	// Arrange
	container := NewContainer()
	reg := container.Register()
	resolver := container.Resolver()
	ctx := context.Background()
	reg.AsType(new(string), func() string { return "ctx-resolved" }, nil)

	// Act
	result := resolver.TypeCtx(ctx, new(string), nil)

	// Assert
	assert.Equal(t, "ctx-resolved", result)
}

func TestResolver_TenantCtx_WhenCalled_ThenDelegatesToTenant(t *testing.T) {
	// Arrange
	container := NewContainer()
	reg := container.Register()
	resolver := container.Resolver()
	ctx := context.Background()
	reg.AsTenant("tenant1", new(string), func() string { return "tenant-ctx-resolved" }, nil)

	// Act
	result := resolver.TenantCtx(ctx, "tenant1", new(string), nil)

	// Assert
	assert.Equal(t, "tenant-ctx-resolved", result)
}
