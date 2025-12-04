# Go Infrastructure

[![Go Report Card](https://goreportcard.com/badge/github.com/janmbaco/go-infrastructure/v2)](https://goreportcard.com/report/github.com/janmbaco/go-infrastructure/v2)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Sponsor](https://img.shields.io/badge/Sponsor-♥-ff69b4)](https://github.com/sponsors/janmbaco)

**Production-ready infrastructure components for Go applications.**  

`go-infrastructure` provides battle-tested, modular components to build robust Go services: dependency injection with context support, live-reloading configuration, structured logging, comprehensive error handling, type-safe events, HTTP/HTTPS servers, database persistence, and more.

## Why go-infrastructure?

**Type-safe** - Leverages Go generics for compile-time safety  
**Modular** - Use only what you need, no forced dependencies  
**Production-ready** - Battle-tested patterns and error handling  
**Context-aware** - First-class support for `context.Context`  
**Well-documented** - Comprehensive docs and examples  
**Actively maintained** - Regular updates and security patches

---

## Features

### Core Infrastructure

| Module | Description | Key Features |
|--------|-------------|--------------|
| **[Dependency Injection](./dependencyinjection)** | Type-safe DI container | Context support, lifetimes (Singleton/Scoped/Transient), generics, multi-tenant |
| **[Configuration](./configuration)** | Dynamic configuration | Live reload, file-based, change events, freeze/restore |
| **[Logging](./logs)** | Structured logging | Multiple levels, console/file output, daily rotation |
| **[Error Handling](./errors)** | Centralized error management | Error catching, validation, try-catch patterns |
| **[Events Manager](./eventsmanager)** | Type-safe pub/sub | Generic-based, parallel/sequential, no reflection |
| **[Server](./server)** | HTTP/HTTPS helpers | Listener builder, SPA support, graceful shutdown |
| **[Persistence](./persistence)** | Database abstraction | GORM integration, typed access, MySQL/PostgreSQL/SQLite/SQL Server |
| **[Crypto](./crypto)** | Encryption utilities | AES-256, secure key management |
| **[Disk](./disk)** | File system utilities | File watching, change notifications, path helpers |

---

## Installation

```bash
go get github.com/janmbaco/go-infrastructure/v2
```

**Requirements:**
- Go 1.24 or higher
- Go modules enabled

---

## Quick Start

A minimal example using a DI container and file-based configuration with automatic reload:

```go
package main

import (
    "fmt"

    di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
    cfgioc "github.com/janmbaco/go-infrastructure/v2/configuration/fileconfig/ioc"
    cfgresolver "github.com/janmbaco/go-infrastructure/v2/configuration/fileconfig/ioc/resolver"
    logsioc "github.com/janmbaco/go-infrastructure/v2/logs/ioc"
    errsioc "github.com/janmbaco/go-infrastructure/v2/errors/ioc"
    eventsioc "github.com/janmbaco/go-infrastructure/v2/eventsmanager/ioc"
    diskioc "github.com/janmbaco/go-infrastructure/v2/disk/ioc"
)

type Config struct {
    Database string `json:"database"`
    Port     int    `json:"port"`
}

func main() {
    container := di.NewBuilder().
        AddModule(logsioc.NewLogsModule()).
        AddModule(errsioc.NewErrorsModule()).
        AddModule(eventsioc.NewEventsModule()).
        AddModule(diskioc.NewDiskModule()).
        AddModule(cfgioc.NewConfigurationModule()).
        MustBuild()

    resolver := container.Resolver()
    handler := cfgresolver.GetFileConfigHandler(
        resolver,
        "config.json",
        &Config{Port: 8080}, // defaults if file is missing
    )

    cfg := handler.GetConfig().(*Config)
    fmt.Printf("App running on port %d\n", cfg.Port)

    // The process keeps running and picks up changes to config.json automatically.
    select {}
}
```

`config.json`:

```json
{
  "database": "postgres://localhost/myapp",
  "port": 8080
}
```

---

## Module Deep Dive

### Dependency Injection

Full-featured DI container with Go generics and context support.

```go
container := di.NewBuilder().
    Register(func(r di.Register) {
        // Context-aware provider
        r.AsScope(new(*Service), func(ctx context.Context, logger logs.Logger) *Service {
            return &Service{
                Logger: logger,
                RequestID: ctx.Value("requestID").(string),
            }
        }, nil)
    }).
    MustBuild()

// Resolve with context
ctx := context.WithValue(context.Background(), "requestID", "req-123")
service := di.ResolveCtx[*Service](ctx, container.Resolver())
```

**Features:** Lifetimes (Type/Scoped/Singleton/Tenant), context propagation, automatic cancellation, generic helpers  
[Full Documentation](./dependencyinjection/README.md) | [Context Example](./dependencyinjection/examples/context_example.go)

### Configuration

File-based configuration with live reload and zero downtime.

```go
type Config struct {
    Database string `json:"database"`
    Port     int    `json:"port"`
}

handler := cfgresolver.GetFileConfigHandler(resolver, "config.json", &Config{Port: 8080})
cfg := handler.GetConfig().(*Config)

// Changes to config.json are automatically detected and applied
```

**Features:** Auto-reload, change events, freeze/restore, validation  
[Full Documentation](./configuration/README.md)

### Logging

Structured logging with multiple outputs and daily rotation.

```go
logger := logs.NewLogger()
logger.SetConsoleLevel(logs.InfoLevel)
logger.SetFileLogLevel(logs.TraceLevel)

logger.Info("Application started", "version", "2.0.0")
logger.WithFields(map[string]interface{}{
    "userID": "123",
    "action": "login",
}).Info("User action")
```

**Features:** Multiple levels, console/file output, daily rotation, structured fields  
[Full Documentation](./logs/README.md)

### Error Handling

Comprehensive error handling with try-catch patterns.

```go
errorCatcher := errors.NewErrorCatcher(logger)

err := errorCatcher.TryCatchError(
    func() error {
        return riskyOperation()
    },
    func(err error) error {
        logger.Error("Operation failed", err)
        return fmt.Errorf("wrapped: %w", err)
    },
)
```

**Features:** Error catching, validation, finally blocks, error wrapping  
[Full Documentation](./errors/README.md)

### Events Manager

Type-safe publish-subscribe with generics.

```go
type UserCreatedEvent struct {
    UserID string
    Email  string
}

eventManager := eventsmanager.NewEventManager()
publisher := eventsmanager.NewPublisher[UserCreatedEvent](eventManager)

subscriptions := eventsmanager.NewSubscriptions[UserCreatedEvent](errorHandler)
subscriptions.Subscribe(func(event UserCreatedEvent) {
    fmt.Printf("User created: %s\n", event.Email)
})

publisher.Publish(UserCreatedEvent{UserID: "123", Email: "user@example.com"})
```

**Features:** Type-safe, parallel/sequential execution, no reflection overhead  
[Full Documentation](./eventsmanager/README.md)

### Server

HTTP/HTTPS server with SPA support and graceful shutdown.

```go
listener := server.NewListenerBuilder().
    SetPort(":8080").
    SetTLSConfig(tlsConfig).
    SetHandler(myHandler).
    Build()

if err := listener.ListenAndServe(); err != nil {
    log.Fatal(err)
}
```

**Features:** TLS support, SPA routing, graceful shutdown, port conflict recovery  
[Full Documentation](./server/README.md)

### Persistence

Type-safe database access with GORM.

```go
db := persistence.NewDB(dbInfo, errorCatcher)
dataAccess := dataaccess.NewDataAccess[User](db)

users, err := dataAccess.SelectRows(&User{Active: true})
err = dataAccess.InsertRow(&User{Name: "John"})
err = dataAccess.UpdateRow(&User{ID: 1}, map[string]interface{}{"Name": "Jane"})
```

**Features:** Typed operations, multiple DB backends, preloading, associations  
[Full Documentation](./persistence/README.md)

---

## Real-World Example: Single Page Application Server

Complete example demonstrating all modules working together.

**Features:**
- Serves static files with SPA routing
- Live configuration reload (no restarts needed)
- Structured logging with rotation
- Graceful shutdown and error recovery
- Docker support for containerized deployments

```bash
# Run with Go
go run ./cmd/singlepageapp -port :8080 -static ./dist -index index.html

# Or with Docker
docker build -f server/facades/Dockerfile -t myapp .
docker run -p 8080:8080 myapp
```

[View Full Example](./cmd/singlepageapp) | [Dockerfile](./server/facades/Dockerfile)

---

## Security

We take security seriously. If you discover a security vulnerability:

- **DO NOT** open a public issue
- Report it via [GitHub Security Advisories](https://github.com/janmbaco/go-infrastructure/security/advisories/new)
- See our [Security Policy](./SECURITY.md) for details

---

## Contributing

We welcome contributions! Whether it's:

- Bug reports
- Feature requests
- Documentation improvements
- Code contributions

**Getting Started:**

1. Read the [Contributing Guide](./CONTRIBUTING.md)
2. Check the [Style Guide](./STYLE_GUIDE.md)
3. Review the [Code of Conduct](./CODE_OF_CONDUCT.md)

```bash
# Fork and clone
git clone https://github.com/YOUR-USERNAME/go-infrastructure.git

# Create a feature branch
git checkout -b feature/amazing-feature

# Make changes and run tests
go test ./...

# Submit a PR
```

---

## Documentation

- [CHANGELOG](./CHANGELOG.md) - Version history and changes
- [UPGRADE GUIDE (v2.0)](./UPGRADE-2.0.md) - Migration from v1.x to v2.x
- [SECURITY](./SECURITY.md) - Security policy and best practices
- [STYLE GUIDE](./STYLE_GUIDE.md) - Coding standards and conventions
- [CONTRIBUTING](./CONTRIBUTING.md) - How to contribute
- [CODE OF CONDUCT](./CODE_OF_CONDUCT.md) - Community guidelines

---

## Support & Sponsorship

If `go-infrastructure` helps you build better Go services, consider supporting its development:

[![Sponsor](https://img.shields.io/badge/Sponsor-♥-ff69b4?style=for-the-badge)](https://github.com/sponsors/janmbaco)

**Your sponsorship helps:**
- Maintain and improve the project
- Create more documentation and examples
- Fix bugs and security issues faster
- Develop new features

Sponsors will be recognized here and in release notes.

---

## Project Stats

- **Version:** v2.0.0+
- **Go Version:** 1.24+
- **License:** Apache 2.0
- **Test Coverage:** 80%+
- **Actively Maintained:** Yes

---

## Acknowledgments

Built with care for the Go community.

Special thanks to all [contributors](https://github.com/janmbaco/go-infrastructure/graphs/contributors) who have helped improve this project.

---

## License

This project is licensed under the **Apache License 2.0** – see the [LICENSE](./LICENSE) file for details.

---

## Contact & Community

- [GitHub Discussions](https://github.com/janmbaco/go-infrastructure/discussions) - Ask questions, share ideas
- [Issue Tracker](https://github.com/janmbaco/go-infrastructure/issues) - Report bugs, request features
- Contact maintainers through GitHub

---

**Made with Go** | **Star us on GitHub** | **[Become a Sponsor](https://github.com/sponsors/janmbaco)**

