package dependencyinjection

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBuilder_WhenCreated_ThenReturnsBuilder(t *testing.T) {
	// Arrange & Act
	builder := NewBuilder()

	// Assert
	assert.NotNil(t, builder)
	assert.NotNil(t, builder.container)
	assert.Empty(t, builder.modules)
	assert.Empty(t, builder.modulesCtx)
}

func TestBuilder_AddModule_WhenAdded_ThenModuleIsInList(t *testing.T) {
	// Arrange
	builder := NewBuilder()
	module := &mockModule{}

	// Act
	result := builder.AddModule(module)

	// Assert
	assert.Equal(t, builder, result)
	assert.Len(t, builder.modules, 1)
	assert.Contains(t, builder.modules, module)
}

func TestBuilder_AddModules_WhenAdded_ThenModulesAreInList(t *testing.T) {
	// Arrange
	builder := NewBuilder()
	module1 := &mockModule{}
	module2 := &mockModule{}

	// Act
	result := builder.AddModules(module1, module2)

	// Assert
	assert.Equal(t, builder, result)
	assert.Len(t, builder.modules, 2)
	assert.Contains(t, builder.modules, module1)
	assert.Contains(t, builder.modules, module2)
}

func TestBuilder_AddModuleWithContext_WhenAdded_ThenModuleIsInList(t *testing.T) {
	// Arrange
	builder := NewBuilder()
	module := &mockModuleWithContext{}

	// Act
	result := builder.AddModuleWithContext(module)

	// Assert
	assert.Equal(t, builder, result)
	assert.Len(t, builder.modulesCtx, 1)
	assert.Contains(t, builder.modulesCtx, module)
}

func TestBuilder_Register_WhenRegistered_ThenCallsRegisterFunc(t *testing.T) {
	// Arrange
	builder := NewBuilder()
	called := false
	registerFunc := func(r Register) {
		called = true
	}

	// Act
	result := builder.Register(registerFunc)

	// Assert
	assert.Equal(t, builder, result)
	assert.True(t, called)
}

func TestBuilder_RegisterCtx_WhenRegistered_ThenCallsRegisterFunc(t *testing.T) {
	// Arrange
	builder := NewBuilder()
	ctx := context.Background()
	called := false
	registerFunc := func(ctx context.Context, r Register) {
		called = true
		assert.Equal(t, ctx, context.Background())
	}

	// Act
	result := builder.RegisterCtx(ctx, registerFunc)

	// Assert
	assert.Equal(t, builder, result)
	assert.True(t, called)
}

func TestBuilder_Build_WhenNoModules_ThenReturnsContainer(t *testing.T) {
	// Arrange
	builder := NewBuilder()

	// Act
	container, err := builder.Build()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, container)
}

func TestBuilder_BuildCtx_WhenModules_ThenRegistersServices(t *testing.T) {
	// Arrange
	builder := NewBuilder()
	ctx := context.Background()
	module := &mockModule{}
	builder.AddModule(module)

	// Act
	container, err := builder.BuildCtx(ctx)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, container)
	assert.True(t, module.registered)
}

func TestBuilder_BuildCtx_WhenModuleFails_ThenReturnsError(t *testing.T) {
	// Arrange
	builder := NewBuilder()
	ctx := context.Background()
	module := &mockModule{shouldFail: true}
	builder.AddModule(module)

	// Act
	container, err := builder.BuildCtx(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, container)
}

func TestBuilder_MustBuild_WhenSuccess_ThenReturnsContainer(t *testing.T) {
	// Arrange
	builder := NewBuilder()

	// Act
	container := builder.MustBuild()

	// Assert
	assert.NotNil(t, container)
}

func TestBuilder_MustBuildCtx_WhenSuccess_ThenReturnsContainer(t *testing.T) {
	// Arrange
	builder := NewBuilder()
	ctx := context.Background()

	// Act
	container := builder.MustBuildCtx(ctx)

	// Assert
	assert.NotNil(t, container)
}

func TestBuilder_MustBuildCtx_WhenFails_ThenPanics(t *testing.T) {
	// Arrange
	builder := NewBuilder()
	module := &mockModule{shouldFail: true}
	builder.AddModule(module)
	ctx := context.Background()

	// Act & Assert
	assert.Panics(t, func() {
		builder.MustBuildCtx(ctx)
	})
}

func TestBuildWithContainer_WhenModules_ThenRegisters(t *testing.T) {
	// Arrange
	container := NewContainer()
	module := &mockModule{}

	// Act
	err := BuildWithContainer(container, module)

	// Assert
	require.NoError(t, err)
	assert.True(t, module.registered)
}

func TestBuildWithContainerCtx_WhenModules_ThenRegisters(t *testing.T) {
	// Arrange
	container := NewContainer()
	ctx := context.Background()
	module := &mockModule{}

	// Act
	err := BuildWithContainerCtx(ctx, container, module)

	// Assert
	require.NoError(t, err)
	assert.True(t, module.registered)
}

func TestBuildWithContainer_WhenModuleFails_ThenReturnsError(t *testing.T) {
	// Arrange
	container := NewContainer()
	module := &mockModule{shouldFail: true}

	// Act
	err := BuildWithContainer(container, module)

	// Assert
	assert.Error(t, err)
}

// Mock types for testing
type mockModule struct {
	registered bool
	shouldFail bool
}

func (m *mockModule) RegisterServices(register Register) error {
	m.registered = true
	if m.shouldFail {
		return assert.AnError
	}
	return nil
}

type mockModuleWithContext struct {
	registered bool
	shouldFail bool
}

func (m *mockModuleWithContext) RegisterServicesCtx(ctx context.Context, register Register) error {
	m.registered = true
	if m.shouldFail {
		return assert.AnError
	}
	return nil
}

func TestNewContainer_WhenCreated_ThenReturnsContainer(t *testing.T) {
	// Arrange & Act
	container := NewContainer()

	// Assert
	assert.NotNil(t, container)
	assert.NotNil(t, container.Register())
	assert.NotNil(t, container.Resolver())
}

func TestContainer_Register_WhenCalled_ThenReturnsRegister(t *testing.T) {
	// Arrange
	container := NewContainer().(*container)

	// Act
	reg := container.Register()

	// Assert
	assert.NotNil(t, reg)
	assert.Equal(t, container.register, reg)
}

func TestContainer_Resolver_WhenCalled_ThenReturnsResolver(t *testing.T) {
	// Arrange
	container := NewContainer().(*container)

	// Act
	res := container.Resolver()

	// Assert
	assert.NotNil(t, res)
	assert.Equal(t, container.resolver, res)
}
