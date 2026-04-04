# Configuration

`github.com/janmbaco/go-infrastructure/v2/configuration`

The `configuration` package provides a file-backed configuration handler with change notifications, rollback support and integration with the rest of `go-infrastructure`.

## What It Includes

- `ConfigHandler` as the main abstraction
- `configuration/fileconfig` as the JSON file implementation
- `configuration/events` for modified and restored event types
- `configuration/fileconfig/ioc` and `configuration/fileconfig/ioc/resolver` for DI wiring
- `Period` for advanced refresh scheduling

## Install

```bash
go get github.com/janmbaco/go-infrastructure/v2/configuration
go get github.com/janmbaco/go-infrastructure/v2/configuration/fileconfig
```

## Quick Start

The usual setup is to wire the required modules through DI and resolve a file-backed `ConfigHandler`.

```go
package main

import (
    "fmt"

    cfgresolver "github.com/janmbaco/go-infrastructure/v2/configuration/fileconfig/ioc/resolver"
    cfgioc "github.com/janmbaco/go-infrastructure/v2/configuration/fileconfig/ioc"
    di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
    diskioc "github.com/janmbaco/go-infrastructure/v2/disk/ioc"
    errorsioc "github.com/janmbaco/go-infrastructure/v2/errors/ioc"
    eventsioc "github.com/janmbaco/go-infrastructure/v2/eventsmanager/ioc"
    logsioc "github.com/janmbaco/go-infrastructure/v2/logs/ioc"
)

type AppConfig struct {
    Port string `json:"port"`
    DSN  string `json:"dsn"`
}

func main() {
    container := di.NewBuilder().
        AddModule(logsioc.NewLogsModule()).
        AddModule(errorsioc.NewErrorsModule()).
        AddModule(eventsioc.NewEventsModule()).
        AddModule(diskioc.NewDiskModule()).
        AddModule(cfgioc.NewConfigurationModule()).
        MustBuild()

    configHandler := cfgresolver.GetFileConfigHandler(
        container.Resolver(),
        "config.json",
        &AppConfig{Port: ":8080"},
    )

    cfg := configHandler.GetConfig().(*AppConfig)
    fmt.Println(cfg.Port, cfg.DSN)
}
```

If `config.json` does not exist, the file-based handler writes the provided defaults first.

## Core Interface

```go
type ConfigHandler interface {
    ModifiedSubscriber
    RestoredSubscriber

    GetConfig() interface{}
    SetConfig(interface{}) error

    Freeze()
    Unfreeze()

    CanRestore() bool
    Restore() error

    SetRefreshTime(Period) error
    ForceRefresh() error
}
```

## Main Behaviors

- `GetConfig` returns the current in-memory config object.
- `SetConfig` updates the current config and persists it back to the file.
- `Freeze` and `Unfreeze` control whether externally detected file changes are applied immediately.
- `CanRestore` and `Restore` let you roll back to the previous config snapshot.
- `SetRefreshTime` and `ForceRefresh` are advanced hooks for refresh scheduling and pending updates.

## File-Based Handler

`configuration/fileconfig.NewFileConfigHandler` is the standard implementation:

```go
func NewFileConfigHandler(
    filePath string,
    defaults interface{},
    errorCatcher errors.ErrorCatcher,
    eventManager *eventsmanager.EventManager,
    filechangeNotifier disk.FileChangedNotifier,
    logger logs.Logger,
) (configuration.ConfigHandler, error)
```

The file-based handler:

- keeps an in-memory config of the same type as `defaults`
- creates the file if it does not exist
- watches the file for write events
- updates the in-memory config when the file changes
- publishes `ModifiedEvent` and `RestoredEvent` through `eventsmanager`

## Subscriptions

`ConfigHandler` supports two callback-based subscriptions:

```go
type ModifiedSubscriber interface {
    ModifiedSubscribe(*func())
    ModifiedUnsubscribe(*func())
}

type RestoredSubscriber interface {
    RestoredSubscribe(*func())
    RestoredUnsubscribe(*func())
}
```

Example:

```go
onModified := func() {
    cfg := configHandler.GetConfig().(*AppConfig)
    fmt.Println("updated port:", cfg.Port)
}

configHandler.ModifiedSubscribe(&onModified)
defer configHandler.ModifiedUnsubscribe(&onModified)
```

## Rollback Errors

`configuration/fileconfig` exposes a typed error for handler failures:

```go
type HandlerError interface {
    errors.CustomError
    GetErrorType() HandlerErrorType
}

const (
    UnexpectedError HandlerErrorType = iota
    OldConfigNilError
)
```

`OldConfigNilError` is returned when `Restore()` is requested but there is no previous config snapshot available.

## About `Period`

`Period` is an advanced helper used by the refresh loop:

```go
type Period struct {
    // unexported fields
}

func (p *Period) IsFinished() bool
```

Most applications using `fileconfig` do not need to construct `Period` directly.

## Related Packages

- `configuration/fileconfig`: JSON-backed implementation
- `configuration/events`: config change event types
- `disk`: file change notifications
- `eventsmanager`: event publication
