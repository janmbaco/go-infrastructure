package dependencyinjection

import (
	"fmt"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"reflect"
	"sync"
)

// Dependencies defines an object responsible to store the regiters of provider of dependencies for a application
type Dependencies interface {
	Set(key DependencyKey, object DependencyObject)
	Get(key DependencyKey) DependencyObject
	Bind(keyFrom DependencyKey, keyTo DependencyKey)
}

type DependencyObject interface {
	Create(params map[string]interface{}, dependencies Dependencies, scopeObjects map[DependencyObject]interface{}) interface{}
}

type DependencyKey struct {
	Tenant string
	Iface  reflect.Type
}

type dependencies struct {
	objects sync.Map
	binds   sync.Map
}

func newDependencies() Dependencies {
	return &dependencies{}
}

func (d *dependencies) Set(key DependencyKey, object DependencyObject) {
	errorschecker.CheckNilParameter(map[string]interface{}{"object": object})
	d.objects.Store(key, object)
}

func (d *dependencies) Get(key DependencyKey) DependencyObject {
	realKey := key
	if bind, ok := d.binds.Load(key); ok {
		realKey = bind.(DependencyKey)
	}
	if object, ok := d.objects.Load(realKey); ok {
		return object.(DependencyObject)
	}
	panic(fmt.Errorf("%v is not registered as a dependency", key.Iface.Name()))
}

func (d *dependencies) Bind(keyFrom DependencyKey, keyTo DependencyKey) {
	d.binds.Store(keyFrom, keyTo)
}



// FileConfigHandlerErrorType is the type of the errors of FileConfigHandler
type dependecyType uint8

const (
	_NewType dependecyType = iota
	_ScopedType
	_Singleton
)

type dependencyObject struct {
	provider    interface{}
	object      interface{}
	dependecyType dependecyType
	argNames    map[uint]string
}

func (do *dependencyObject) Create(params map[string]interface{}, dependencies Dependencies, scopedObjects map[DependencyObject]interface{}) interface{} {

	if do.object != nil {
		return do.object
	}

	if obj, isContained :=  scopedObjects[do]; isContained {
		return obj
	} 

	functionValue := reflect.ValueOf(do.provider)
	functionType := reflect.TypeOf(do.provider)
	if functionType.Kind() != reflect.Func {
		panic("The provider must be a Func!")
	}
	args := make([]reflect.Value, 0)
	if functionType.NumIn() > 0 {
		for i := 0; i < functionType.NumIn(); i++ {
			var name = do.argNames[uint(i)]
			if object, isInParamas := params[do.argNames[uint(i)]]; name != "" && isInParamas {
				args = append(args, reflect.ValueOf(object))
			} else {
				args = append(args, reflect.ValueOf(dependencies.Get(DependencyKey{Iface: functionType.In(i)}).Create(params, dependencies, scopedObjects)))
			}
		}
	}

	result := functionValue.Call(args)[0].Interface()
	if result != nil {
		switch do.dependecyType { 
			case _Singleton:
				do.object = result
			case _ScopedType:
				scopedObjects[do] = result
			}
	}

	return result
}
