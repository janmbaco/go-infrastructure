# Configuration

File-based configuration with hot-reload, rollback, and typed events — wired to the rest of `go-infrastructure`.

This module is made of:

- `configuration` – core interfaces (`ConfigHandler`, `Period`, subscriptions)
- `configuration/fileconfig` – JSON file–backed implementation with hot reload
- `configuration/events` – config change events integrated with `eventsmanager`
- `configuration/fileconfig/ioc` – DI module + resolver helpers

> Typical usage: you define a config struct, point the library at a JSON file, and your app sees updated values whenever the file changes — no restart required.   

---

## When to use this

Use the configuration module when you want:

- A single source of truth for app settings
- JSON config files that **auto-reload** on change
- The ability to **subscribe** to configuration changes
- A safe way to **rollback** to the previous config if something goes wrong
- Integration with `logs`, `errors`, `eventsmanager`, and `disk` via DI

---

## Quick Start: JSON config with hot reload

```go
package main

import (
    "fmt"
    di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
    cfgIoc "github.com/janmbaco/go-infrastructure/v2/configuration/fileconfig/ioc"
    cfgResolver "github.com/janmbaco/go-infrastructure/v2/configuration/fileconfig/ioc/resolver"
    logsIoc "github.com/janmbaco/go-infrastructure/v2/logs/ioc"
    errorsIoc "github.com/janmbaco/go-infrastructure/v2/errors/ioc"
    eventsIoc "github.com/janmbaco/go-infrastructure/v2/eventsmanager/ioc"
    diskIoc "github.com/janmbaco/go-infrastructure/v2/disk/ioc"
)

type AppConfig struct {
    Port     int    `json:"port"`
    Database string `json:"database"`
}

func main() {
    container := di.NewBuilder().
        AddModule(logsIoc.NewLogsModule()).
        AddModule(errorsIoc.NewErrorsModule()).
        AddModule(eventsIoc.NewEventsModule()).
        AddModule(diskIoc.NewDiskModule()).
        AddModule(cfgIoc.NewConfigurationModule()).
        MustBuild()

    resolver := container.Resolver()

    // Build a file-backed config handler with defaults
    configHandler := cfgResolver.GetFileConfigHandler(
        resolver,
        "config.json",
        &AppConfig{
            Port: 8080,
        },
    )

    cfg := configHandler.GetConfig().(*AppConfig)
    fmt.Printf("Server starting on port %d (DB: %s)\n", cfg.Port, cfg.Database)

    // The handler will keep updating cfg when config.json changes.
    select {} // keep the process alive
}
````

`config.json` is a standard JSON file:

```json
{
  "port": 8081,
  "database": "postgres://localhost/myapp"
}
```

When you edit and save `config.json`, the handler reloads the file, updates the in-memory config, and emits a *Modified* event underneath.

---

## Core concepts

### `ConfigHandler` interface

`ConfigHandler` lives in `configuration` and is the central abstraction:

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

**What it gives you:**

* `GetConfig()`
  Returns the current config as `interface{}`. You usually cast it to your struct:

  ```go
  cfg := configHandler.GetConfig().(*AppConfig)
  ```

* `SetConfig(value interface{}) error`
  Manually replace the current config in memory (e.g., after validating a dynamic change). Implementations can validate or reject invalid values.

* `Freeze()` / `Unfreeze()`
  Temporarily stop/resume applying new configs. This is useful if you’re in a critical section and don’t want the configuration to change underneath you (for example, while reconfiguring other services).

* `CanRestore() bool` + `Restore() error`
  Implementations (like `fileconfig`) keep the previous valid configuration around. If a new version causes trouble, you can roll back:

  ```go
  if configHandler.CanRestore() {
      if err := configHandler.Restore(); err != nil {
          // handle/ log rollback failure
      }
  }
  ```

  If there is no “old” config to go back to, implementations can return a typed error (see `HandlerErrorType` below).

* `SetRefreshTime(period Period) error`
  Associates a `Period` with the handler so implementations can decide when it’s OK to refresh again. This is primarily an advanced hook for throttling/controlling refresh behavior; in most apps you don’t need to call it directly.

* `ForceRefresh() error`
  Ask the handler to refresh its config immediately (e.g., re-read from disk) regardless of internal timing. For file-based configs this means “re-parse the file and update state now”, and it will typically emit a Modified or error.

---

### Subscribing to config events

`ConfigHandler` embeds two subscription interfaces:

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

Subscriptions are simple callbacks stored as function pointers:

```go
// Subscribe to "config modified" notifications
onModified := func() {
    cfg := configHandler.GetConfig().(*AppConfig)
    fmt.Println("Config updated, current port:", cfg.Port)
}
configHandler.ModifiedSubscribe(&onModified)
defer configHandler.ModifiedUnsubscribe(&onModified)

// Subscribe to "config restored" notifications
onRestored := func() {
    cfg := configHandler.GetConfig().(*AppConfig)
    fmt.Println("Config restored to previous version:", cfg.Port)
}
configHandler.RestoredSubscribe(&onRestored)
defer configHandler.RestoredUnsubscribe(&onRestored)
```

These callbacks run when the handler detects a change (e.g., after `config.json` is updated and successfully parsed) or when a restore happens.

---

### `Period`

```go
type Period struct {
    // contains filtered or unexported fields
}

func (p *Period) IsFinished() bool
```

`Period` is an opaque helper used by implementations to decide whether a “refresh interval” has finished. You see it in:

```go
SetRefreshTime(period Period) error
```

In normal use with `fileconfig`, you don’t need to construct or manipulate `Period` yourself; the file-based handler is driven primarily by file change notifications. If you implement your own `ConfigHandler`, `Period` gives you a reusable way to express “is it time to refresh yet?”.

---

## File-based configuration (`configuration/fileconfig`)

This package provides the standard JSON file implementation of `ConfigHandler`.

### `NewFileConfigHandler`

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

**Responsibilities:**

* Load the config from `filePath` as JSON into the given `defaults` type
* Start watching the file for changes (`disk.FileChangedNotifier`)
* When the file changes:

  * Re-load and parse JSON
  * If successful:

    * Swap the in-memory config
    * Emit a *Modified* event ({see below})
  * If something goes wrong:

    * Wrap the error as a `HandlerError`
    * Log via `logs.Logger`
    * Use `errors.ErrorCatcher` so the rest of your app doesn’t crash

You rarely call `NewFileConfigHandler` directly; instead, you let DI build and wire it via the configuration module (see “DI integration” below).

### Typed errors: `HandlerError` and `HandlerErrorType`

```go
type HandlerError interface {
    errors.CustomError
    GetErrorType() HandlerErrorType
}

type HandlerErrorType uint8

const (
    UnexpectedError HandlerErrorType = iota
    OldConfigNilError
)
```

* `UnexpectedError` – generic catch-all for unexpected failures (I/O, JSON parsing, etc.)
* `OldConfigNilError` – used when the handler cannot find “old” configuration (for example, when attempting a restore and there is nothing to roll back to).

You can detect and branch on these when handling errors from `ConfigHandler` methods:

```go
if err := configHandler.Restore(); err != nil {
    if hErr, ok := err.(fileconfig.HandlerError); ok {
        switch hErr.GetErrorType() {
        case fileconfig.UnexpectedError:
            // log or alert: something really wrong
        case fileconfig.OldConfigNilError:
            // nothing to restore to
        }
    }
}
```

---

## Configuration events (`configuration/events`)

The `events` subpackage defines typed events that plug into `eventsmanager`.

### Event types

```go
type ModifiedEvent struct{}

type RestoredEvent struct{}
```

Each event type implements the generic event API used by `eventsmanager`:

```go
func (e ModifiedEvent) GetEventArgs() ModifiedEvent
func (e ModifiedEvent) IsParallelPropagation() bool
func (e ModifiedEvent) StopPropagation() bool

func (e RestoredEvent) GetEventArgs() RestoredEvent
func (e RestoredEvent) IsParallelPropagation() bool
func (e RestoredEvent) StopPropagation() bool
```

* `ModifiedEvent` – fired when configuration changes (for example, after a file reload).
* `RestoredEvent` – fired when configuration is rolled back to a previous version.

### Event handlers

```go
type ModifiedEventHandler struct { /* ... */ }

func NewModifiedEventHandler(
    subscriptions eventsmanager.Subscriptions[ModifiedEvent],
) *ModifiedEventHandler

func (m *ModifiedEventHandler) ModifiedSubscribe(subscription *func())
func (m *ModifiedEventHandler) ModifiedUnsubscribe(subscription *func())

type RestoredEventHandler struct { /* ... */ }

func NewRestoredEventHandler(
    subscriptions eventsmanager.Subscriptions[RestoredEvent],
) *RestoredEventHandler

func (r *RestoredEventHandler) RestoredSubscribe(subscription *func())
func (r *RestoredEventHandler) RestoredUnsubscribe(subscription *func())
```

These are the building blocks behind the `ModifiedSubscriber` / `RestoredSubscriber` methods exposed on `ConfigHandler`.

You can also hook into config events directly via `eventsmanager` if you already use it elsewhere:

```go
import (
    "github.com/janmbaco/go-infrastructure/v2/eventsmanager"
    cfgEvents "github.com/janmbaco/go-infrastructure/v2/configuration/events"
    eventsIoc "github.com/janmbaco/go-infrastructure/v2/eventsmanager/ioc"
    eventsResolver "github.com/janmbaco/go-infrastructure/v2/eventsmanager/ioc/resolver"
)

// ... build container with eventsIoc.NewEventsModule() etc.

resolver := container.Resolver()
eventMgr := eventsResolver.GetEventManager(resolver)

// Subscribe using the generic eventsmanager API
subs := eventsmanager.NewSubscriptions[cfgEvents.ModifiedEvent]()
subs.Add(func(evt cfgEvents.ModifiedEvent) {
    cfg := configHandler.GetConfig().(*AppConfig)
    fmt.Println("Config changed via eventsmanager; new port:", cfg.Port)
})

publisher := eventsmanager.NewPublisher(subs, logger)
eventsmanager.Register(eventMgr, publisher)

// Now when FileConfigHandler publishes ModifiedEvent, your subscriber runs
```

---

## DI integration (`configuration/fileconfig/ioc`)

To avoid manually wiring `NewFileConfigHandler` and all its collaborators, you use the IoC module and resolver helper.

### Module: `NewConfigurationModule`

The module in `configuration/fileconfig/ioc` plugs into the `dependencyinjection` builder:

```go
import (
    di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
    cfgIoc "github.com/janmbaco/go-infrastructure/v2/configuration/fileconfig/ioc"
)

container := di.NewBuilder().
    AddModule(cfgIoc.NewConfigurationModule()).
    MustBuild()
```

This module registers all services needed by file-based configuration:

* The file configuration handler implementation
* The configuration events handlers
* Their dependencies from other modules (`logs`, `errors`, `eventsmanager`, `disk`), assuming you add those modules too (see the Quick Start example).

### Resolver helper: `GetFileConfigHandler`

```go
import (
    cfgResolver "github.com/janmbaco/go-infrastructure/v2/configuration/fileconfig/ioc/resolver"
)

handler := cfgResolver.GetFileConfigHandler(
    container.Resolver(),
    "config.json",
    &AppConfig{Port: 8080},
)
```

Signature:

```go
func GetFileConfigHandler(
    resolver dependencyinjection.Resolver,
    filePath string,
    defaults interface{},
) configuration.ConfigHandler
```

This helper:

* Resolves `ErrorCatcher`, `EventManager`, `FileChangedNotifier`, and `Logger` from DI
* Calls `NewFileConfigHandler` with those dependencies and your file path + defaults
* Returns a fully wired `ConfigHandler` ready to use

This is the recommended entrypoint for file configuration in a DI-based app.

---

## Putting it all together

A typical pattern for production services:

1. **Define your config struct**

   ```go
   type AppConfig struct {
       Port        string `json:"port"`
       StaticPath  string `json:"static_path"`
       Index       string `json:"index"`
       LogLevel    string `json:"log_level"`
   }
   ```

2. **Build the container with infrastructure modules**

   ```go
   container := di.NewBuilder().
       AddModule(logsIoc.NewLogsModule()).
       AddModule(errorsIoc.NewErrorsModule()).
       AddModule(eventsIoc.NewEventsModule()).
       AddModule(diskIoc.NewDiskModule()).
       AddModule(cfgIoc.NewConfigurationModule()).
       MustBuild()
   ```

3. **Resolve a `ConfigHandler`**

   ```go
   configHandler := cfgResolver.GetFileConfigHandler(
       container.Resolver(),
       "app.json",
       &AppConfig{
           Port:       ":8080",
           StaticPath: "./dist",
           Index:      "index.html",
           LogLevel:   "info",
       },
   )
   ```

4. **Subscribe to changes**

   ```go
   onModified := func() {
       cfg := configHandler.GetConfig().(*AppConfig)
       fmt.Println("Reloading server with new config:", cfg)
       // reconfigure other services here...
   }
   configHandler.ModifiedSubscribe(&onModified)
   ```

5. **Use the config everywhere (via DI)**
   You can inject `ConfigHandler` into services, or resolve it where needed, and always read the latest `*AppConfig` through `GetConfig()`.

---

## Summary

The configuration module gives you:

* A **generic `ConfigHandler` interface** for managing configuration and change notifications.
* A **file-based implementation** (`fileconfig`) that hot-reloads JSON files, with typed errors.
* **Typed events** (`ModifiedEvent`, `RestoredEvent`) that integrate with `eventsmanager`.
* A **DI module and resolver helper** so you can plug everything in with a couple of lines.

You get robust configuration plumbing with minimal code, while still having full control over how your app reacts to configuration changes.

