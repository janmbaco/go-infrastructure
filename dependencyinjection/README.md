# Dependency Injection

`github.com/janmbaco/go-infrastructure/v2/dependencyinjection`

The `dependencyinjection` package provides a lightweight DI container for Go applications. It supports typed registration and resolution, multiple lifetimes, parameterized providers, tenant-specific registrations and context-aware resolution.

## What It Includes

- `Builder` to compose a container from modules and ad-hoc registrations
- `Container` exposing `Register()` and `Resolver()`
- `Register` methods for type, scoped, singleton and tenant lifetimes
- generic helpers such as `RegisterSingleton[T]`, `Resolve[T]` and `ResolveCtx[T]`
- parameter injection through named arguments
- provider functions that can receive `context.Context`
- `Module` and `ModuleWithContext` abstractions

## Install

```bash
go get github.com/janmbaco/go-infrastructure/v2/dependencyinjection
```

## Mental Model

The package follows a simple flow:

1. use `Builder` to collect modules and registrations
2. build a `Container`
3. resolve dependencies through `Resolver`

```go
builder := dependencyinjection.NewBuilder()
container := builder.MustBuild()
resolver := container.Resolver()
```

The container itself registers `Container`, `Register` and `Resolver` as singletons, so they can also be resolved from providers if needed.

## Quick Start

```go
package main

import (
    di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
)

type Logger interface {
    Info(message string)
}

type logger struct{}

func (l *logger) Info(message string) {}

type UserService struct {
    logger Logger
}

func main() {
    container := di.NewBuilder().
        Register(func(r di.Register) {
            di.RegisterSingleton[Logger](r, func() Logger {
                return &logger{}
            })

            di.RegisterTypeWithParams[*UserService](r, func(log Logger) *UserService {
                return &UserService{
                    logger: log,
                }
            }, nil)
        }).
        MustBuild()

    _ = di.MustResolve[*UserService](container.Resolver())
}
```

The lower-level `Register.AsType`, `AsScope` and `AsSingleton` APIs are also available when you need explicit control over provider signatures.

## Recommended Registration Style

For most code, the generic helpers in `register_generics.go` are the easiest API to use:

```go
di.RegisterType[*Service](r, func() *Service { ... })
di.RegisterScoped[*RequestState](r, func() *RequestState { ... })
di.RegisterSingleton[Logger](r, func() Logger { ... })

di.RegisterTypeWithParams[*Repository](r, func(db *DB, log Logger) *Repository {
    return NewRepository(db, log)
}, nil)

di.RegisterTypeWithParams[*DB](r, func(dsn string, log Logger) *DB {
    return NewDB(dsn, log)
}, map[int]string{0: "dsn"})
```

When you need more control, use the underlying `Register` interface:

```go
type Register interface {
    AsType(iface, provider interface{}, argNames map[int]string)
    AsScope(iface, provider interface{}, argNames map[int]string)
    AsSingleton(iface, provider interface{}, argNames map[int]string)
    AsTenant(tenant string, iface, provider interface{}, argNames map[int]string)
    AsSingletonTenant(tenant string, iface, provider interface{}, argNames map[int]string)
    Bind(ifaceFrom, ifaceTo interface{})
}
```

`iface` is the type key being registered. A common pattern is `new(*MyType)` or `new(MyInterface)`.

## Lifetimes

The container supports three base lifetimes plus tenant variants:

- `AsType`: creates a new instance every time the dependency is resolved
- `AsScope`: reuses the same instance within a single resolution graph
- `AsSingleton`: creates the instance once and reuses it for the life of the container
- `AsTenant`: like `AsType`, but isolated by tenant key
- `AsSingletonTenant`: like `AsSingleton`, but isolated by tenant key

### How Scoped Works

The easiest way to understand `AsScope` is to think in terms of a single resolution graph.

If `OrderService` depends on `RequestContext` and `OrderRepository`, and `OrderRepository` also depends on `RequestContext`, then resolving `OrderService` creates a graph like this:

```text
Resolve(OrderService)
OrderService
├─ RequestContext
└─ OrderRepository
   └─ RequestContext
```

What changes with each lifetime is how those two `RequestContext` nodes behave:

- `AsType`: `OrderService.RequestContext` and `OrderRepository.RequestContext` are different instances
- `AsScope`: both references point to the same instance during that one resolution
- `AsSingleton`: both references point to the same instance, and that same instance is also reused in later resolutions

Example:

```go
container := di.NewBuilder().
    Register(func(r di.Register) {
        r.AsScope(new(*RequestContext), func() *RequestContext {
            return &RequestContext{}
        }, nil)

        r.AsType(new(*OrderRepository), func(ctx *RequestContext) *OrderRepository {
            return &OrderRepository{RequestContext: ctx}
        }, nil)

        r.AsType(new(*OrderService), func(ctx *RequestContext, repo *OrderRepository) *OrderService {
            return &OrderService{RequestContext: ctx, Repository: repo}
        }, nil)
    }).
    MustBuild()

service := container.Resolver().Type(new(*OrderService), nil).(*OrderService)
_ = service
```

Because `RequestContext` is scoped, the same instance is reused across the dependencies created while resolving `*OrderService`.

## Parameter Injection

Providers can receive named parameters. `argNames` maps provider argument positions to parameter names.

```go
container := di.NewBuilder().
    Register(func(r di.Register) {
        r.AsType(new(*DB), func(dsn string) *DB {
            return NewDB(dsn)
        }, map[int]string{0: "dsn"})
    }).
    MustBuild()

db := di.ResolveWithParams[*DB](container.Resolver(), map[string]interface{}{
    "dsn": "postgres://localhost/app",
})

_ = db
```

If a named parameter is not provided, the container tries to resolve that argument as another dependency.

## Context-Aware Providers

Providers may accept `context.Context` as their first argument. When you resolve with `ResolveCtx`, `Resolver.TypeCtx` or `Resolver.TenantCtx`, that context is propagated through the whole resolution graph.

```go
container := di.NewBuilder().
    Register(func(r di.Register) {
        r.AsScope(new(*RequestInfo), func(ctx context.Context) *RequestInfo {
            requestID, _ := ctx.Value("request_id").(string)
            return &RequestInfo{RequestID: requestID}
        }, nil)
    }).
    MustBuild()

ctx := context.WithValue(context.Background(), "request_id", "req-123")
info := di.ResolveCtx[*RequestInfo](ctx, container.Resolver())
_ = info
```

If resolution is attempted with a cancelled or expired context, the container panics. The same happens when a provider returns `(value, error)` and the error is non-nil.

The `AsTypeCtx`, `AsScopeCtx`, `AsSingletonCtx`, `AsTenantCtx` and `AsSingletonTenantCtx` methods exist on `Register`, but they currently delegate to their non-`Ctx` counterparts. The important runtime behavior comes from context-aware provider signatures and `*Ctx` resolution.

## Tenant Registrations

Tenant registrations let you keep separate implementations or instances per tenant key:

```go
container := di.NewBuilder().
    Register(func(r di.Register) {
        di.RegisterSingletonTenant[*Mailer](r, "acme", func() *Mailer {
            return NewMailer("smtp.acme.local")
        })

        di.RegisterSingletonTenant[*Mailer](r, "globex", func() *Mailer {
            return NewMailer("smtp.globex.local")
        })
    }).
    MustBuild()

acmeMailer := di.ResolveTenant[*Mailer](container.Resolver(), "acme")
globexMailer := di.ResolveTenant[*Mailer](container.Resolver(), "globex")

_, _ = acmeMailer, globexMailer
```

Tenant registrations are keyed by both type and tenant name.

## Binding

`Bind` lets one registration reuse another type key:

```go
container := di.NewBuilder().
    Register(func(r di.Register) {
        r.AsSingleton(new(*consoleLogger), func() *consoleLogger {
            return &consoleLogger{}
        }, nil)

        r.Bind(new(Logger), new(*consoleLogger))
    }).
    MustBuild()

logger := di.Resolve[Logger](container.Resolver())
_ = logger
```

This is useful when the concrete type is the registered provider and the interface should resolve to it.

## Modules

Modules package related registrations together:

```go
type UserModule struct{}

func (m *UserModule) RegisterServices(register di.Register) error {
    di.RegisterSingleton[Logger](register, func() Logger {
        return &logger{}
    })
    di.RegisterTypeWithParams[*UserService](register, func(log Logger) *UserService {
        return &UserService{logger: log}
    }, nil)
    return nil
}

container := di.NewBuilder().
    AddModule(&UserModule{}).
    MustBuild()
```

For modules that need context during registration, implement `ModuleWithContext` and use `AddModuleWithContext`.

```go
type SecretsModule struct{}

func (m *SecretsModule) RegisterServicesCtx(ctx context.Context, register di.Register) error {
    secret := ctx.Value("dsn").(string)
    di.RegisterSingleton[*DB](register, func() *DB {
        return NewDB(secret)
    })
    return nil
}
```

## Builder and Container APIs

`Builder`:

```go
func NewBuilder() *Builder
func (b *Builder) AddModule(module Module) *Builder
func (b *Builder) AddModules(modules ...Module) *Builder
func (b *Builder) AddModuleWithContext(module ModuleWithContext) *Builder
func (b *Builder) Register(registerFunc func(Register)) *Builder
func (b *Builder) RegisterCtx(ctx context.Context, registerFunc func(context.Context, Register)) *Builder
func (b *Builder) Build() (Container, error)
func (b *Builder) BuildCtx(ctx context.Context) (Container, error)
func (b *Builder) MustBuild() Container
func (b *Builder) MustBuildCtx(ctx context.Context) Container
```

`Container`:

```go
type Container interface {
    Register() Register
    Resolver() Resolver
}
```

`BuildWithContainer` and `BuildWithContainerCtx` are available when you want to register modules into an existing container.

## Resolution APIs

Low-level resolver:

```go
type Resolver interface {
    Type(iface interface{}, params map[string]interface{}) interface{}
    Tenant(tenantName string, iface interface{}, params map[string]interface{}) interface{}
    TypeCtx(ctx context.Context, iface interface{}, params map[string]interface{}) interface{}
    TenantCtx(ctx context.Context, tenantName string, iface interface{}, params map[string]interface{}) interface{}
}
```

Generic helpers:

- `Resolve[T](resolver)`
- `ResolveWithParams[T](resolver, params)`
- `ResolveTenant[T](resolver, tenant)`
- `ResolveTenantWithParams[T](resolver, tenant, params)`
- `ResolveCtx[T](ctx, resolver)`
- `ResolveWithParamsCtx[T](ctx, resolver, params)`
- `ResolveTenantCtx[T](ctx, resolver, tenant)`
- `ResolveTenantWithParamsCtx[T](ctx, resolver, tenant, params)`
- `MustResolve[T](resolver)`
- `MustResolveCtx[T](ctx, resolver)`
- `TryResolve[T](resolver)`
- `TryResolveCtx[T](ctx, resolver)`

The generic helpers are usually the most convenient choice for application code.

## Notes and Caveats

- Resolution failures panic. That includes missing registrations, provider errors and cancelled contexts.
- Providers must be functions.
- Scoped instances are cached only within a single resolution graph.
- Singleton instances are created on first successful resolution.
- Context-aware providers should declare `context.Context` as the first parameter.

## Related Files

- `examples/context_example.go`: end-to-end example of context propagation
- `register_generics.go`: generic registration helpers
- `resolver_generics.go`: generic resolution helpers
