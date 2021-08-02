package dependencyinjection

import "github.com/janmbaco/go-infrastructure/errors/errorschecker"

type Resolver interface {
	Type(iface interface{}, params map[string]interface{}) interface{}
	Tenant(tenantName string, iface interface{}, params map[string]interface{}) interface{}
}

type resolver struct {
	dependencies Dependencies
}

func newResolver(dependencies Dependencies) Resolver {
	return &resolver{dependencies: dependencies}
}

func (r *resolver) Type(iface interface{}, params map[string]interface{}) interface{} {
	errorschecker.CheckNilParameter(map[string]interface{}{"iface": iface})
	return r.dependencies.Get(iface, params)
}

func (r *resolver) Tenant(tenantName string, iface interface{}, params map[string]interface{}) interface{} {
	errorschecker.CheckNilParameter(map[string]interface{}{"iface": iface})
	return r.dependencies.GetTenant(tenantName, iface, params)
}
