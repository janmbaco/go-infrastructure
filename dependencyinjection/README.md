# Dependency Injection (`dependencyinjection`)

The `dependencyinjection` package provides a **type-safe DI container** for Go applications.

It gives you:

- a fluent **Builder → Container → Resolver** pipeline  
- a clean **Register / Resolver** API  
- **lifetimes**: Type, Scoped, Singleton, Tenant, SingletonTenant  
- **generics** for type-safe registration & resolution  
- **context-aware** factories (`…Ctx` variants)  
- **parameterized** factories (`WithParams` variants)  
- **tenant-based** registrations (multi-tenant support)  
- a **Module** abstraction to package your registrations

```go
import di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
````

---

## Mental model

Keep this model in mind:

1. **Builder** – configure modules and registrations.
2. **Container** – built once; holds all registrations.
3. **Resolver** – used at runtime to resolve dependencies.

```go
builder := di.NewBuilder()
// add modules and registrations...
container := builder.MustBuild()
resolver := container.Resolver()
```

You typically build the container once at startup and then pass around the `Resolver`.

---

## Quick Start

### Define a service

```go
package main

import (
    "context"

    di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
    logsIoc "github.com/janmbaco/go-infrastructure/v2/logs/ioc"
    "github.com/janmbaco/go-infrastructure/v2/logs"
)

type UserService struct {
    logger logs.Logger
}

func NewUserService(logger logs.Logger) *UserService {
    return &UserService{logger: logger}
}

func (s *UserService) Start(ctx context.Context) error {
    s.logger.Info("user service started")
    return nil
}

func main() {
    ctx := context.Background()

    container := di.NewBuilder().
        AddModule(logsIoc.NewLogsModule()).
        Register(func(r di.Register) {
            // singleton service
            di.RegisterSingleton[*UserService](r, func() *UserService {
                logger := di.Resolve[logs.Logger](r.Resolver())
                return NewUserService(logger)
            })
        }).
        MustBuild()

    svc := di.MustResolve[*UserService](container.Resolver())
    _ = svc.Start(ctx)
}
```

What this shows:

* `NewBuilder` creates a builder.
* `AddModule` wires in infrastructure modules (logs, etc.).
* `Register` lets you add your own registrations using generics.
* `MustBuild` builds or panics if something fails.
* `MustResolve[T]` gets a ready instance of `T`.

---

## Lifetimes: Type, Scoped, Singleton, Tenant, SingletonTenant

### Summary

* **Type**
  A new instance **every time** the dependency is needed.

* **Scoped (per resolution graph)**
  One instance **per resolution graph**.
  If `A` is a scoped dependency and you resolve `C` once, any `A` needed inside that single resolution (directly or via other dependencies like `B`) will be the **same instance**.

* **Singleton**
  One instance **for the entire container**.

* **Tenant / SingletonTenant**
  Separate instances for each **tenant key** (`string`), either per resolution graph (Tenant) or global (SingletonTenant).

---

### Type vs Scoped vs Singleton – concrete example

Imagine:

* `A` is a dependency of `B`
* `A` is also a dependency of `C`
* `B` is a dependency of `C`

So when you resolve `C`, the graph is:

```text
C
├─ A
└─ B
   └─ A
```

We’ll use a simple `A` with an ID:

```go
type A struct {
    ID int
}

var nextID int

func newA() *A {
    nextID++
    return &A{ID: nextID}
}

type B struct {
    A *A
}

type C struct {
    A *A
    B *B
}
```

#### 1) `A` registered as **Type**

```go
register.AsType(
    new(*A),
    func() *A {
        return newA()
    },
    nil,
)

register.AsType(
    new(*B),
    func(a *A) *B {
        return &B{A: a}
    },
    nil,
)

register.AsType(
    new(*C),
    func(a *A, b *B) *C {
        return &C{A: a, B: b}
    },
    nil,
)

c := resolver.Type(new(*C), nil).(*C)

fmt.Println(c.A.ID, c.B.A.ID)
fmt.Println("c.A == c.B.A ?", c.A == c.B.A)
```

Possible output:

```text
1 2
c.A == c.B.A ? false
```

* `Type` ⇒ each usage of `A` gets a **new instance**.
* The `A` in `C` is **different** from the `A` in `B`.

#### 2) `A` registered as **Scoped**

```go
register.AsScope(
    new(*A),
    func() *A {
        return newA()
    },
    nil,
)

// B and C can stay as AsType
register.AsType(
    new(*B),
    func(a *A) *B {
        return &B{A: a}
    },
    nil,
)

register.AsType(
    new(*C),
    func(a *A, b *B) *C {
        return &C{A: a, B: b}
    },
    nil,
)

c1 := resolver.Type(new(*C), nil).(*C)

fmt.Println(c1.A.ID, c1.B.A.ID)
fmt.Println("c1.A == c1.B.A ?", c1.A == c1.B.A)
```

Possible output:

```text
1 1
c1.A == c1.B.A ? true
```

Now:

* **Within a single resolution of `C`**, `A` is created once and reused.
* The `A` in `C` and the `A` in `B` are the **same instance**.

If you resolve `C` twice:

```go
c1 := resolver.Type(new(*C), nil).(*C)
c2 := resolver.Type(new(*C), nil).(*C)

fmt.Println("c1.A == c2.A ?", c1.A == c2.A)
fmt.Println("c1.B.A == c2.B.A ?", c1.B.A == c2.B.A)
```

You will get:

```text
c1.A == c2.A ? false
c1.B.A == c2.B.A ? false
```

> **Scoped** means “shared inside one resolution graph”, not global.

#### 3) `A` registered as **Singleton**

```go
register.AsSingleton(
    new(*A),
    func() *A {
        return newA()
    },
    nil,
)

c1 := resolver.Type(new(*C), nil).(*C)
c2 := resolver.Type(new(*C), nil).(*C)

fmt.Println(c1.A == c1.B.A) // true – shared within C1
fmt.Println(c2.A == c2.B.A) // true – shared within C2
fmt.Println(c1.A == c2.A)   // true – shared globally (singleton)
```

---

## Register API

### Low-level `Register` interface

```go
type Register interface {
    AsType(iface, provider interface{}, argNames map[int]string)
    AsScope(iface, provider interface{}, argNames map[int]string)
    AsSingleton(iface, provider interface{}, argNames map[int]string)
    AsTenant(tenant string, iface, provider interface{}, argNames map[int]string)
    AsSingletonTenant(tenant string, iface, provider interface{}, argNames map[int]string)
    Bind(ifaceFrom, ifaceTo interface{})

    // Context-aware methods
    AsTypeCtx(ctx context.Context, iface, provider interface{}, argNames map[int]string)
    AsScopeCtx(ctx context.Context, iface, provider interface{}, argNames map[int]string)
    AsSingletonCtx(ctx context.Context, iface, provider interface{}, argNames map[int]string)
    AsTenantCtx(ctx context.Context, tenant string, iface, provider interface{}, argNames map[int]string)
    AsSingletonTenantCtx(ctx context.Context, tenant string, iface, provider interface{}, argNames map[int]string)
}
```

You rarely need this directly; prefer the **generic helpers** below for type safety.

---

## Generic registration helpers

All these work on `di.Register` and use Go generics.

### Type

```go
di.RegisterType[T](r Register, factory func() T)
di.RegisterTypeCtx[T](ctx context.Context, r Register, factory func(context.Context) (T, error))
di.RegisterTypeWithParams[T](r Register, factory interface{}, argNames map[int]string)
```

### Scoped

```go
di.RegisterScoped[T](r Register, factory func() T)
di.RegisterScopedCtx[T](ctx context.Context, r Register, factory func(context.Context) (T, error))
di.RegisterScopedWithParams[T](r Register, factory interface{}, argNames map[int]string)
```

### Singleton

```go
di.RegisterSingleton[T](r Register, factory func() T)
di.RegisterSingletonCtx[T](ctx context.Context, r Register, factory func(context.Context) (T, error))
di.RegisterSingletonWithParams[T](r Register, factory interface{}, argNames map[int]string)
```

### Tenant / SingletonTenant

```go
di.RegisterTenant[T](r Register, tenant string, factory func() T)
di.RegisterTenantCtx[T](ctx context.Context, r Register, tenant string, factory func(context.Context) (T, error))
di.RegisterTenantWithParams[T](r Register, tenant string, factory interface{}, argNames map[int]string)

di.RegisterSingletonTenant[T](r Register, tenant string, factory func() T)
di.RegisterSingletonTenantCtx[T](ctx context.Context, r Register, tenant string, factory func(context.Context) (T, error))
di.RegisterSingletonTenantWithParams[T](r Register, tenant string, factory interface{}, argNames map[int]string)
```

All signatures are listed in the official docs.

---

## Parameters (`WithParams`)

When your factory needs external values (DSNs, secrets, limits…), use the `WithParams` variants.

You:

1. Declare parameter names by position (`argNames`), and
2. Provide a map at resolve time.

```go
// Register with named parameters
di.RegisterTypeWithParams[*DB](
    r,
    func(dsn string, maxOpen int) *DB {
        return OpenDB(dsn, maxOpen)
    },
    map[int]string{
        0: "dsn",
        1: "maxOpen",
    },
)

// Resolve with params
db := di.ResolveWithParams[*DB](container.Resolver(), map[string]interface{}{
    "dsn":     "postgres://user:pass@host/db",
    "maxOpen": 10,
})
```

Works for all `*WithParams` helpers: Type, Scoped, Singleton, Tenant, SingletonTenant.

---

## Context Support

### Overview

The `dependencyinjection` package provides **full context.Context support** throughout the dependency resolution chain:

- **Context propagation**: Context flows through the entire dependency graph
- **Cancellation**: Respects `ctx.Done()` and stops resolution immediately
- **Timeouts**: Honors context deadlines during dependency creation
- **Values**: Pass request-scoped data (request IDs, tenant info, tracing) through context
- **Automatic detection**: Providers can optionally accept `context.Context` as first parameter

### Context-aware factories (`…Ctx`)

If construction may block or fail (network, disk, etc.), use the `…Ctx` versions.

```go
ctx := context.Background()

di.RegisterSingletonCtx[*Client](ctx, r, func(ctx context.Context) (*Client, error) {
    client, err := NewClient(ctx) // respects deadlines / cancellations
    return client, err
})
```

Nearly every registration helper has a `…Ctx` sibling that accepts a factory `func(context.Context) (T, error)`.

### How providers receive context

Providers can **optionally** accept `context.Context` as their **first parameter**. The DI container automatically detects this and passes the resolution context:

```go
// Provider without context
di.RegisterType[*Service](r, func() *Service {
    return &Service{}
})

// Provider with context - automatically detected
di.RegisterType[*Service](r, func(ctx context.Context) *Service {
    // ctx is automatically passed during resolution
    requestID := ctx.Value("requestID")
    return &Service{RequestID: requestID}
})

// Provider with context and dependencies
di.RegisterType[*Service](r, func(ctx context.Context, logger logs.Logger, db *DB) *Service {
    // ctx comes first, then dependencies
    return &Service{
        Logger: logger,
        DB:     db,
        RequestID: ctx.Value("requestID"),
    }
})
```

### Context propagation through nested dependencies

Context automatically propagates through the entire dependency graph:

```go
type Database struct {
    RequestID string
}

type Repository struct {
    DB *Database
}

type Service struct {
    Repo *Repository
}

// All providers can access the same context
di.RegisterType[*Database](r, func(ctx context.Context) *Database {
    return &Database{RequestID: ctx.Value("requestID").(string)}
})

di.RegisterType[*Repository](r, func(ctx context.Context, db *Database) *Repository {
    // ctx is the same one passed to Database
    return &Repository{DB: db}
})

di.RegisterType[*Service](r, func(ctx context.Context, repo *Repository) *Service {
    // ctx propagates through the entire chain
    return &Service{Repo: repo}
})

// Resolve with context
ctx := context.WithValue(context.Background(), "requestID", "req-12345")
service := di.ResolveCtx[*Service](ctx, resolver)
// service.Repo.DB.RequestID == "req-12345"
```

### Cancellation and timeouts

The container checks for context cancellation at each resolution step:

```go
// Example: timeout during resolution
ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()

di.RegisterType[*SlowService](r, func(ctx context.Context) (*SlowService, error) {
    time.Sleep(200 * time.Millisecond) // exceeds timeout
    return &SlowService{}, nil
})

// This will panic with "context error during dependency resolution: context deadline exceeded"
service := di.ResolveCtx[*SlowService](ctx, resolver)
```

```go
// Example: cancellation
ctx, cancel := context.WithCancel(context.Background())

go func() {
    time.Sleep(50 * time.Millisecond)
    cancel() // Cancel during resolution
}()

// Will panic if cancellation happens during resolution
service := di.ResolveCtx[*ComplexService](ctx, resolver)
```

### Using context values for request-scoped data

Pass request-specific data through context:

```go
type contextKey string

const (
    requestIDKey contextKey = "requestID"
    tenantIDKey  contextKey = "tenantID"
    userIDKey    contextKey = "userID"
)

// In your HTTP handler
func handleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    ctx = context.WithValue(ctx, requestIDKey, generateRequestID())
    ctx = context.WithValue(ctx, tenantIDKey, extractTenant(r))
    ctx = context.WithValue(ctx, userIDKey, extractUser(r))
    
    // Resolve request-scoped services with context
    service := di.ResolveCtx[*RequestService](ctx, resolver)
    service.Process()
}

// Service can access request context
di.RegisterScoped[*RequestService](r, func(ctx context.Context, logger logs.Logger) *RequestService {
    requestID := ctx.Value(requestIDKey).(string)
    tenantID := ctx.Value(tenantIDKey).(string)
    
    logger.WithFields(map[string]interface{}{
        "requestID": requestID,
        "tenantID":  tenantID,
    }).Info("creating request service")
    
    return &RequestService{
        RequestID: requestID,
        TenantID:  tenantID,
        Logger:    logger,
    }
})
```

### Methods without context

For backward compatibility and convenience, methods without `Ctx` suffix use `context.Background()` internally:

```go
// These two are equivalent when you don't need context features:
service1 := di.Resolve[*Service](resolver)
service2 := di.ResolveCtx[*Service](context.Background(), resolver)
```

But if your provider accepts `context.Context`, even `Resolve()` will pass `context.Background()` to it:

```go
di.RegisterType[*Service](r, func(ctx context.Context) *Service {
    // ctx will be context.Background() when using Resolve()
    // ctx will be your custom context when using ResolveCtx()
    return &Service{}
})

service1 := di.Resolve[*Service](resolver)           // receives context.Background()
service2 := di.ResolveCtx[*Service](ctx, resolver)   // receives your ctx
```

### Best practices

1. **Use context for request-scoped operations**: HTTP requests, gRPC calls, database transactions
2. **Pass timeouts for expensive operations**: Network calls, external APIs, database operations
3. **Propagate tracing/logging context**: Request IDs, correlation IDs, trace spans
4. **Don't use context for application configuration**: Use parameters or separate config system
5. **Singleton caveat**: Singletons created with context only use it once (at creation time)

---

## Tenants

Tenants let you host multiple variants of the same type under different keys (e.g., different SMTP servers, DBs, etc.).

```go
di.RegisterTenant[*Mailer](r, "acme", func() *Mailer {
    return NewMailer("smtp.acme")
})

di.RegisterTenant[*Mailer](r, "globex", func() *Mailer {
    return NewMailer("smtp.globex")
})

acmeMailer := di.ResolveTenant[*Mailer](container.Resolver(), "acme")
globexMailer := di.ResolveTenant[*Mailer](container.Resolver(), "globex")
```

For global shared instances per tenant, use `RegisterSingletonTenant`.

There are also `WithParams` and `…Ctx` variants for more advanced scenarios.

---

## Resolver API

### Low-level `Resolver` interface

```go
type Resolver interface {
    Type(iface interface{}, params map[string]interface{}) interface{}
    Tenant(tenantName string, iface interface{}, params map[string]interface{}) interface{}

    // Context-aware
    TypeCtx(ctx context.Context, iface interface{}, params map[string]interface{}) interface{}
    TenantCtx(ctx context.Context, tenantName string, iface interface{}, params map[string]interface{}) interface{}
}
```

This is the non-generic core that all the helper functions build on.

### Generic resolution helpers

* `Resolve[T](resolver Resolver) T`
* `ResolveCtx[T](ctx context.Context, resolver Resolver) T`
* `ResolveWithParams[T](resolver Resolver, params map[string]interface{}) T`
* `ResolveWithParamsCtx[T](ctx context.Context, resolver Resolver, params map[string]interface{}) T`
* `ResolveTenant[T](resolver Resolver, tenant string) T`
* `ResolveTenantCtx[T](ctx context.Context, resolver Resolver, tenant string) T`
* `ResolveTenantWithParams[T](resolver Resolver, tenant string, params map[string]interface{}) T`
* `ResolveTenantWithParamsCtx[T](ctx context.Context, resolver Resolver, tenant string, params map[string]interface{}) T`

### Safety helpers

Choose how strict you want to be:

* `MustResolve[T](resolver)` – panics if not registered.
* `MustResolveCtx[T](ctx, resolver)` – same, but context-aware.
* `TryResolve[T](resolver) (T, bool)` – returns `(zero, false)` if missing.
* `TryResolveCtx[T](ctx, resolver) (T, bool)` – same with context.
* `ResolveOrDefault[T](resolver, defaultValue T) T` – fallback to a default.

---

## Builder & Container

### Builder

```go
builder := di.NewBuilder().
    AddModule(logsModule).
    AddModule(errorsModule).
    Register(func(r di.Register) {
        di.RegisterSingleton[*MyService](r, NewMyService)
    })
```

`Builder` methods:

```go
type Builder struct { /* ... */ }

func NewBuilder() *Builder

func (b *Builder) AddModule(module Module) *Builder
func (b *Builder) AddModuleWithContext(module ModuleWithContext) *Builder
func (b *Builder) AddModules(modules ...Module) *Builder

func (b *Builder) Register(registerFunc func(Register)) *Builder
func (b *Builder) RegisterCtx(ctx context.Context, registerFunc func(context.Context, Register)) *Builder

func (b *Builder) Build() (Container, error)
func (b *Builder) BuildCtx(ctx context.Context) (Container, error)
func (b *Builder) MustBuild() Container
func (b *Builder) MustBuildCtx(ctx context.Context) Container
```

You usually use `MustBuild()` for app startup.

### Container

```go
type Container interface {
    Register() Register
    Resolver() Resolver
}

func NewContainer() Container
```

You typically get a `Container` from `NewBuilder().MustBuild()`.
`NewContainer` is available for advanced scenarios (manually building with `BuildWithContainer`).

---

## Modules

Modules group registrations logically.

```go
type Module interface {
    // RegisterServices registers all services in this module with the register
    RegisterServices(register Register) error
}
```

Example:

```go
type UserModule struct{}

func (m *UserModule) RegisterServices(r di.Register) error {
    di.RegisterSingleton[*UserRepository](r, NewUserRepository)
    di.RegisterType[*UserService](r, NewUserService)
    return nil
}
```

Usage:

```go
container := di.NewBuilder().
    AddModule(&UserModule{}).
    MustBuild()
```

### Functional modules

You don’t need a struct; you can use function adapters:

```go
var userModule di.Module = di.ModuleFunc(func(r di.Register) error {
    di.RegisterSingleton[*UserRepository](r, NewUserRepository)
    return nil
})

var userModuleCtx di.ModuleWithContext = di.ModuleFuncCtx(
    func(ctx context.Context, r di.Register) error {
        di.RegisterSingletonCtx[*Client](ctx, r, newClientFromCtx)
        return nil
    },
)
```

Then:

```go
di.NewBuilder().
    AddModule(userModule).
    AddModuleWithContext(userModuleCtx).
    MustBuild()
```

All types are exposed in the package docs.

---

## Using an existing container

For advanced composition, you can build on top of an existing container:

```go
err := di.BuildWithContainer(container, moduleA, moduleB)
if err != nil {
    // handle error
}
```

Context-aware variant:

```go
err := di.BuildWithContainerCtx(ctx, container, moduleA, moduleB)
```

This is useful when you need to extend a container that was created elsewhere.

---

## Advanced types

For most applications you only need:

* `Builder`, `Container`, `Register`, `Resolver`
* `Module`, `ModuleWithContext`, `ModuleFunc`, `ModuleFuncCtx`
* Generic registration/resolution helpers.

For more advanced scenarios the package exposes:

```go
type Dependencies interface {
    Set(key DependencyKey, object DependencyObject)
    Get(key DependencyKey) DependencyObject
    Bind(keyFrom DependencyKey, keyTo DependencyKey)
}

type DependencyKey struct {
    Iface  reflect.Type
    Tenant string
}

type DependencyObject interface {
    Create(params map[string]interface{}, dependencies Dependencies, scopeObjects map[DependencyObject]interface{}) interface{}
}
```

These are the low-level building blocks that back the DI container. Most users never need them directly.

---

## Testing

Because constructors are plain Go functions, you can unit test your services without the container:

```go
type fakeLogger struct{}

func (fakeLogger) Info(msg string, args ...any) {}
// implement other methods as no-ops...

func TestUserService(t *testing.T) {
    svc := NewUserService(fakeLogger{})
    // assert behaviour...
}
```

For integration-style tests, you can spin up a small container:

```go
container := di.NewBuilder().
    Register(func(r di.Register) {
        di.RegisterSingleton[*FakeMailer](r, NewFakeMailer)
        di.RegisterSingleton[*ServiceUnderTest](r, NewServiceUnderTest)
    }).
    MustBuild()

svc := di.MustResolve[*ServiceUnderTest](container.Resolver())
// run integration test with the real DI wiring
```

---

## API cheat sheet

**Core types**

* `Builder`
* `Container`
* `Register`
* `Resolver`
* `Module`, `ModuleWithContext`, `ModuleFunc`, `ModuleFuncCtx`
* `Dependencies`, `DependencyKey`, `DependencyObject`

**Builder**

* `NewBuilder()`
* `(*Builder).AddModule`, `AddModuleWithContext`, `AddModules`
* `(*Builder).Register`, `RegisterCtx`
* `(*Builder).Build`, `BuildCtx`, `MustBuild`, `MustBuildCtx`

**Container**

* `NewContainer()`
* `Container.Register()`
* `Container.Resolver()`

**Registration (generics)**

* `RegisterType`, `RegisterTypeCtx`, `RegisterTypeWithParams`
* `RegisterScoped`, `RegisterScopedCtx`, `RegisterScopedWithParams`
* `RegisterSingleton`, `RegisterSingletonCtx`, `RegisterSingletonWithParams`
* `RegisterTenant`, `RegisterTenantCtx`, `RegisterTenantWithParams`
* `RegisterSingletonTenant`, `RegisterSingletonTenantCtx`, `RegisterSingletonTenantWithParams`

**Resolution (generics)**

* `Resolve`, `ResolveCtx`
* `ResolveWithParams`, `ResolveWithParamsCtx`
* `ResolveTenant`, `ResolveTenantCtx`
* `ResolveTenantWithParams`, `ResolveTenantWithParamsCtx`
* `MustResolve`, `MustResolveCtx`
* `TryResolve`, `TryResolveCtx`
* `ResolveOrDefault`

**Build on existing container**

* `BuildWithContainer`, `BuildWithContainerCtx`

---

With this DI module you:

* keep all wiring in one place (the container and modules),
* fully control lifetimes (Type, Scoped, Singleton, Tenant, SingletonTenant),
* get type-safe registration and resolution via generics,
* and stay idiomatic Go (plain structs, constructors, and explicit wiring).
