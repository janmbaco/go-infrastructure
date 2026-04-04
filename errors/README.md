# Errors

`github.com/janmbaco/go-infrastructure/v2/errors`

The `errors` package groups simple utilities for wrapping errors, validating arguments and handling error-returning functions in a consistent way.

## What It Includes

- `CustomError` and `CustomizableError`
- `ValidateNotNil`
- `ErrorCatcher`
- `errors/ioc` for DI registration

## Install

```bash
go get github.com/janmbaco/go-infrastructure/v2/errors
```

## Custom Errors

`CustomError` carries both a public message and an internal error:

```go
type CustomError interface {
    error
    GetMessage() string
    GetInternalError() error
}
```

The default implementation is `CustomizableError`:

```go
err := &errors.CustomizableError{
    Message:       "user not found",
    InternalError: fmt.Errorf("record not found"),
}
```

This is useful when you want a stable external message but still keep the original cause.

## Validate Not Nil

`ValidateNotNil` checks a set of named parameters and returns an error when one or more are nil:

```go
if err := errors.ValidateNotNil(map[string]interface{}{
    "user": user,
    "db":   db,
}); err != nil {
    return err
}
```

The returned error includes the calling function name and the missing parameter names.

## ErrorCatcher

`ErrorCatcher` is a helper for functions that return `error`. It does not recover panics; it only coordinates execution, optional callbacks and, in one case, logging.

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

Create it directly:

```go
logger := logs.NewLogger()
errorCatcher := errors.NewErrorCatcher(logger)
```

If `logger` is `nil`, the catcher falls back to a no-op logger.

## Common Patterns

Run a function and handle the returned error:

```go
err := errorCatcher.TryCatchError(
    func() error {
        return processFile("input.json")
    },
    func(err error) {
        fmt.Println("processing failed:", err)
    },
)
```

Always run cleanup:

```go
err := errorCatcher.TryCatchErrorAndFinally(
    func() error {
        return doWork()
    },
    func(err error) {
        fmt.Println("work failed:", err)
    },
    func() {
        fmt.Println("cleanup")
    },
)
```

Log an error and keep going:

```go
_ = errorCatcher.OnErrorContinue(func() error {
    return sendMetrics()
})
```

`OnErrorContinue` logs through `logs.ErrorLogger` and returns `nil`, even when the wrapped function fails.

## DI Integration

`errors/ioc` registers `errors.ErrorCatcher` as a singleton:

```go
container := di.NewBuilder().
    AddModule(logsioc.NewLogsModule()).
    AddModule(errorsioc.NewErrorsModule()).
    MustBuild()
```

Because `NewErrorCatcher` depends on `logs.ErrorLogger`, the usual setup is to register `logs/ioc` too.

## Related Files

- `customerror.go`: custom error contracts
- `validation.go`: nil validation helper
- `errorcatcher.go`: execution helpers
