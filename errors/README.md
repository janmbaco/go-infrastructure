# Errors

`github.com/janmbaco/go-infrastructure/v2/errors`

Centralized error utilities for Go applications:

- `ErrorCatcher` – “try/catch/finally” style helpers **without** panic/recover
- `ValidateNotNil` – defensive parameter validation with readable messages
- `CustomError` / `CustomizableError` – wrap internal errors with user-friendly messages
- IoC module to inject a shared `ErrorCatcher` across your app

Use it standalone, or plug it into the `dependencyinjection` container alongside `logs`.

---

## Install

```bash
go get github.com/janmbaco/go-infrastructure/v2/errors
````

If you use DI:

```bash
go get github.com/janmbaco/go-infrastructure/v2/errors/ioc
go get github.com/janmbaco/go-infrastructure/v2/errors/ioc/resolver
```

---

## Quick Start: Centralized error handling

```go
package main

import (
    "fmt"

    di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
    errorsIoc "github.com/janmbaco/go-infrastructure/v2/errors/ioc"
    errorsResolver "github.com/janmbaco/go-infrastructure/v2/errors/ioc/resolver"
    logsIoc "github.com/janmbaco/go-infrastructure/v2/logs/ioc"
)

func processData(data []byte) error {
    if len(data) == 0 {
        return fmt.Errorf("empty data")
    }
    // Process...
    return nil
}

func main() {
    // Build container with logging + error handling
    container := di.NewBuilder().
        AddModule(logsIoc.NewLogsModule()).
        AddModule(errorsIoc.NewErrorsModule()).
        MustBuild()

    resolver := container.Resolver()

    // Resolve shared ErrorCatcher
    errorCatcher := errorsResolver.GetErrorCatcher(resolver)

    // Execute logic with centralized error handling & logging
    errorCatcher.TryCatchError(
        func() error { return processData([]byte{}) },
        func(err error) {
            // Error is already logged by the ErrorCatcher
            fmt.Println("Operation failed:", err)
        },
    )
}
```

What you get:

* No panic/recover
* One place to deal with logging & error formatting
* Reusable pattern you can apply across your app

---

## Core Types

### `CustomError` and `CustomizableError`

`CustomError` lets you carry **two layers** of information:

* A human-friendly message
* The underlying internal error

```go
import "github.com/janmbaco/go-infrastructure/v2/errors"

func loadUser(id string) error {
    // Imagine this comes from a DB client
    internalErr := fmt.Errorf("record not found")

    return &errors.CustomizableError{
        Message:       "user not found",
        InternalError: internalErr,
    }
}
```

The interface:

```go
type CustomError interface {
    error
    GetMessage() string
    GetInternalError() error
}
```

Typical usage:

* Wrap low-level errors (DB, HTTP, etc.) with a clear message for logs/UX
* Keep the original error for troubleshooting or error inspection

Example:

```go
err := loadUser("123")
if cErr, ok := err.(errors.CustomError); ok {
    fmt.Println("Public message:", cErr.GetMessage())
    fmt.Println("Internal error:", cErr.GetInternalError())
}
```

---

### `ValidateNotNil`

Simple defensive programming helper for function parameters:

```go
import "github.com/janmbaco/go-infrastructure/v2/errors"

func ProcessUser(user *User, db *Database) error {
    // This returns an error if any of these values is nil
    if err := errors.ValidateNotNil(map[string]interface{}{
        "user": user,
        "db":   db,
    }); err != nil {
        return err
    }

    // Safe to use user and db here
    return nil
}
```

Key points:

* Accepts a map of `"parameterName" -> value`
* Returns `nil` if everything is non-nil
* Returns an error listing which parameters were `nil` otherwise

Use it at the top of public functions or constructors to fail fast on programming errors.

---

## ErrorCatcher

`ErrorCatcher` gives you structured ways to **run functions that return `error`** and handle them consistently:

```go
type ErrorCatcher interface {
    HandleError(err error, errorfn func(error)) error
    HandleErrorWithFinally(err error, errorfn func(error), finallyfn func()) error
    TryCatchError(tryfn func() error, errorfn func(error)) error
    TryFinally(tryfn func() error, finallyfn func()) error
    TryCatchErrorAndFinally(tryfn func() error, errorfn func(error), finallyfn func()) error
    OnErrorContinue(tryfn func() error) error
    CatchError(err error, errorfn func(error)) error
    CatchErrorAndFinally(err error, errorfn func(error), finallyfn func()) error
}
```

The concrete implementation is created with:

```go
func NewErrorCatcher(logger logs.ErrorLogger) errors.ErrorCatcher
```

The logger is used internally to record errors (including stack traces, depending on your logging setup).

### Patterns

#### 1. `TryCatchError`: classic try/catch

```go
err := errorCatcher.TryCatchError(
    func() error {
        return doWork()
    },
    func(err error) {
        // Handle the error (already logged)
        fmt.Println("doWork failed:", err)
    },
)
if err != nil {
    // Optional: propagate or wrap
}
```

Good for:

* Running operations that may fail
* Guaranteeing logging and uniform error handling logic

#### 2. `TryCatchErrorAndFinally`: try/catch/finally

```go
err := errorCatcher.TryCatchErrorAndFinally(
    func() error {
        return openAndProcessFile("input.txt")
    },
    func(err error) {
        fmt.Println("Processing failed:", err)
    },
    func() {
        fmt.Println("Cleaning up resources...")
        // close files, release locks, etc.
    },
)
```

* `tryfn` – your main operation
* `errorfn` – runs when `tryfn` returns an error
* `finallyfn` – runs **always**, success or failure

#### 3. `TryFinally`: cleanup regardless of errors

```go
err := errorCatcher.TryFinally(
    func() error {
        return updateCache()
    },
    func() {
        fmt.Println("Finished attempting cache update.")
    },
)
```

Use this when:

* You don’t need custom error handling at the call site
* You still want a `finally` step (logging, metric, cleanup)

#### 4. `OnErrorContinue`: log and keep going

```go
_ = errorCatcher.OnErrorContinue(func() error {
    return sendMetrics()
})
// Even if sendMetrics fails, execution continues
```

This is handy for non-critical operations where failures should be logged but not stop the flow (background metrics, telemetry, optional notifications, etc.).

#### 5. `HandleError` / `CatchError`: deal with existing error values

If you already have an `err` from some call:

```go
err := repo.Save(user)

err = errorCatcher.HandleError(err, func(e error) {
    // Only called when err != nil
    fmt.Println("Save failed:", e)
})
if err != nil {
    return err
}
```

`CatchError` / `CatchErrorAndFinally` offer the same semantics but named from the “catch” perspective:

```go
err := repo.Save(user)

return errorCatcher.CatchErrorAndFinally(
    err,
    func(e error) { fmt.Println("Save failed:", e) },
    func()       { fmt.Println("Done trying to save user") },
)
```

---

## Using with Dependency Injection (`errors/ioc`)

The `errors/ioc` package gives you a DI module to register `ErrorCatcher` as part of the container.

### Register the module

```go
import (
    di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
    errorsIoc "github.com/janmbaco/go-infrastructure/v2/errors/ioc"
    logsIoc "github.com/janmbaco/go-infrastructure/v2/logs/ioc"
)

container := di.NewBuilder().
    AddModule(logsIoc.NewLogsModule()).
    AddModule(errorsIoc.NewErrorsModule()).
    MustBuild()
```

`ErrorsModule`:

* Implements `dependencyinjection.Module`
* Registers an `errors.ErrorCatcher` that uses the shared `logs.ErrorLogger` from the `logs` module

### Resolving `ErrorCatcher` with the resolver helper

```go
import (
    errorsResolver "github.com/janmbaco/go-infrastructure/v2/errors/ioc/resolver"
)

resolver := container.Resolver()
errorCatcher := errorsResolver.GetErrorCatcher(resolver)
```

Signature:

```go
func GetErrorCatcher(resolver dependencyinjection.Resolver) errors.ErrorCatcher
```

You can also resolve it directly using generics if you prefer:

```go
import (
    "github.com/janmbaco/go-infrastructure/v2/errors"
    di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
)

errorCatcher := di.Resolve[errors.ErrorCatcher](resolver)
```

Either way, it’s the same singleton instance registered by `ErrorsModule`.

---

## Testing error flows

In unit tests you often want to:

* Assert that error callbacks are invoked
* Avoid writing real logs

You have two options:

1. **Use the real `ErrorCatcher` with a fake logger**

   ```go
   type fakeErrorLogger struct {
       lastError error
   }

   func (f *fakeErrorLogger) TryError(err error) {
       f.lastError = err
   }

   // implement the other methods of logs.ErrorLogger as no-ops...

   func TestTryCatchError(t *testing.T) {
       logger := &fakeErrorLogger{}
       catcher := errors.NewErrorCatcher(logger)

       called := false
       catcher.TryCatchError(
           func() error { return fmt.Errorf("boom") },
           func(err error) {
               called = true
           },
       )

       if !called {
           t.Fatal("expected error callback to be called")
       }
       if logger.lastError == nil {
           t.Fatal("expected error to be logged")
       }
   }
   ```

2. **Inject a fake `ErrorCatcher` via DI**

   For higher-level tests, you can register your own implementation of `errors.ErrorCatcher` in the container and skip `ErrorsModule` entirely for that test.

---

## Summary

The `errors` package gives you:

* A unified way to **validate inputs** (`ValidateNotNil`)
* A small abstraction for **wrapping errors** with user-friendly messages (`CustomError`, `CustomizableError`)
* A powerful `ErrorCatcher` interface to implement consistent **try/catch/finally** patterns without resorting to panic/recover
* A ready-to-use **DI module** so you can inject a shared `ErrorCatcher` wired to your logger

Combine it with the `logs` and `dependencyinjection` modules to keep error handling consistent, observable, and easy to test across your whole application.
