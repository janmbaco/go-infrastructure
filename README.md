# Go Infrastructure

[![Go Report Card](https://goreportcard.com/badge/github.com/janmbaco/go-infrastructure)](https://goreportcard.com/report/github.com/janmbaco/go-infrastructure)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

**Build production-ready Go applications faster.** Common infrastructure solved: logging, configuration reload, dependency injection, events, and error handling—all working together out of the box.


---

## Why go-infrastructure?

**No more boilerplate** - DI container, logging, and config reload already wired  
**Live config reload** - Change JSON files, app adapts automatically without restart  
**Clean error handling** - No panic/recover antipatterns, errors are values  
**Type-safe events** - Generic-based pub/sub with compile-time type checking  

---

## Quick Start

### Example 1: Auto-Reloading Configuration

Your app reads `config.json` and **automatically reloads** when the file changes—no restart needed.

**config.json:**
```json
{
  "database": "postgres://localhost/myapp",
  "port": 8080
}
```

**main.go:**
```go
package main

import (
    "fmt"
    di "github.com/janmbaco/go-infrastructure/dependencyinjection"
    "github.com/janmbaco/go-infrastructure/configuration/fileconfig/ioc"
    configResolver "github.com/janmbaco/go-infrastructure/configuration/fileconfig/ioc/resolver"
    logsIoc "github.com/janmbaco/go-infrastructure/logs/ioc"
    errorsIoc "github.com/janmbaco/go-infrastructure/errors/ioc"
    eventsIoc "github.com/janmbaco/go-infrastructure/eventsmanager/ioc"
    diskIoc "github.com/janmbaco/go-infrastructure/disk/ioc"
)

type Config struct {
    Database string `json:"database"`
    Port     int    `json:"port"`
}

func main() {
    container := di.NewBuilder().
        AddModule(logsIoc.NewLogsModule()).
        AddModule(errorsIoc.NewErrorsModule()).
        AddModule(eventsIoc.NewEventsModule()).
        AddModule(diskIoc.NewDiskModule()).
        AddModule(ioc.NewConfigurationModule()).
        MustBuild()
    
    resolver := container.Resolver()
    configHandler := configResolver.GetFileConfigHandler(
        resolver,
        "config.json",
        &Config{Port: 8080}, // Default if file missing
    )
    
    cfg := configHandler.GetConfig().(*Config)
    fmt.Printf("App running on port %d\n", cfg.Port)
    
    // App automatically picks up changes to config.json
    select {} // Keep running
}
```

**What you get:**
- Changes to `config.json` detected automatically
- No restart needed
- Safe defaults if file is missing

---

### Example 2: Clean Error Handling

No panic/recover antipatterns. Catch errors centrally and log them.

```go
package main

import (
    "fmt"
    di "github.com/janmbaco/go-infrastructure/dependencyinjection"
    errorsIoc "github.com/janmbaco/go-infrastructure/errors/ioc"
    errorsResolver "github.com/janmbaco/go-infrastructure/errors/ioc/resolver"
    logsIoc "github.com/janmbaco/go-infrastructure/logs/ioc"
)

func processData(data []byte) error {
    if len(data) == 0 {
        return fmt.Errorf("empty data")
    }
    // Process...
    return nil
}

func main() {
    container := di.NewBuilder().
        AddModule(logsIoc.NewLogsModule()).
        AddModule(errorsIoc.NewErrorsModule()).
        MustBuild()
    
    errorCatcher := errorsResolver.GetErrorCatcher(container.Resolver())
    
    errorCatcher.TryCatchError(
        func() error { return processData([]byte{}) },
        func(err error) {
            // Error automatically logged
            fmt.Println("Operation failed:", err)
        },
    )
}
```

**What you get:**
- Centralized error logging
- No panic-based control flow
- Stack traces in logs automatically

---

### Example 3: Event-Driven Architecture

React to configuration changes or business events with type-safe subscribers.

```go
package main

import (
    "fmt"
    di "github.com/janmbaco/go-infrastructure/dependencyinjection"
    "github.com/janmbaco/go-infrastructure/eventsmanager"
    "github.com/janmbaco/go-infrastructure/configuration/events"
    eventsIoc "github.com/janmbaco/go-infrastructure/eventsmanager/ioc"
    eventsResolver "github.com/janmbaco/go-infrastructure/eventsmanager/ioc/resolver"
    logsIoc "github.com/janmbaco/go-infrastructure/logs/ioc"
    logsResolver "github.com/janmbaco/go-infrastructure/logs/ioc/resolver"
    errorsIoc "github.com/janmbaco/go-infrastructure/errors/ioc"
)

func main() {
    container := di.NewBuilder().
        AddModule(logsIoc.NewLogsModule()).
        AddModule(errorsIoc.NewErrorsModule()).
        AddModule(eventsIoc.NewEventsModule()).
        MustBuild()
    
    resolver := container.Resolver()
    eventMgr := eventsResolver.GetEventManager(resolver)
    logger := logsResolver.GetLogger(resolver)
    
    // Subscribe to config changes
    modifiedSubs := eventsmanager.NewSubscriptions[events.ModifiedEvent]()
    modifiedSubs.Add(func(evt events.ModifiedEvent) {
        fmt.Println("Config changed! Reloading services...")
    })
    
    modifiedPub := eventsmanager.NewPublisher(modifiedSubs, logger)
    eventsmanager.Register(eventMgr, modifiedPub)
    
    // Publish events
    eventsmanager.Publish(eventMgr, events.ModifiedEvent{})
}
```

**What you get:**
- React to config changes without polling
- Type-safe event handlers (compile-time checks)
- Automatic error handling in subscribers

---

## Complete Example: Production-Ready SPA Server

See how all modules work together in a real production application: a Single Page Application server with live config reload.

This example is included in `cmd/singlepageapp` and can be deployed as a Docker container.

### The Application

```go
package main

import (
    "flag"
    "github.com/janmbaco/go-infrastructure/logs"
    "github.com/janmbaco/go-infrastructure/server/facades"
)

func main() {
    port := flag.String("port", ":8080", "port to listen on")
    staticPath := flag.String("static", "./static", "path to static files")
    index := flag.String("index", "index.html", "index file name")
    flag.Parse()

    logger := logs.NewLogger()
    logger.Info("Starting SPA server on " + *port)
    
    facades.SinglePageAppStart(*port, *staticPath, *index)
}
```

### What This Gives You

**Automatic Configuration Reload** - Edit config file, server reconfigures without restart  
**Production Logging** - Daily rotated logs with configurable levels  
**Error Recovery** - Crashes are caught and logged, server stays running  
**Event-Driven** - React to config changes via type-safe events  
**Dockerized** - Ready-to-deploy container image

### How It Works

The `SinglePageAppStart` facade wires together:

1. **DI Container** - Manages all dependencies (logger, error catcher, config handler)
2. **File Config Handler** - Monitors `{app}.json` for changes
3. **Event Manager** - Publishes config change events
4. **Listener** - HTTP server that restarts on config changes
5. **Single Page App Handler** - Serves static files with SPA routing

### Configuration File

**app.json:**
```json
{
  "port": ":8080",
  "static_path": "./dist",
  "index": "index.html"
}
```

Edit this file while the server is running - it automatically reloads.

### Docker Deployment

**Using Official Image:**
```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/dist:/app/static \
  janmbaco/singlepageapp:latest \
  -port :8080 \
  -static /app/static \
  -index index.html
```

**Build Your Own:**
```bash
docker build -f server/facades/Dockerfile -t myapp:latest .
```

The official `janmbaco/singlepageapp:latest` image is built from this repository and published to Docker Hub on each release.

### Use Cases

**React/Vue/Angular Apps** - Zero-config SPA serving  
**Static Site Hosting** - Production-ready with logging  
**Microservices Frontend** - Containerized SPA delivery  
**Development Servers** - Live reload on config changes

See `cmd/singlepageapp/main.go` and `server/facades/singlepageapp_facade.go` for full implementation.

---

## Installation

```bash
go get github.com/janmbaco/go-infrastructure
```

---

## Common Patterns

### Pattern 1: Build a Modular DI Container

```go
package main

import (
    di "github.com/janmbaco/go-infrastructure/dependencyinjection"
    logsIoc "github.com/janmbaco/go-infrastructure/logs/ioc"
    errorsIoc "github.com/janmbaco/go-infrastructure/errors/ioc"
)

// Your custom module
type MyAppModule struct{}

func (m *MyAppModule) RegisterServices(register di.Register) error {
    register.AsSingleton(
        new(*UserService),
        func() *UserService { return &UserService{} },
        nil,
    )
    return nil
}

func main() {
    // Build container with multiple modules
    container := di.NewBuilder().
        AddModule(logsIoc.NewLogsModule()).
        AddModule(errorsIoc.NewErrorsModule()).
        AddModule(&MyAppModule{}).
        MustBuild()
    
    resolver := container.Resolver()
    // Use services...
}
```

---

### Pattern 2: React to Configuration Changes

```go
package main

import (
    "github.com/janmbaco/go-infrastructure/eventsmanager"
    "github.com/janmbaco/go-infrastructure/configuration/events"
)

func main() {
    // ... setup container and configHandler ...
    
    // Subscribe to config changes
    modifiedSubs := eventsmanager.NewSubscriptions[events.ModifiedEvent]()
    modifiedSubs.Add(func(evt events.ModifiedEvent) {
        newCfg := configHandler.GetConfig().(*AppConfig)
        updateServices(newCfg)
    })
    
    modifiedPub := eventsmanager.NewPublisher(modifiedSubs, logger)
    eventsmanager.Register(eventManager, modifiedPub)
}

func updateServices(cfg *AppConfig) {
    // Reconfigure services without restart
}
```

---

### Pattern 3: Type-Safe Generics

**Dependency Injection with Generics:**

```go
import di "github.com/janmbaco/go-infrastructure/dependencyinjection"

// Type-safe registration
di.RegisterSingleton[*UserService](register, func() *UserService {
    return &UserService{}
})

// Type-safe resolution
service := di.Resolve[*UserService](resolver)
// No type assertion needed
```

**Event Management with Generics:**

```go
import "github.com/janmbaco/go-infrastructure/eventsmanager"

type UserCreatedEvent struct {
    UserID string
}

// Type-safe subscriptions
subs := eventsmanager.NewSubscriptions[UserCreatedEvent]()
subs.Add(func(evt UserCreatedEvent) {
    // evt is typed, no cast needed
    fmt.Println("User created:", evt.UserID)
})

pub := eventsmanager.NewPublisher(subs, logger)
eventsmanager.Register(eventMgr, pub)

// Type-safe publish
eventsmanager.Publish(eventMgr, UserCreatedEvent{UserID: "123"})
```

### Dependency Injection

---

## Core Modules

### Dependency Injection

Modular DI container with type-safe registration and resolution.

#### Builder Pattern

```go
import (
    di "github.com/janmbaco/go-infrastructure/dependencyinjection"
    logsIoc "github.com/janmbaco/go-infrastructure/logs/ioc"
    errorsIoc "github.com/janmbaco/go-infrastructure/errors/ioc"
)

// Build container with pre-defined modules
container := di.NewBuilder().
    AddModule(logsIoc.NewLogsModule()).
    AddModule(errorsIoc.NewErrorsModule()).
    AddModule(customModule).
    MustBuild()

// Or with custom registration
container := di.NewBuilder().
    Register(func(r di.Register) {
        r.AsSingleton(new(*MyService), NewMyService, nil)
    }).
    MustBuild()
```

#### Module Definition

```go
type MyModule struct{}

func NewMyModule() di.Module {
    return &MyModule{}
}

func (m *MyModule) RegisterServices(register di.Register) error {
    register.AsSingleton(
        new(*MyService),
        func() *MyService { return &MyService{} },
        nil,
    )
    return nil
}
```

#### Generic Registration

```go
// Type-safe registration
di.RegisterSingleton[*MyService](register, func() *MyService {
    return &MyService{}
})

// Type-safe resolution
service := di.Resolve[*MyService](resolver)
```

#### Lifetime Scopes

```go
// Singleton: One instance per container
register.AsSingleton(new(*Logger), NewLogger, nil)

// Transient: New instance per request
register.AsTransient(new(*Service), NewService, nil)

// Scoped: One instance per scope
register.AsScoped(new(*RequestService), NewRequestService, nil)
```

---

### Logging

Structured logging with file rotation and configurable levels.

#### Configuration

```go
import "github.com/janmbaco/go-infrastructure/logs"

logger := logs.NewLogger()

// Set log levels
logger.SetConsoleLevel(logs.Info)    // Console: Info and above
logger.SetFileLogLevel(logs.Trace)   // File: Everything

// Set log directory
logger.SetDir("./logs")
```

#### Log Levels

| Level | Value | Method | Description |
|-------|-------|--------|-------------|
| Off | 0 | - | No logging |
| Fatal | 1 | `Fatal(msg)` | Critical errors (exits program) |
| Error | 2 | `Error(msg)` | Errors |
| Warning | 3 | `Warning(msg)` | Warnings |
| Info | 4 | `Info(msg)` | Informational messages |
| Trace | 5 | `Trace(msg)` | Detailed debugging |

#### Usage

```go
logger.Trace("Entering function")
logger.Info("User logged in: " + username)
logger.Warning("Cache miss for key: " + key)
logger.Error("Failed to connect: " + err.Error())
logger.Fatal("Database unavailable") // Exits program
```

Logs are written to `{logs_dir}/YYYY-MM-DD.log` with daily rotation.

---

### Configuration

File-based configuration with automatic reload on changes.

#### Setup

```go
import (
    "github.com/janmbaco/go-infrastructure/configuration/fileconfig/ioc"
    configResolver "github.com/janmbaco/go-infrastructure/configuration/fileconfig/ioc/resolver"
)

type AppConfig struct {
    Port     int    `json:"port"`
    Database string `json:"database"`
}

// Build container with configuration module
container := di.NewBuilder().
    AddModule(logsIoc.NewLogsModule()).
    AddModule(errorsIoc.NewErrorsModule()).
    AddModule(eventsIoc.NewEventsModule()).
    AddModule(diskIoc.NewDiskModule()).
    AddModule(ioc.NewConfigurationModule()).
    MustBuild()

// Get configuration handler
configHandler := configResolver.GetFileConfigHandler(
    resolver,
    "config.json",
    &AppConfig{Port: 8080}, // Defaults
)

// Get current configuration
config := configHandler.GetConfig().(*AppConfig)
fmt.Println("Port:", config.Port)
```

#### Configuration Events

```go
import (
    "github.com/janmbaco/go-infrastructure/eventsmanager"
    "github.com/janmbaco/go-infrastructure/configuration/events"
)

// Subscribe to config changes
modifiedSubs := eventsmanager.NewSubscriptions[events.ModifiedEvent]()
modifiedSubs.Add(func(evt events.ModifiedEvent) {
    newConfig := configHandler.GetConfig().(*AppConfig)
    fmt.Println("Config updated. New port:", newConfig.Port)
})

modifiedPub := eventsmanager.NewPublisher(modifiedSubs, logger)
eventsmanager.Register(eventMgr, modifiedPub)

// Subscribe to config restoration
restoredSubs := eventsmanager.NewSubscriptions[events.RestoredEvent]()
restoredSubs.Add(func(evt events.RestoredEvent) {
    fmt.Println("Config restored to previous version")
})

restoredPub := eventsmanager.NewPublisher(restoredSubs, logger)
eventsmanager.Register(eventMgr, restoredPub)
```

---

### Error Handling

Centralized error catching without panic-based antipatterns.

#### Error Catcher

```go
import (
    "github.com/janmbaco/go-infrastructure/errors"
    errorsResolver "github.com/janmbaco/go-infrastructure/errors/ioc/resolver"
)

errorCatcher := errorsResolver.GetErrorCatcher(resolver)

// Try-Catch pattern
errorCatcher.TryCatchError(
    func() error {
        // Risky operation
        return someOperation()
    },
    func(err error) {
        // Handle error (automatically logged)
        fmt.Println("Caught error:", err)
    },
)
```

#### Validation

```go
import "github.com/janmbaco/go-infrastructure/errors"

func ProcessUser(user *User, db *Database) error {
    // Validate parameters (returns error, no panic)
    if err := errors.ValidateNotNil(map[string]interface{}{
        "user": user,
        "db":   db,
    }); err != nil {
        return err
    }
    
    // Process...
    return nil
}
```

---

### Event Management

Type-safe publisher/subscriber pattern with generics.

#### Event Definition

```go
import "github.com/janmbaco/go-infrastructure/eventsmanager"

type UserCreatedEvent struct {
    UserID string
    Email  string
}

// Events are passed directly, no wrapper needed
```

#### Publish & Subscribe

```go
import (
    "github.com/janmbaco/go-infrastructure/eventsmanager"
    eventsResolver "github.com/janmbaco/go-infrastructure/eventsmanager/ioc/resolver"
)

eventMgr := eventsResolver.GetEventManager(resolver)

// Create subscriptions with type safety
subscriptions := eventsmanager.NewSubscriptions[UserCreatedEvent]()
subscriptions.Add(func(event UserCreatedEvent) {
    fmt.Printf("User created: %s (%s)\n", event.UserID, event.Email)
})

// Create publisher
publisher := eventsmanager.NewPublisher(subscriptions, logger)
eventsmanager.Register(eventMgr, publisher)

// Publish event (type-safe)
eventsmanager.Publish(eventMgr, UserCreatedEvent{
    UserID: "123",
    Email:  "user@example.com",
})
```

---

### Server

HTTP/HTTPS server with graceful shutdown and configuration reload.

#### HTTP Server

```go
import (
    "net/http"
    "github.com/janmbaco/go-infrastructure/server"
    serverResolver "github.com/janmbaco/go-infrastructure/server/ioc/resolver"
)

// Create listener builder
listenerBuilder := serverResolver.GetListenerBuilder(resolver, configHandler)
listenerBuilder.SetBootstrapper(func(cfg interface{}, serverSetter *server.ServerSetter) error {
    serverSetter.Name = "MyServer"
    serverSetter.Addr = ":8080"
    serverSetter.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    return nil
})

listener, err := listenerBuilder.GetListener()
if err != nil {
    log.Fatal(err)
}

// Start server (blocks until stopped)
finish := listener.Start()
if err := <-finish; err != nil {
    log.Fatal(err)
}
```

#### Single Page App Serving

```go
import "github.com/janmbaco/go-infrastructure/server"

spaHandler := server.NewSinglePageApp("./dist", "index.html")
http.Handle("/", spaHandler)
```

---

### Persistence

GORM-based data access layer with multi-database support.

#### Database Setup

```go
import (
    "github.com/janmbaco/go-infrastructure/persistence/orm_base"
    "gorm.io/gorm"
)

type User struct {
    ID   uint   `gorm:"primaryKey"`
    Name string
}

dbInfo := &orm_base.DatabaseInfo{
    Engine:       orm_base.Postgres,
    Host:         "localhost",
    Port:         "5432",
    Name:         "myapp",
    UserName:     "postgres",
    UserPassword: "password",
}

db := orm_base.NewDB(
    dialectorResolver,
    dbInfo,
    &gorm.Config{},
    []interface{}{&User{}}, // Auto-migrate
)
```

#### Supported Databases

| Engine | Connection String Format |
|--------|--------------------------|
| PostgreSQL | `host=X port=Y user=Z password=W dbname=D` |
| MySQL | `user:password@tcp(host:port)/dbname` |
| SQLite | `./test.db` |
| SQL Server | `sqlserver://user:password@host:port?database=D` |

---

### Crypto

AES encryption/decryption service.

```go
import (
    "github.com/janmbaco/go-infrastructure/crypto"
    cryptoResolver "github.com/janmbaco/go-infrastructure/crypto/ioc/resolver"
)

key := []byte("my-32-byte-encryption-key!!!") // Must be 16, 24, or 32 bytes
cipher := cryptoResolver.GetCipher(resolver, key)

// Encrypt
plaintext := []byte("Secret message")
encrypted := cipher.Encrypt(plaintext)

// Decrypt
decrypted := cipher.Decrypt(encrypted)
fmt.Println(string(decrypted)) // "Secret message"
```

---

### Disk Utilities

File change notification and utilities.

```go
import (
    "github.com/janmbaco/go-infrastructure/disk"
    diskResolver "github.com/janmbaco/go-infrastructure/disk/ioc/resolver"
)

fileNotifier := diskResolver.GetFileChangedNotifier(resolver, "config.json")
fileNotifier.Subscribe(func() {
    fmt.Println("File changed!")
})
```

---

## Architecture

### Module-Based Design

```
┌─────────────────────────────────────────────────┐
│            Application Layer                    │
│  (Your application using go-infrastructure)     │
└───────────────────┬─────────────────────────────┘
                    │
┌───────────────────▼─────────────────────────────┐
│         Dependency Injection Container          │
│  (Builder → Modules → Register → Resolver)      │
└───────────────────┬─────────────────────────────┘
                    │
        ┌───────────┼───────────┬─────────┬─────────────┐
        │           │           │         │             │
┌───────▼────┐ ┌───▼────┐ ┌───▼────┐ ┌─▼──────┐ ┌────▼─────┐
│   Logs     │ │ Errors │ │ Events │ │ Config │ │  Server  │
│   Module   │ │ Module │ │ Module │ │ Module │ │  Module  │
└────────────┘ └────────┘ └────────┘ └────────┘ └──────────┘
```

### Key Principles

**Separation of Concerns** - Each module handles a specific infrastructure concern  
**Dependency Injection** - Services resolved via DI container  
**Event-Driven** - Configuration changes and errors trigger events  
**Type Safety** - Generic functions for compile-time type checking  
**Error Handling** - Error-based (not panic-based) for better control

---

## Testing

### Unit Testing with DI

```go
package mypackage_test

import (
    "testing"
    di "github.com/janmbaco/go-infrastructure/dependencyinjection"
)

func TestMyService(t *testing.T) {
    container := di.NewBuilder().
        Register(func(r di.Register) {
            // Register test dependencies
            r.AsSingleton(new(*MockLogger), func() *MockLogger {
                return &MockLogger{}
            }, nil)
        }).
        MustBuild()
    
    service := container.Resolver().Type(new(*MyService), nil).(*MyService)
    
    // Test service...
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run specific package
go test ./dependencyinjection
```

---

## Requirements

- **Go**: 1.24+
- **OS**: Linux, macOS, Windows
- **Dependencies**: See `go.mod` for full list

### Key Dependencies

- `github.com/fsnotify/fsnotify` - File watching
- `gorm.io/gorm` - ORM framework
- `gorm.io/driver/*` - Database drivers

---

## Troubleshooting

### Issue: Module Import Error

**Error:**
```
cannot find package "github.com/janmbaco/go-infrastructure/..."
```

**Solution:**
```bash
go get github.com/janmbaco/go-infrastructure
go mod tidy
```

---

### Issue: Event Subscribers Not Called

**Symptom:** Published events don't trigger subscribers

**Solution:** Ensure you're using the same `EventManager` instance:
```go
// Wrong - different publisher instances
subs1 := eventsmanager.NewSubscriptions[MyEvent]()
pub1 := eventsmanager.NewPublisher(subs1, logger)

subs2 := eventsmanager.NewSubscriptions[MyEvent]()
pub2 := eventsmanager.NewPublisher(subs2, logger) // Different instance

// Correct - resolve from DI
eventMgr := eventsResolver.GetEventManager(resolver)
eventsmanager.Register(eventMgr, publisher)
eventsmanager.Publish(eventMgr, event)
```

---

### Issue: Config Not Reloading

**Symptom:** Changes to `config.json` don't take effect

**Checklist:**
1. File watcher is running (check logs for "watching file")
2. Event subscribers are registered before file changes
3. Config file is being edited, not replaced (use edit/save, not mv/cp)

---

### Issue: ValidateNotNil Returns Error

**Error:**
```
parameters cannot be nil: user
```

**Solution:** Check that all parameters are non-nil before calling the function:
```go
if user == nil {
    return fmt.Errorf("user cannot be nil")
}
```

---

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Follow Go conventions: Use `gofmt`, `golangci-lint`
4. Write tests: Maintain coverage
5. Commit with clear messages: `git commit -m 'Add feature: XYZ'`
6. Open a Pull Request: Include description and tests

### Development Setup

```bash
# Clone repository
git clone https://github.com/janmbaco/go-infrastructure.git
cd go-infrastructure

# Install dependencies
go mod download

# Run tests
go test ./...

# Run linter (if installed)
golangci-lint run

# Check coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## Additional Resources

- [UPGRADE-2.0.md](./UPGRADE-2.0.md) - Migration guide from v1.x
- [CHANGELOG.md](./CHANGELOG.md) - Version history and release notes
- [GitHub Issues](https://github.com/janmbaco/go-infrastructure/issues) - Report bugs or request features

---

## License

This project is licensed under the **Apache License 2.0** - see the [LICENSE](LICENSE) file for details.

---


