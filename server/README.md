# Server

`github.com/janmbaco/go-infrastructure/v2/server`

Production-ready HTTP/HTTPS server with graceful shutdown, configuration integration, and single-page application support.

The server module provides:

- Fluent `ListenerBuilder` for server configuration
- HTTP/HTTPS support with TLS
- Single-page application (SPA) handler with client-side routing
- Graceful shutdown with connection draining
- Dynamic configuration with live reload
- DI module to inject configured servers across your app

---

## Install

```bash
go get github.com/janmbaco/go-infrastructure/v2/server
```

---

## Quick Start

### Simple HTTP server

```go
package main

import (
    "net/http"
    "github.com/janmbaco/go-infrastructure/v2/server"
)

func main() {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

    listener := server.NewListenerBuilder().
        SetPort(":8080").
        SetHandler(handler).
        Build()

    if err := listener.ListenAndServe(); err != nil {
        panic(err)
    }
}
```

---

## Core API

### ListenerBuilder

The `ListenerBuilder` interface provides a fluent API for configuring servers:

```go
type ListenerBuilder interface {
    SetPort(port string) ListenerBuilder
    SetHandler(handler http.Handler) ListenerBuilder
    SetTLSConfig(config *tls.Config) ListenerBuilder
    SetServerSetter(setter ServerSetter) ListenerBuilder
    Build() Listener
}
```

Create a builder:

```go
builder := server.NewListenerBuilder()
```

Configure and build:

```go
listener := builder.
    SetPort(":8080").
    SetHandler(handler).
    Build()
```

### Listener

The `Listener` interface manages server lifecycle:

```go
type Listener interface {
    ListenAndServe() error
    Shutdown(ctx context.Context) error
    Port() string
    IsHTTPS() bool
}
```

* **ListenAndServe** – starts the server (blocks until shutdown or error)
* **Shutdown** – gracefully shuts down with context timeout
* **Port** – returns the configured port
* **IsHTTPS** – returns `true` if TLS is configured

---

## HTTPS support

Add TLS configuration to enable HTTPS:

```go
listener := server.NewListenerBuilder().
    SetPort(":8443").
    SetHandler(handler).
    SetTLSConfig(&tls.Config{
        Certificates: []tls.Certificate{cert},
        MinVersion:   tls.VersionTLS12,
    }).
    Build()

if err := listener.ListenAndServe(); err != nil {
    panic(err)
}
```

Load certificates:

```go
cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
if err != nil {
    panic(err)
}
```

---

## Single-page application (SPA) support

The server module includes built-in support for SPAs (React, Vue, Angular, etc.).

### Quick start with `facades.SinglePageAppStart`

```go
package main

import (
    "flag"
    "github.com/janmbaco/go-infrastructure/v2/server/facades"
)

func main() {
    port := flag.String("port", ":8080", "server port")
    static := flag.String("static", "./dist", "static files directory")
    index := flag.String("index", "index.html", "index file")
    flag.Parse()

    facades.SinglePageAppStart(*port, *static, *index)
}
```

This helper function:

* Serves static assets (JS, CSS, images)
* Returns `index.html` for all unmatched routes (enables client-side routing)
* Handles 404s gracefully
* Prevents path traversal attacks

### Using the SPA handler directly

For more control:

```go
spaHandler := server.NewSinglePageApp("./dist", "index.html")

listener := server.NewListenerBuilder().
    SetPort(":8080").
    SetHandler(spaHandler).
    Build()

listener.ListenAndServe()
```

The `NewSinglePageApp` function returns an `http.Handler` that:

```go
// Static file serving with fallback to index.html
// Handles:
// - /assets/app.js → serves from staticPath
// - /app/route → returns index.html (client-side routing)
// - / → returns index.html
```

---

## Graceful shutdown

Handle OS signals for clean shutdown:

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/janmbaco/go-infrastructure/v2/server"
)

func main() {
    listener := server.NewListenerBuilder().
        SetPort(":8080").
        SetHandler(myHandler()).
        Build()

    // Start server in goroutine
    go func() {
        if err := listener.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := listener.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exited")
}

func myHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello!"))
    })
}
```

---

## Configuration integration

Integrate with the `configuration` module for dynamic updates:

```go
package main

import (
    "net/http"

    di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
    "github.com/janmbaco/go-infrastructure/v2/server/ioc"
    "github.com/janmbaco/go-infrastructure/v2/server/ioc/resolver"
    configioc "github.com/janmbaco/go-infrastructure/v2/configuration/fileconfig/ioc"
)

type Config struct {
    Port string `json:"port"`
    TLS  struct {
        Enabled  bool   `json:"enabled"`
        CertFile string `json:"certFile"`
        KeyFile  string `json:"keyFile"`
    } `json:"tls"`
}

func main() {
    container := di.NewBuilder().
        AddModule(configioc.NewConfigurationModule()).
        AddModule(ioc.NewServerModule()).
        Register(func(r di.Register) {
            r.AsSingleton(new(http.Handler), func() http.Handler {
                return myAppHandler()
            }, nil)
        }).
        MustBuild()

    builder := resolver.GetListenerBuilder(container.Resolver())
    listener := builder.Build()

    if err := listener.ListenAndServe(); err != nil {
        panic(err)
    }
}

func myAppHandler() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`{"status":"ok"}`))
    })
    return mux
}
```

Configuration file format (`config.json`):

```json
{
  "server": {
    "port": ":8080",
    "tls": {
      "enabled": true,
      "certFile": "/path/to/cert.pem",
      "keyFile": "/path/to/key.pem"
    }
  }
}
```

---

## DI integration (`server/ioc`)

The `server/ioc` package provides a DI module for the container.

### ServerModule

```go
type ServerModule struct{}

func NewServerModule() *ServerModule
func (m *ServerModule) RegisterServices(register dependencyinjection.Register) error
```

`ServerModule` implements `dependencyinjection.Module` and registers server components.

### Using with `dependencyinjection`

```go
import (
    di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
    serverIoc "github.com/janmbaco/go-infrastructure/v2/server/ioc"
    serverResolver "github.com/janmbaco/go-infrastructure/v2/server/ioc/resolver"
)

func main() {
    container := di.NewBuilder().
        AddModule(serverIoc.NewServerModule()).
        MustBuild()

    builder := serverResolver.GetListenerBuilder(container.Resolver())
    listener := builder.Build()

    listener.ListenAndServe()
}
```

---

## Error handling

The server module defines specific error types:

### ListenerBuilderError

Configuration errors during builder setup:

```go
type ListenerBuilderError interface {
    error
    GetErrorType() ListenerBuilderErrorType
}

type ListenerBuilderErrorType uint8

const (
    BuilderUnexpected ListenerBuilderErrorType = iota
    PortNotSet
    HandlerNotSet
)
```

* `BuilderUnexpected` – internal builder error
* `PortNotSet` – `SetPort` was not called
* `HandlerNotSet` – `SetHandler` was not called

### ListenerError

Runtime errors during server operation:

```go
type ListenerError interface {
    error
    GetErrorType() ListenerErrorType
}

type ListenerErrorType uint8

const (
    ListenerUnexpected ListenerErrorType = iota
    BindError
    ListenError
    ServeError
)
```

* `ListenerUnexpected` – internal listener error
* `BindError` – failed to bind to port
* `ListenError` – failed to start listening
* `ServeError` – error during request handling

Example:

```go
if err := listener.ListenAndServe(); err != nil {
    if lErr, ok := err.(server.ListenerError); ok {
        switch lErr.GetErrorType() {
        case server.BindError:
            logger.Error("port already in use:", lErr)
        case server.ServeError:
            logger.Error("serving error:", lErr)
        default:
            logger.Error("listener error:", lErr)
        }
    }
}
```

---

## Examples

### Complete SPA application

See the working example: [`cmd/singlepageapp`](../cmd/singlepageapp)

Features demonstrated:

* Static file serving
* SPA routing support
* Configuration management
* Docker deployment
* Graceful shutdown

Run it:

```bash
go run ./cmd/singlepageapp -port :8080 -static ./dist -index index.html
```

### Docker deployment

Build and run in container:

```bash
# Build image
docker build -f server/facades/Dockerfile -t myapp .

# Run container
docker run -p 8080:8080 -v ./config.json:/app/config.json myapp
```

---

## Security best practices

### Always use HTTPS in production

```go
SetTLSConfig(&tls.Config{
    MinVersion: tls.VersionTLS12,
    CipherSuites: []uint16{
        tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
    },
})
```

### Implement request timeouts

```go
srv := &http.Server{
    ReadTimeout:  15 * time.Second,
    WriteTimeout: 15 * time.Second,
    IdleTimeout:  60 * time.Second,
}
```

### Use secure headers

```go
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("X-Frame-Options", "DENY")
w.Header().Set("Content-Security-Policy", "default-src 'self'")
```

---

## Health checks

### Basic health check

```go
mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
})
```

### Readiness check

```go
mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
    // Check dependencies (DB, cache, etc.)
    if allReady {
        w.WriteHeader(http.StatusOK)
    } else {
        w.WriteHeader(http.StatusServiceUnavailable)
    }
})
```

---

## Testing

Example test for SPA handler:

```go
func TestSinglePageApp_ServeHTTP(t *testing.T) {
    handler := server.NewSinglePageApp("./testdata", "index.html")
    
    tests := []struct {
        name           string
        path           string
        expectedStatus int
    }{
        {"root path", "/", http.StatusOK},
        {"static file", "/app.js", http.StatusOK},
        {"SPA route", "/app/dashboard", http.StatusOK},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("GET", tt.path, nil)
            rec := httptest.NewRecorder()
            
            handler.ServeHTTP(rec, req)
            
            assert.Equal(t, tt.expectedStatus, rec.Code)
        })
    }
}
```

---

## Troubleshooting

### Port already in use

```
Error: bind: address already in use
```

Check for processes using the port:

```bash
# Linux/Mac
lsof -i :8080

# Windows
netstat -ano | findstr :8080
```

### TLS certificate errors

```
Error: tls: failed to find any PEM data
```

Verify certificate files exist and are valid PEM format:

```bash
openssl x509 -in cert.pem -text -noout
```

### SPA routes return 404

Ensure you're using `NewSinglePageApp` handler, not a basic file server.

---

## Summary

The `server` module gives you:

* A **fluent builder** for configuring HTTP/HTTPS servers
* Built-in **SPA support** with client-side routing
* **Graceful shutdown** with context-based timeout
* **Configuration integration** via the `configuration` module
* **DI module** to wire servers into your application
* **Type-safe error handling** with specific error types

Use it to build production-ready web services with minimal boilerplate and maximum flexibility.
