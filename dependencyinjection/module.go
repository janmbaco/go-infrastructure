package dependencyinjection
import "context"

// Module represents a self-contained unit that can register its dependencies
type Module interface {
	// RegisterServices registers all services in this module with the register
	RegisterServices(register Register) error
}

// ModuleFunc is a function adapter for Module interface
type ModuleFunc func(Register) error

// RegisterServices implements Module interface
func (f ModuleFunc) RegisterServices(register Register) error {
	return f(register)
}

// ModuleWithContext represents a module that requires context
type ModuleWithContext interface {
	// RegisterServicesCtx registers all services with context support
	RegisterServicesCtx(ctx context.Context, register Register) error
}

// ModuleFuncCtx is a function adapter for ModuleWithContext interface
type ModuleFuncCtx func(context.Context, Register) error

// RegisterServicesCtx implements ModuleWithContext interface
func (f ModuleFuncCtx) RegisterServicesCtx(ctx context.Context, register Register) error {
	return f(ctx, register)
}
