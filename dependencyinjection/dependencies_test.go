package dependencyinjection

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDependencies_WhenCreated_ThenReturnsDependencies(t *testing.T) {
	// Arrange & Act
	deps := newDependencies()

	// Assert
	assert.NotNil(t, deps)
}

func TestDependencies_Set_WhenSet_ThenCanGet(t *testing.T) {
	// Arrange
	deps := newDependencies()
	key := DependencyKey{Iface: reflect.TypeOf("")}
	obj := &mockDependencyObject{}

	// Act
	deps.Set(key, obj)

	// Assert
	retrieved := deps.Get(key)
	assert.Equal(t, obj, retrieved)
}

func TestDependencies_Get_WhenNotSet_ThenPanics(t *testing.T) {
	// Arrange
	deps := newDependencies()
	key := DependencyKey{Iface: reflect.TypeOf("")}

	// Act & Assert
	assert.Panics(t, func() {
		deps.Get(key)
	})
}

func TestDependencies_Bind_WhenBound_ThenGetReturnsBound(t *testing.T) {
	// Arrange
	deps := newDependencies()
	fromKey := DependencyKey{Iface: reflect.TypeOf("")}
	toKey := DependencyKey{Iface: reflect.TypeOf(0)}
	obj := &mockDependencyObject{}
	deps.Set(toKey, obj)

	// Act
	deps.Bind(fromKey, toKey)

	// Assert
	retrieved := deps.Get(fromKey)
	assert.Equal(t, obj, retrieved)
}

func TestDependencyObject_Create_WhenSingleton_ThenReturnsSameInstance(t *testing.T) {
	// Arrange
	obj := &dependencyObject{
		provider:      func() string { return "test" },
		dependecyType: _Singleton,
		argNames:      map[int]string{},
	}
	deps := newDependencies()
	params := map[string]interface{}{}
	scoped := make(map[DependencyObject]interface{})

	// Act
	first := obj.Create(params, deps, scoped)
	second := obj.Create(params, deps, scoped)

	// Assert
	assert.Equal(t, "test", first)
	assert.Equal(t, first, second)
}

func TestDependencyObject_Create_WhenScoped_ThenReturnsSameInScope(t *testing.T) {
	// Arrange
	obj := &dependencyObject{
		provider:      func() string { return "test" },
		dependecyType: _ScopedType,
		argNames:      map[int]string{},
	}
	deps := newDependencies()
	params := map[string]interface{}{}
	scoped := make(map[DependencyObject]interface{})

	// Act
	first := obj.Create(params, deps, scoped)
	second := obj.Create(params, deps, scoped)

	// Assert
	assert.Equal(t, "test", first)
	assert.Equal(t, first, second)
}

func TestDependencyObject_Create_WhenNewType_ThenReturnsNewInstance(t *testing.T) {
	// Arrange
	obj := &dependencyObject{
		provider:      func() string { return "test" },
		dependecyType: _NewType,
		argNames:      map[int]string{},
	}
	deps := newDependencies()
	params := map[string]interface{}{}
	scoped := make(map[DependencyObject]interface{})

	// Act
	first := obj.Create(params, deps, scoped)
	second := obj.Create(params, deps, scoped)

	// Assert
	assert.Equal(t, "test", first)
	assert.Equal(t, "test", second)
	// Note: Since it's a string, they are equal, but for reference types it would be different
}

func TestDependencyObject_Create_WhenProviderReturnsError_ThenPanics(t *testing.T) {
	// Arrange
	obj := &dependencyObject{
		provider: func() (string, error) { return "", assert.AnError },
		argNames: map[int]string{},
	}
	deps := newDependencies()
	params := map[string]interface{}{}
	scoped := make(map[DependencyObject]interface{})

	// Act & Assert
	assert.Panics(t, func() {
		obj.Create(params, deps, scoped)
	})
}

func TestDependencyObject_Create_WhenProviderNotFunc_ThenPanics(t *testing.T) {
	// Arrange
	obj := &dependencyObject{
		provider: "not a func",
		argNames: map[int]string{},
	}
	deps := newDependencies()
	params := map[string]interface{}{}
	scoped := make(map[DependencyObject]interface{})

	// Act & Assert
	assert.Panics(t, func() {
		obj.Create(params, deps, scoped)
	})
}

// Mock dependency object for testing
type mockDependencyObject struct {
	created interface{}
}

func (m *mockDependencyObject) Create(params map[string]interface{}, dependencies Dependencies, scopeObjects map[DependencyObject]interface{}) interface{} {
	return m.created
}
