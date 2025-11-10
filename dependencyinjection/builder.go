package dependencyinjection

import "context"

// Builder provides a fluent API for building a container with modules
type Builder struct {
	container  Container
	modules    []Module
	modulesCtx []ModuleWithContext
}

// NewBuilder creates a new container builder
func NewBuilder() *Builder {
	return &Builder{
		container:  NewContainer(),
		modules:    make([]Module, 0),
		modulesCtx: make([]ModuleWithContext, 0),
	}
}

// AddModule adds a module to the builder
func (b *Builder) AddModule(module Module) *Builder {
	b.modules = append(b.modules, module)
	return b
}

// AddModules adds multiple modules to the builder
func (b *Builder) AddModules(modules ...Module) *Builder {
	b.modules = append(b.modules, modules...)
	return b
}

// AddModuleWithContext adds a context-aware module
func (b *Builder) AddModuleWithContext(module ModuleWithContext) *Builder {
	b.modulesCtx = append(b.modulesCtx, module)
	return b
}

// Register allows direct registration before building
func (b *Builder) Register(registerFunc func(Register)) *Builder {
	registerFunc(b.container.Register())
	return b
}

// RegisterCtx allows direct registration with context
func (b *Builder) RegisterCtx(ctx context.Context, registerFunc func(context.Context, Register)) *Builder {
	registerFunc(ctx, b.container.Register())
	return b
}

// Build builds the container by registering all modules
func (b *Builder) Build() (Container, error) {
	return b.BuildCtx(context.Background())
}

// BuildCtx builds the container with context
func (b *Builder) BuildCtx(ctx context.Context) (Container, error) {
	register := b.container.Register()

	for _, module := range b.modules {
		if err := module.RegisterServices(register); err != nil {
			return nil, err
		}
	}

	for _, module := range b.modulesCtx {
		if err := module.RegisterServicesCtx(ctx, register); err != nil {
			return nil, err
		}
	}

	return b.container, nil
}

// MustBuild builds the container and panics on error
func (b *Builder) MustBuild() Container {
	return b.MustBuildCtx(context.Background())
}

// MustBuildCtx builds with context and panics on error
func (b *Builder) MustBuildCtx(ctx context.Context) Container {
	container, err := b.BuildCtx(ctx)
	if err != nil {
		panic(err)
	}
	return container
}

// BuildWithContainer builds using an existing container
func BuildWithContainer(container Container, modules ...Module) error {
	return BuildWithContainerCtx(context.Background(), container, modules...)
}

// BuildWithContainerCtx builds with context using an existing container
func BuildWithContainerCtx(ctx context.Context, container Container, modules ...Module) error {
	register := container.Register()

	for _, module := range modules {
		if err := module.RegisterServices(register); err != nil {
			return err
		}
	}

	return nil
}
