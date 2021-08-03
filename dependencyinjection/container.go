package dependencyinjection

// Container defines an object responsible to contains the dependencies of a application
type Container interface {
	Register() Register
	Resolver() Resolver
}

type container struct {
	register Register
	resolver Resolver
}

// Register gets the object responsible to register dependencies
func (c *container) Register() Register {
	return c.register
}

// Resolver gets the object responsilbe to resolver dependencies
func (c *container) Resolver() Resolver {
	return c.resolver
}

// NewContainer returns a container
func NewContainer() Container {
	deps := newDependencies()
	container := &container{register: newRegister(deps), resolver: newResolver(deps)}
	container.register.AsSingleton(new(Container), func() Container { return container }, nil)
	container.register.AsSingleton(new(Register), container.Register(), nil)
	container.register.AsSingleton(new(Resolver), container.Resolver(), nil)
	return container
}
