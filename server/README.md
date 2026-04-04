# Server

`github.com/janmbaco/go-infrastructure/v2/server`

The `server` package provides config-driven listeners for HTTP and gRPC servers, plus a lightweight single-page application handler.

## What It Includes

- `ListenerBuilder` to wire listeners from configuration
- `Listener` with `Start()` and `Stop()`
- HTTP and gRPC bootstrapping through `ServerSetter`
- Config validation and application hooks for live reload scenarios
- `NewSinglePageApp` and `facades.SinglePageAppStart`
- DI integration through `server/ioc`

## Install

```bash
go get github.com/janmbaco/go-infrastructure/v2/server
```

## Core API

```go
type ListenerBuilder interface {
    SetBootstrapper(BootstrapperFunc) ListenerBuilder
    SetGrpcDefinitions(GrpcDefinitionsFunc) ListenerBuilder
    SetConfigValidatorFunc(ConfigValidatorFunc) ListenerBuilder
    SetConfigApplicatorFunc(ConfigApplicatorFunc) ListenerBuilder
    GetListener() (Listener, error)
}

type Listener interface {
    Start() chan ListenerError
    Stop()
}
```

The bootstrapper receives the current config and fills a `ServerSetter`:

```go
type ServerSetter struct {
    Handler      http.Handler
    TLSConfig    *tls.Config
    TLSNextProto map[string]func(*http.Server, *tls.Conn, http.Handler)
    Name         string
    Addr         string
    ServerType   ServerType
}
```

Available server types:

```go
const (
    HTTPServer ServerType = iota
    GRpcSever
)
```

`GRpcSever` is the current exported gRPC constant name in the package.

## HTTP Quick Start

```go
package main

import (
    "net/http"
    "os"
    "os/signal"
    "syscall"

    configioc "github.com/janmbaco/go-infrastructure/v2/configuration/fileconfig/ioc/resolver"
    "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
    serverioc "github.com/janmbaco/go-infrastructure/v2/server/ioc"
    serverresolver "github.com/janmbaco/go-infrastructure/v2/server/ioc/resolver"
    "github.com/janmbaco/go-infrastructure/v2/server"
)

type Config struct {
    Address string `json:"address"`
}

func main() {
    container := dependencyinjection.NewBuilder().
        AddModules(serverioc.ConfigureServerModules()...).
        MustBuild()

    resolver := container.Resolver()
    configHandler := configioc.GetFileConfigHandler(resolver, "server.json", &Config{
        Address: ":8080",
    })

    listener, err := serverresolver.GetListenerBuilder(resolver, configHandler).
        SetBootstrapper(func(config interface{}, serverSetter *server.ServerSetter) error {
            cfg := config.(*Config)

            mux := http.NewServeMux()
            mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
                _, _ = w.Write([]byte("ok"))
            })

            serverSetter.Name = "http-api"
            serverSetter.Addr = cfg.Address
            serverSetter.Handler = mux
            return nil
        }).
        GetListener()
    if err != nil {
        panic(err)
    }

    finish := listener.Start()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    select {
    case err := <-finish:
        if err != nil {
            panic(err)
        }
    case <-quit:
        listener.Stop()
        if err := <-finish; err != nil {
            panic(err)
        }
    }
}
```

## gRPC Listener

For gRPC, set `ServerType` to `server.GRpcSever` and register protobuf services with `SetGrpcDefinitions`.

```go
listener, err := serverresolver.GetListenerBuilder(resolver, configHandler).
    SetBootstrapper(func(config interface{}, serverSetter *server.ServerSetter) error {
        cfg := config.(*Config)
        serverSetter.Name = "grpc-api"
        serverSetter.Addr = cfg.Address
        serverSetter.ServerType = server.GRpcSever
        return nil
    }).
    SetGrpcDefinitions(func(grpcServer *grpc.Server) {
        pb.RegisterGreeterServer(grpcServer, greeter)
    }).
    GetListener()
```

If `ServerType` is `server.GRpcSever` and no gRPC definitions are provided, `GetListener()` returns `NilGrpcDefinitionsError`.

## Live Configuration Hooks

`ListenerBuilder` exposes two optional hooks for config updates:

- `SetConfigValidatorFunc`: decide whether a new config should be applied.
- `SetConfigApplicatorFunc`: apply config changes without forcing a restart when possible.

If no applicator is provided, or if `ConfigApplication.NeedsRestart` is set to `true`, the listener restarts after config changes.

```go
listener, err := serverresolver.GetListenerBuilder(resolver, configHandler).
    SetBootstrapper(bootstrapper).
    SetConfigValidatorFunc(func(config interface{}) (bool, error) {
        cfg := config.(*Config)
        return cfg.Address != "", nil
    }).
    SetConfigApplicatorFunc(func(config interface{}, app *server.ConfigApplication) error {
        *app.NeedsRestart = true
        return nil
    }).
    GetListener()
```

## SPA Support

Serve static assets with SPA fallback:

```go
handler := server.NewSinglePageApp("./dist", "index.html")
```

For the opinionated executable-style bootstrap, use:

```go
facades.SinglePageAppStart(":8080", "./dist", "index.html")
```

The facade creates a config file next to the executable, wires the required modules and starts the listener.

## Error Types

Builder errors:

```go
const (
    UnexpectedBuilderError ListenerBuilderErrorType = iota
    NilBootstraperError
    NilGrpcDefinitionsError
)
```

Runtime listener errors:

```go
const (
    UnexpectedError ListenerErrorType = iota
    AddressNotConfigured
)
```

## Related Packages

- `server/ioc`: DI module and convenience module set
- `server/ioc/resolver`: helper to resolve `ListenerBuilder`
- `server/facades`: executable-style entry points
