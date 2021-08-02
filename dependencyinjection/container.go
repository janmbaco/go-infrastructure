package dependencyinjection

type Container struct {
	Register Register
	Resolver Resolver
}

func NewContainer() *Container {
	deps := newDependencies()
	return &Container{Register: newRegister(deps), Resolver: newResolver(deps)}
}
