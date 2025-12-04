package dependencyinjection

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestResolver_TypeCtx_WhenProviderAcceptsContext_ThenPassesContext tests that
// providers can receive context.Context as first parameter
func TestResolver_TypeCtx_WhenProviderAcceptsContext_ThenPassesContext(t *testing.T) {
	// Arrange
	type ctxKey string
	const testKey ctxKey = "testKey"
	const testValue = "testValue"

	ctx := context.WithValue(context.Background(), testKey, testValue)

	container := NewContainer()
	var receivedCtx context.Context

	container.Register().AsType(new(string), func(c context.Context) string {
		receivedCtx = c
		return c.Value(testKey).(string)
	}, nil)

	// Act
	result := container.Resolver().TypeCtx(ctx, new(string), nil)

	// Assert
	assert.NotNil(t, receivedCtx)
	assert.Equal(t, testValue, result)
	assert.Equal(t, testValue, receivedCtx.Value(testKey))
}

// TestResolver_TypeCtx_WhenContextCanceled_ThenPanics tests that
// resolution respects context cancellation
func TestResolver_TypeCtx_WhenContextCanceled_ThenPanics(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	container := NewContainer()
	container.Register().AsType(new(string), func() string {
		return "test"
	}, nil)

	// Act & Assert
	assert.Panics(t, func() {
		container.Resolver().TypeCtx(ctx, new(string), nil)
	})
}

// TestResolver_TypeCtx_WhenContextTimedOut_ThenPanics tests that
// resolution respects context timeout
func TestResolver_TypeCtx_WhenContextTimedOut_ThenPanics(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(10 * time.Millisecond) // Ensure timeout has passed

	container := NewContainer()
	container.Register().AsType(new(string), func() string {
		return "test"
	}, nil)

	// Act & Assert
	assert.Panics(t, func() {
		container.Resolver().TypeCtx(ctx, new(string), nil)
	})
}

// TestResolver_TypeCtx_WhenProviderWithContextReturnsError_ThenPanics tests
// error handling with context-aware providers
func TestResolver_TypeCtx_WhenProviderWithContextReturnsError_ThenPanics(t *testing.T) {
	// Arrange
	ctx := context.Background()
	expectedError := errors.New("provider error")

	container := NewContainer()
	container.Register().AsType(new(string), func(c context.Context) (string, error) {
		return "", expectedError
	}, nil)

	// Act & Assert
	assert.Panics(t, func() {
		container.Resolver().TypeCtx(ctx, new(string), nil)
	})
}

// TestResolver_TypeCtx_WhenNestedDependenciesWithContext_ThenPropagatesContext tests
// that context is propagated through nested dependency resolution
func TestResolver_TypeCtx_WhenNestedDependenciesWithContext_ThenPropagatesContext(t *testing.T) {
	// Arrange
	type ctxKey string
	const testKey ctxKey = "requestID"
	const requestID = "req-12345"

	ctx := context.WithValue(context.Background(), testKey, requestID)

	type Inner struct {
		Value string
	}

	type Outer struct {
		Inner *Inner
	}

	container := NewContainer()

	// Register inner service that uses context
	container.Register().AsType(new(*Inner), func(c context.Context) *Inner {
		return &Inner{Value: c.Value(testKey).(string)}
	}, nil)

	// Register outer service that depends on inner and also uses context
	container.Register().AsType(new(*Outer), func(c context.Context, inner *Inner) *Outer {
		// Verify context is passed correctly
		assert.Equal(t, requestID, c.Value(testKey))
		return &Outer{Inner: inner}
	}, nil)

	// Act
	result := container.Resolver().TypeCtx(ctx, new(*Outer), nil).(*Outer)

	// Assert
	assert.NotNil(t, result)
	assert.NotNil(t, result.Inner)
	assert.Equal(t, requestID, result.Inner.Value)
}

// TestResolver_TenantCtx_WhenProviderAcceptsContext_ThenPassesContext tests
// that tenant resolution also supports context
func TestResolver_TenantCtx_WhenProviderAcceptsContext_ThenPassesContext(t *testing.T) {
	// Arrange
	type ctxKey string
	const testKey ctxKey = "tenant"
	const tenantName = "tenant1"
	const tenantValue = "Tenant One"

	ctx := context.WithValue(context.Background(), testKey, tenantValue)

	container := NewContainer()
	container.Register().AsTenant(tenantName, new(string), func(c context.Context) string {
		return c.Value(testKey).(string)
	}, nil)

	// Act
	result := container.Resolver().TenantCtx(ctx, tenantName, new(string), nil)

	// Assert
	assert.Equal(t, tenantValue, result)
}

// TestResolver_Type_WhenNoContext_ThenUsesBackgroundContext tests
// that methods without Ctx suffix use background context internally
func TestResolver_Type_WhenNoContext_ThenUsesBackgroundContext(t *testing.T) {
	// Arrange
	container := NewContainer()
	var receivedCtx context.Context

	container.Register().AsType(new(string), func(c context.Context) string {
		receivedCtx = c
		return "test"
	}, nil)

	// Act
	result := container.Resolver().Type(new(string), nil)

	// Assert
	assert.Equal(t, "test", result)
	assert.NotNil(t, receivedCtx)
	// Background context should not be nil
	assert.Equal(t, context.Background(), receivedCtx)
}

// TestResolver_TypeCtx_WhenSingletonWithContext_ThenContextOnlyUsedFirstTime tests
// that singleton instances cache the result and context is only used on first creation
func TestResolver_TypeCtx_WhenSingletonWithContext_ThenContextOnlyUsedFirstTime(t *testing.T) {
	// Arrange
	type ctxKey string
	const testKey ctxKey = "counter"

	ctx1 := context.WithValue(context.Background(), testKey, 1)
	ctx2 := context.WithValue(context.Background(), testKey, 2)

	container := NewContainer()
	callCount := 0

	container.Register().AsSingleton(new(*int), func(c context.Context) *int {
		callCount++
		value := c.Value(testKey).(int)
		return &value
	}, nil)

	// Act
	result1 := container.Resolver().TypeCtx(ctx1, new(*int), nil).(*int)
	result2 := container.Resolver().TypeCtx(ctx2, new(*int), nil).(*int)

	// Assert
	assert.Equal(t, 1, callCount, "Provider should only be called once for singleton")
	assert.Equal(t, 1, *result1, "First resolution should use ctx1")
	assert.Equal(t, 1, *result2, "Second resolution should return cached singleton")
	assert.Equal(t, result1, result2, "Should return same instance")
}

// TestResolver_TypeCtx_WhenMixedProvidersWithAndWithoutContext_ThenBothWork tests
// that providers with and without context can coexist
func TestResolver_TypeCtx_WhenMixedProvidersWithAndWithoutContext_ThenBothWork(t *testing.T) {
	// Arrange
	type ctxKey string
	const testKey ctxKey = "value"

	ctx := context.WithValue(context.Background(), testKey, "from-context")

	type WithCtx struct {
		Value string
	}

	type WithoutCtx struct {
		Value string
	}

	container := NewContainer()

	container.Register().AsType(new(*WithCtx), func(c context.Context) *WithCtx {
		return &WithCtx{Value: c.Value(testKey).(string)}
	}, nil)

	container.Register().AsType(new(*WithoutCtx), func() *WithoutCtx {
		return &WithoutCtx{Value: "hardcoded"}
	}, nil)

	// Act
	resultWithCtx := container.Resolver().TypeCtx(ctx, new(*WithCtx), nil).(*WithCtx)
	resultWithoutCtx := container.Resolver().TypeCtx(ctx, new(*WithoutCtx), nil).(*WithoutCtx)

	// Assert
	assert.Equal(t, "from-context", resultWithCtx.Value)
	assert.Equal(t, "hardcoded", resultWithoutCtx.Value)
}
