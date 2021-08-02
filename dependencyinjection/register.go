package dependencyinjection

import (
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"sync"
)

type Register interface {
	AsType(iface, provider interface{}, argNames map[uint]string)
	AsSingleton(iface, provider interface{}, argNames map[uint]string)
	AsTenant(tenant string, iface, provider interface{}, argNames map[uint]string)
	AsSingletonTenant(tenant string, iface, provider interface{}, argNames map[uint]string)
	Bind(ifaceFrom, ifaceTo interface{})
}

type register struct {
	dependencies Dependencies
	binds        sync.Map
}

func newRegister(dependencies Dependencies) Register {
	return &register{dependencies: dependencies}
}

func (r *register) AsType(iface, provider interface{}, argNames map[uint]string) {
	errorschecker.CheckNilParameter(map[string]interface{}{"iface": iface, "provider": provider})
	r.dependencies.Set(iface, provider, argNames)
}

func (r *register) AsSingleton(iface, provider interface{}, argNames map[uint]string) {
	errorschecker.CheckNilParameter(map[string]interface{}{"iface": iface, "provider": provider})
	r.dependencies.SetAsSingleton(iface, provider, argNames)
}

func (r *register) AsTenant(tenant string, iface, provider interface{}, argNames map[uint]string) {
	errorschecker.CheckNilParameter(map[string]interface{}{"iface": iface, "provider": provider})
	r.dependencies.SetTenant(tenant, iface, provider, argNames)
}

func (r *register) AsSingletonTenant(tenant string, iface, provider interface{}, argNames map[uint]string) {
	errorschecker.CheckNilParameter(map[string]interface{}{"iface": iface, "provider": provider})
	r.dependencies.SetTenantAsSingleton(tenant, iface, provider, argNames)
}

func (r *register) Bind(ifaceFrom, ifaceTo interface{}) {
	errorschecker.CheckNilParameter(map[string]interface{}{"ifaceFrom": ifaceFrom, "ifaceTo": ifaceTo})
	r.dependencies.Bind(ifaceFrom, ifaceTo)
}
