# Events Manager (`eventsmanager`)

`github.com/janmbaco/go-infrastructure/v2/eventsmanager`

A small, generic **event bus** built on Go generics:

- Type-safe events (no `interface{}` casting)
- Pluggable publishers with logging
- Simple subscription model
- Integration with `logs`, `disk`, `configuration`, etc.
- IoC module to inject a shared `EventManager` across your app   

Typical use: you define an event type, add subscribers, and publish events through a shared `EventManager`.

---

## When to use this

Use `eventsmanager` when you want:

- Decoupled communication between parts of your app (publish/subscribe)
- Typed events (`UserCreatedEvent`, `FileChangedEvent`, `ConfigModifiedEvent`, …)
- Optional parallel dispatch and propagation control
- Integration with existing infrastructure modules (e.g., configuration + disk watcher already use it).   

---

## Quick Start

### 1. Define an event type

Implement the `EventObject[T]` interface:

```go
import "github.com/janmbaco/go-infrastructure/v2/eventsmanager"

type UserCreatedEvent struct {
    UserID string
    Email  string

    stopPropagation bool
}

func (e UserCreatedEvent) GetEventArgs() UserCreatedEvent {
    // Return the payload you want subscribers to receive
    return e
}

func (e UserCreatedEvent) StopPropagation() bool {
    // If true, publishing stops after this event; keep false for most cases
    return e.stopPropagation
}

func (e UserCreatedEvent) IsParallelPropagation() bool {
    // If true, subscribers may be invoked in parallel
    return false
}
````

`EventObject[T]`:

```go
type EventObject[T any] interface {
    GetEventArgs() T
    StopPropagation() bool
    IsParallelPropagation() bool
}
```

---

### 2. Create subscriptions and a publisher

```go
import (
    "fmt"

    "github.com/janmbaco/go-infrastructure/v2/eventsmanager"
    "github.com/janmbaco/go-infrastructure/v2/logs"
)

func main() {
    logger := logs.NewLogger()

    // Create subscriptions for this event type
    subs := eventsmanager.NewSubscriptions[UserCreatedEvent]()

    // Add subscribers (handlers)
    subs.Add(func(evt UserCreatedEvent) {
        fmt.Println("Sending welcome email to:", evt.Email)
    })

    subs.Add(func(evt UserCreatedEvent) {
        fmt.Println("Auditing new user:", evt.UserID)
    })

    // Create a publisher using these subscriptions
    publisher := eventsmanager.NewPublisher(subs, logger)

    // Create an EventManager and register the publisher for UserCreatedEvent
    em := eventsmanager.NewEventManager()
    eventsmanager.Register(em, publisher)

    // Publish an event
    eventsmanager.Publish(em, UserCreatedEvent{
        UserID: "123",
        Email:  "user@example.com",
    })
}
```

What’s happening:

* `Subscriptions[T]` stores all subscriber functions for event type `T`.
* `Publisher[T]` wraps those subscriptions and logs any errors.
* `EventManager` maps event types to their publishers.
* `Register` associates a `Publisher[T]` with an `EventManager` for type `T`.
* `Publish` sends an event through the registered publisher.

---

## Core API

### EventManager

```go
type EventManager struct {
    // contains filtered or unexported fields
}

func NewEventManager() *EventManager
```

`EventManager` manages publishers keyed by event type.

You normally create a single `EventManager` for your application (or per subsystem) and pass it around or resolve it from DI.

---

### EventObject[T]

```go
type EventObject[T any] interface {
    GetEventArgs() T
    StopPropagation() bool
    IsParallelPropagation() bool
}
```

* **GetEventArgs** – returns the payload that subscribers will receive.
* **StopPropagation** – if `true`, the publisher stops invoking further subscribers. Use this to short-circuit.
* **IsParallelPropagation** – if `true`, the publisher may dispatch subscribers in parallel (implementation-dependent; often `false` for simple events).

Any struct that implements these methods can be used as an event type.

---

### Subscriptions[T]

```go
type Subscriptions[T EventObject[T]] interface {
    SubscriptionsGetter[T]
    Add(subscribeFunc func(T)) error
    Remove(subscribeFunc func(T)) error
}

func NewSubscriptions[T EventObject[T]]() Subscriptions[T]
```

`Subscriptions` holds the set of subscriber functions for a given event type:

```go
subs := eventsmanager.NewSubscriptions[UserCreatedEvent]()

err := subs.Add(func(evt UserCreatedEvent) {
    fmt.Println("User created:", evt.UserID)
})
if err != nil {
    // handle bad signature errors, etc.
}
```

* `Add` – register a subscriber
* `Remove` – unregister it later if needed
* `SubscriptionsGetter[T]` – exposes `GetAlls() []func(T)` to inspect all subscriptions

---

### Publisher[T]

```go
type Publisher[T EventObject[T]] interface {
    Publish(event T)
}

func NewPublisher[T EventObject[T]](
    subscriptions SubscriptionsGetter[T],
    logger logs.Logger,
) Publisher[T]
```

`Publisher` is responsible for:

* Reading all subscribers from `SubscriptionsGetter[T]`
* Calling them with the event’s arguments (`GetEventArgs()`)
* Respecting `StopPropagation()` and `IsParallelPropagation()`
* Logging any errors through `logs.Logger`

You normally don’t implement `Publisher` yourself; you use `NewPublisher` and configure it with:

* A `Subscriptions[T]` (or anything implementing `SubscriptionsGetter[T]`)
* A logger (from the `logs` module)

---

### Top-level helpers: Register & Publish

```go
func Register[T EventObject[T]](em *EventManager, publisher Publisher[T])

func Publish[T EventObject[T]](em *EventManager, event T)
```

* `Register` associates a `Publisher[T]` with the `EventManager` for type `T`.
* `Publish` locates the appropriate publisher and dispatches the event.

This is the most convenient way to use the system:

```go
em := eventsmanager.NewEventManager()

subs := eventsmanager.NewSubscriptions[UserCreatedEvent]()
pub  := eventsmanager.NewPublisher(subs, logger)

eventsmanager.Register(em, pub)
eventsmanager.Publish(em, UserCreatedEvent{UserID: "123", Email: "user@example.com"})
```

---

## Error handling: SubscriptionsError

`SubscriptionsError` is a typed error for subscription operations:

```go
type SubscriptionsError interface {
    errors.CustomError
    GetErrorType() SubscriptionsErrorType
}

type SubscriptionsErrorType uint8

const (
    Unexpected SubscriptionsErrorType = iota
    BadFunctionSignature
    FunctionNoSubscribed
)
```

Meaning:

* `Unexpected` – something went wrong internally (panic recovery, etc.).
* `BadFunctionSignature` – e.g., `Add` was called with a function that doesn’t match `func(T)`.
* `FunctionNoSubscribed` – `Remove` was called with a function that isn’t currently subscribed.

Example:

```go
if err := subs.Add(func(evt UserCreatedEvent) {}); err != nil {
    if sErr, ok := err.(eventsmanager.SubscriptionsError); ok {
        switch sErr.GetErrorType() {
        case eventsmanager.BadFunctionSignature:
            logger.Error("bad subscriber signature:", sErr)
        default:
            logger.Error("subscriptions error:", sErr)
        }
    }
}
```

---

## DI integration (`eventsmanager/ioc`)

The `eventsmanager/ioc` package provides an IoC module for the DI container.

### EventsModule

```go
type EventsModule struct{}

func NewEventsModule() *EventsModule
func (m *EventsModule) RegisterServices(register dependencyinjection.Register) error
```

`EventsModule` implements `dependencyinjection.Module` and registers an `*EventManager` so you can inject it anywhere.

### Using with `dependencyinjection`

```go
import (
    di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
    eventsIoc "github.com/janmbaco/go-infrastructure/v2/eventsmanager/ioc"
    eventsmanager "github.com/janmbaco/go-infrastructure/v2/eventsmanager"
)

func main() {
    container := di.NewBuilder().
        AddModule(eventsIoc.NewEventsModule()).
        MustBuild()

    resolver := container.Resolver()

    // Resolve the shared EventManager:
    em := di.Resolve[*eventsmanager.EventManager](resolver)

    // Now register publishers and publish events as usual
    subs := eventsmanager.NewSubscriptions[UserCreatedEvent]()
    pub  := eventsmanager.NewPublisher(subs, logger)

    eventsmanager.Register(em, pub)
    eventsmanager.Publish(em, UserCreatedEvent{UserID: "1", Email: "user@example.com"})
}
```

This is the same `EventManager` used by other modules:

* `disk.FileChangedNotifier` publishes `FileChangedEvent` through it
* `configuration/events` publishes `ModifiedEvent` and `RestoredEvent`

So you can subscribe to file changes and config changes using the same pattern.

---

## Example: reacting to configuration & disk events

Because `configuration` and `disk` are built on top of `eventsmanager`, you can subscribe to their events just like any other custom event.

```go
import (
    "fmt"

    "github.com/janmbaco/go-infrastructure/v2/eventsmanager"
    cfgEvents "github.com/janmbaco/go-infrastructure/v2/configuration/events"
    "github.com/janmbaco/go-infrastructure/v2/configuration/fileconfig/ioc/resolver"
)

// Assume you already built a DI container with:
// - eventsmanager/ioc.NewEventsModule()
// - configuration/fileconfig/ioc.NewConfigurationModule()
// - logs, errors, disk, etc.

em := eventsResolver.GetEventManager(container.Resolver())
configHandler := cfgResolver.GetFileConfigHandler(
    container.Resolver(),
    "app.json",
    &AppConfig{Port: 8080},
)

// Subscribe to ModifiedEvent
cfgSubs := eventsmanager.NewSubscriptions[cfgEvents.ModifiedEvent]()
cfgSubs.Add(func(evt cfgEvents.ModifiedEvent) {
    cfg := configHandler.GetConfig().(*AppConfig)
    fmt.Println("Config changed; new port:", cfg.Port)
})

cfgPublisher := eventsmanager.NewPublisher(cfgSubs, logger)
eventsmanager.Register(em, cfgPublisher)

// File changes → disk publishes FileChangedEvent → config reloads → configuration/events publishes ModifiedEvent → your subscriber runs.
```

---

## Summary

The `eventsmanager` module gives you:

* A **generic, type-safe event system** based on Go generics
* Simple **subscriptions** with error checking
* **Publishers** that can respect parallelism & stop-propagation flags
* A central **EventManager** that wires event types to publishers
* Integration with other infrastructure modules (config, disk, etc.)
* An **IoC module** so you can inject and reuse the same `EventManager` across your app

Use it whenever you need decoupled, observable flows in your Go services, without falling back to `interface{}` or global variables.
