# Events Manager

`github.com/janmbaco/go-infrastructure/v2/eventsmanager`

The `eventsmanager` package is a small generic event bus. It stores subscriptions by event type, publishes events through registered publishers and can recover panics in event handlers.

## What It Includes

- `EventObject[T]`
- `Subscriptions[T]`
- `Publisher[T]`
- `EventManager`
- top-level `Register` and `Publish` helpers
- `eventsmanager/ioc` for DI registration

## Install

```bash
go get github.com/janmbaco/go-infrastructure/v2/eventsmanager
```

## Event Contract

Each event type implements `EventObject[T]`:

```go
type EventObject[T any] interface {
    GetEventArgs() T
    StopPropagation() bool
    IsParallelPropagation() bool
}
```

Example:

```go
type UserCreatedEvent struct {
    UserID string
    Email  string
}

func (e UserCreatedEvent) GetEventArgs() UserCreatedEvent { return e }
func (e UserCreatedEvent) StopPropagation() bool         { return false }
func (e UserCreatedEvent) IsParallelPropagation() bool   { return false }
```

## Quick Start

```go
logger := logs.NewLogger()

subscriptions := eventsmanager.NewSubscriptions[UserCreatedEvent]()
_ = subscriptions.Add(func(event UserCreatedEvent) {
    fmt.Println("welcome:", event.Email)
})

publisher := eventsmanager.NewPublisher(subscriptions, logger)

eventManager := eventsmanager.NewEventManager()
eventsmanager.Register(eventManager, publisher)

eventsmanager.Publish(eventManager, UserCreatedEvent{
    UserID: "123",
    Email:  "user@example.com",
})
```

## Subscriptions

```go
type Subscriptions[T EventObject[T]] interface {
    SubscriptionsGetter[T]
    Add(subscribeFunc func(T)) error
    Remove(subscribeFunc func(T)) error
}
```

`Add` stores a handler for the event type. `Remove` returns `FunctionNoSubscribed` if the function was not registered.

## Publisher

```go
type Publisher[T EventObject[T]] interface {
    Publish(event T)
}
```

The default publisher:

- reads all registered handlers
- dispatches sequentially or in parallel depending on `IsParallelPropagation()`
- stops early when `StopPropagation()` returns `true`
- recovers panics from handlers and logs them through `logs.Logger`

## EventManager

```go
type EventManager struct {
    // unexported fields
}

func NewEventManager() *EventManager
func Register[T EventObject[T]](em *EventManager, publisher Publisher[T])
func Publish[T EventObject[T]](em *EventManager, event T)
```

`EventManager` is just the registry that maps event types to publishers.

## Subscription Errors

```go
type SubscriptionsError interface {
    errors.CustomError
    GetErrorType() SubscriptionsErrorType
}

const (
    Unexpected SubscriptionsErrorType = iota
    BadFunctionSignature
    FunctionNoSubscribed
)
```

In the current implementation, `Remove` is the main path that can return a subscription error.

## DI Integration

`eventsmanager/ioc` registers `*eventsmanager.EventManager` as a singleton:

```go
container := di.NewBuilder().
    AddModule(eventsioc.NewEventsModule()).
    MustBuild()
```

Other modules such as `configuration` and `disk` rely on this registration.

## Related Files

- `eventobject.go`: event contract
- `subscriptions.go`: subscription storage
- `publisher.go`: dispatch logic
- `eventmanager.go`: publisher registry
