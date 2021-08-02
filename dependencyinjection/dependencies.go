package dependencyinjection

import (
	"errors"
	"fmt"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"reflect"
	"sync"
)

type Dependencies interface {
	Set(iface, provider interface{}, argsName map[uint]string)
	SetTenant(tenant string, iface, provider interface{}, argsName map[uint]string)
	SetAsSingleton(iface, provider interface{}, argsName map[uint]string)
	SetTenantAsSingleton(tenant string, iface, provider interface{}, argsName map[uint]string)
	Bind(ifaceFrom, ifaceTo interface{})
	Get(iface interface{}, params map[string]interface{}) interface{}
	GetTenant(tenant string, iface interface{}, params map[string]interface{}) interface{}
}

type dependencies struct {
	objects sync.Map
	binds   sync.Map
}

func newDependencies() Dependencies {
	return &dependencies{}
}

func (d *dependencies) Set(iface, provider interface{}, argNames map[uint]string) {
	errorschecker.CheckNilParameter(map[string]interface{}{"iface": iface, "provider": provider})
	d.objects.Store(&dependencyKey{iface: iface}, &dependencyObject{provider: provider, dependencies: d, argNames: argNames})
}

func (d *dependencies) SetTenant(tenant string, iface, provider interface{}, argNames map[uint]string) {
	errorschecker.CheckNilParameter(map[string]interface{}{"iface": iface, "provider": provider})
	d.objects.Store(&dependencyKey{tenant: tenant, iface: iface}, &dependencyObject{provider: provider, dependencies: d, argNames: argNames})
}

func (d *dependencies) SetAsSingleton(iface, provider interface{}, argNames map[uint]string) {
	errorschecker.CheckNilParameter(map[string]interface{}{"iface": iface, "provider": provider})
	d.objects.Store(&dependencyKey{iface: iface}, &dependencyObject{provider: provider, isSingleton: true, dependencies: d, argNames: argNames})
}

func (d *dependencies) SetTenantAsSingleton(tenant string, iface, provider interface{}, argNames map[uint]string) {
	errorschecker.CheckNilParameter(map[string]interface{}{"iface": iface, "provider": provider})
	d.objects.Store(&dependencyKey{tenant: tenant, iface: iface}, &dependencyObject{provider: provider, isSingleton: true, dependencies: d, argNames: argNames})
}

func (d *dependencies) Bind(ifaceFrom, ifaceTo interface{}) {
	errorschecker.CheckNilParameter(map[string]interface{}{"ifaceFrom": ifaceFrom, "ifaceTo": ifaceTo})
	d.binds.Store(reflect.Indirect(reflect.ValueOf(ifaceFrom)).Type(), reflect.Indirect(reflect.ValueOf(ifaceTo)).Type())
}

func (d *dependencies) Get(iface interface{}, params map[string]interface{}) interface{} {
	errorschecker.CheckNilParameter(map[string]interface{}{"iface": iface})
	return d.getAutorResolveObject("", d.getType(reflect.Indirect(reflect.ValueOf(iface)).Type()), params)

}

func (d *dependencies) GetTenant(tenant string, iface interface{}, params map[string]interface{}) interface{} {
	errorschecker.CheckNilParameter(map[string]interface{}{"iface": iface})
	return d.getAutorResolveObject(tenant, d.getType(reflect.Indirect(reflect.ValueOf(iface)).Type()), params)
}

func (d *dependencies) getType(iface reflect.Type) reflect.Type {
	if ifaceBind, ok := d.binds.Load(iface); ok {
		return ifaceBind.(reflect.Type)
	}
	return iface
}

func (d *dependencies) getAutorResolveObject(tenant string, typ reflect.Type, params map[string]interface{}) interface{} {
	var object *dependencyObject
	d.objects.Range(func(key, value interface{}) bool {
		storeTenant := key.(*dependencyKey).tenant
		storeIface := key.(*dependencyKey).iface
		if storeTenant == tenant && reflect.Indirect(reflect.ValueOf(storeIface)).Type() == typ {
			object = value.(*dependencyObject)
			return false
		}
		return true
	})
	if object == nil {
		panic(errors.New(fmt.Sprintf("%v is not registered as a dependency", typ.Name())))
	}

	return object.AutoResolve(params)
}

type dependencyKey struct {
	tenant string
	iface  interface{}
}

type dependencyObject struct {
	provider     interface{}
	object       interface{}
	isSingleton  bool
	argNames     map[uint]string
	dependencies *dependencies
}

func (do *dependencyObject) AutoResolve(params map[string]interface{}) interface{} {
	if do.isSingleton && do.object != nil {
		return do.object
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
				args = append(args, reflect.ValueOf(do.dependencies.getAutorResolveObject("", do.dependencies.getType(functionType.In(i)), params)))
			}
		}
	}

	result := functionValue.Call(args)[0].Interface()
	if result != nil && do.isSingleton {
		do.object = result
	}
	return result
}
