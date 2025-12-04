# Logs

`github.com/janmbaco/go-infrastructure/v2/logs`

Production-oriented logging for Go applications, designed to work standalone or plugged into the `dependencyinjection` container.

- Log levels: `Trace`, `Info`, `Warning`, `Error`, `Fatal`
- Separate thresholds for console and file output
- Daily log file rotation: `logs/YYYY-MM-DD.log`
- Helper methods for common patterns (`Trace`, `Infof`, `TryError`…)
- DI module to inject a shared logger across your app

---

## Install

```bash
go get github.com/janmbaco/go-infrastructure/v2/logs
````

---

## Quick Start (standalone)

```go
package main

import "github.com/janmbaco/go-infrastructure/v2/logs"

func main() {
    logger := logs.NewLogger()

    // Where to write log files
    logger.SetDir("./logs")

    // Levels: what goes to console vs file
    logger.SetConsoleLevel(logs.Info)  // show Info+ on stdout
    logger.SetFileLogLevel(logs.Trace) // write everything to file

    logger.Info("app started")
    logger.Trace("debug details for local dev")

    // Fatal logs and then exits the process
    // logger.Fatal("unrecoverable error")
}
```

This gives you:

* Human-readable output on the console
* Daily-rotated log files under `./logs`
* A single API to use everywhere (`logs.Logger`)

---

## Log Levels

`LogLevel` is an `int`-based enum:

```go
const (
    Trace logs.LogLevel = iota
    Info
    Warning
    Error
    Fatal
)
```

Recommended usage:

| Level   | Typical use                                 |
| ------- | ------------------------------------------- |
| Trace   | Very detailed, noisy debug info             |
| Info    | Normal application flow (“user logged in…”) |
| Warning | Suspicious situations, degraded behavior    |
| Error   | Real failures that need attention           |
| Fatal   | Critical errors that terminate the process  |

You normally:

* Set **console** to `Info` (or `Warning`) in production
* Set **file** to `Trace` in production so you keep full history on disk

---

## Logger API

The main interface:

```go
type Logger interface {
    ErrorLogger

    Println(level LogLevel, message string)
    Printlnf(level LogLevel, format string, a ...interface{})

    Trace(message string)
    Tracef(format string, a ...interface{})
    Info(message string)
    Infof(format string, a ...interface{})
    Warning(message string)
    Warningf(format string, a ...interface{})
    Error(message string)
    Errorf(format string, a ...interface{})
    Fatal(message string)
    Fatalf(format string, a ...interface{})

    SetConsoleLevel(level LogLevel)
    SetFileLogLevel(level LogLevel)

    GetErrorLogger() *log.Logger

    SetDir(dir string)
    Mute()
    Unmute()
}
```

### Level-specific helpers

```go
logger.Trace("starting handler X")
logger.Infof("user %s logged in", username)
logger.Warning("cache miss, falling back to DB")
logger.Errorf("failed to connect to DB: %v", err)
logger.Fatal("configuration is invalid") // logs & exits
```

Under the hood these call `Println` / `Printlnf` with the matching `LogLevel`.

### Controlling outputs

```go
// Only Warning+ on console
logger.SetConsoleLevel(logs.Warning)

// But keep everything in the file
logger.SetFileLogLevel(logs.Trace)
```

* Messages **below** the configured level are **ignored** for that output (console/file).
* You can tune console and file independently.

### Log directory & rotation

```go
logger.SetDir("./logs")
```

* Log entries are written to files like:

  ```text
  ./logs/2025-12-01.log
  ./logs/2025-12-02.log
  ```

* Rotation is **daily**; each day goes to its own file.

Make sure:

* The directory exists and your process can write there.
* For Docker/Kubernetes, mount a volume for that path if you want persistence.

### Mute / Unmute

Sometimes you want to temporarily silence logs (e.g., noisy integration tests):

```go
logger.Mute()
defer logger.Unmute()

// These calls become no-ops while muted
logger.Info("this won’t be written")
logger.Error("neither will this")
```

Call `Unmute()` to restore normal behavior.

---

## ErrorLogger shortcuts

`Logger` embeds `ErrorLogger`, which focuses on logging `error` values:

```go
type ErrorLogger interface {
    PrintError(level LogLevel, err error)
    TryPrintError(level LogLevel, err error)

    TryTrace(err error)
    TryInfo(err error)
    TryWarning(err error)
    TryError(err error)
    TryFatal(err error)
}
```

### Common patterns

#### 1. Log an error at a specific level

```go
if err != nil {
    logger.PrintError(logs.Error, err)
}
```

#### 2. Log only if the error is non-nil

```go
// This checks for nil internally
logger.TryError(err)
```

This is equivalent to:

```go
if err != nil {
    logger.Error(err.Error())
}
```

…but shorter and less error-prone.

#### 3. “Fire and forget” error logging

```go
if err := doWork(); err != nil {
    logger.TryWarning(err) // log but keep running
}
```

You’ll use these a lot with the `errors` module and `ErrorCatcher`.

---

## Using with Dependency Injection (logs/ioc)

The `logs/ioc` package provides a DI module so you can inject a shared logger anywhere.

### Register the module

```go
import (
    di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
    logsIoc "github.com/janmbaco/go-infrastructure/v2/logs/ioc"
    "github.com/janmbaco/go-infrastructure/v2/logs"
)

func main() {
    container := di.NewBuilder().
        AddModule(logsIoc.NewLogsModule()).
        MustBuild()

    resolver := container.Resolver()

    logger := di.Resolve[logs.Logger](resolver)
    logger.Info("Logger resolved from DI")
}
```

`LogsModule`:

* Registers a **singleton** `logs.Logger` for the whole container
* Makes it available for constructor injection wherever you need logging

### Injecting into your own services

```go
type UserService struct {
    logger logs.Logger
}

func NewUserService(logger logs.Logger) *UserService {
    return &UserService{logger: logger}
}

func (s *UserService) CreateUser(name string) error {
    s.logger.Infof("creating user %s", name)
    // ...
    return nil
}
```

Register with DI:

```go
container := di.NewBuilder().
    AddModule(logsIoc.NewLogsModule()).
    Register(func(r di.Register) {
        di.RegisterSingleton[*UserService](r, func() *UserService {
            return &UserService{
                logger: di.Resolve[logs.Logger](r.Resolver()),
            }
        })
    }).
    MustBuild()

svc := di.Resolve[*UserService](container.Resolver())
svc.CreateUser("alice")
```

Everywhere you resolve `logs.Logger` you get the **same instance**, so:

* Settings like `SetDir`, `SetConsoleLevel`, `SetFileLogLevel` apply app-wide
* Logs from all services go to the same files

---

## Working with the `errors` module

The `errors` package’s `ErrorCatcher` is built on top of `logs.ErrorLogger`, so using `LogsModule` gives it what it needs automatically.

Example:

```go
import (
    "fmt"

    di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
    errorsIoc "github.com/janmbaco/go-infrastructure/v2/errors/ioc"
    errorsResolver "github.com/janmbaco/go-infrastructure/v2/errors/ioc/resolver"
    logsIoc "github.com/janmbaco/go-infrastructure/v2/logs/ioc"
)

func riskyOperation() error {
    return fmt.Errorf("something went wrong")
}

func main() {
    container := di.NewBuilder().
        AddModule(logsIoc.NewLogsModule()).
        AddModule(errorsIoc.NewErrorsModule()).
        MustBuild()

    errorCatcher := errorsResolver.GetErrorCatcher(container.Resolver())

    errorCatcher.TryCatchError(
        riskyOperation,
        func(err error) {
            // err is already logged using the shared logger
            fmt.Println("Operation failed:", err)
        },
    )
}
```

No extra wiring: as long as you add both modules, the `ErrorCatcher` logs to the same logger used everywhere else.

---

## Integrating with libraries that need *log.Logger

Sometimes a third-party library expects a `*log.Logger`.
Use `GetErrorLogger()` to hook the go-infrastructure logger into it:

```go
import (
    stdlog "log"

    "github.com/janmbaco/go-infrastructure/v2/logs"
)

func main() {
    logger := logs.NewLogger()
    logger.SetDir("./logs")

    stdLogger := logger.GetErrorLogger()

    // Example: some library accepts *log.Logger
    somePkg.SetLogger(stdLogger)
}
```

This keeps all logging (your code + the library) going through the same log files and levels.

---

## Testing

For unit tests you often don’t want real files, or you just want to assert messages.

You can:

1. **Use the real logger but mute it**

   ```go
   func TestSomething(t *testing.T) {
       logger := logs.NewLogger()
       logger.Mute() // avoid noisy test output

       // pass logger into your services, or resolve from DI
   }
   ```

2. **Swap in a fake implementation**

   ```go
   type fakeLogger struct{}

   func (fakeLogger) PrintError(level logs.LogLevel, err error)              {}
   func (fakeLogger) TryPrintError(level logs.LogLevel, err error)          {}
   func (fakeLogger) TryTrace(err error)                                     {}
   func (fakeLogger) TryInfo(err error)                                      {}
   func (fakeLogger) TryWarning(err error)                                   {}
   func (fakeLogger) TryError(err error)                                     {}
   func (fakeLogger) TryFatal(err error)                                     {}

   func (fakeLogger) Println(level logs.LogLevel, message string)            {}
   func (fakeLogger) Printlnf(level logs.LogLevel, format string, a ...any)  {}
   func (fakeLogger) Trace(message string)                                   {}
   func (fakeLogger) Tracef(format string, a ...any)                         {}
   func (fakeLogger) Info(message string)                                    {}
   func (fakeLogger) Infof(format string, a ...any)                          {}
   func (fakeLogger) Warning(message string)                                 {}
   func (fakeLogger) Warningf(format string, a ...any)                       {}
   func (fakeLogger) Error(message string)                                   {}
   func (fakeLogger) Errorf(format string, a ...any)                         {}
   func (fakeLogger) Fatal(message string)                                   {}
   func (fakeLogger) Fatalf(format string, a ...any)                         {}
   func (fakeLogger) SetConsoleLevel(level logs.LogLevel)                    {}
   func (fakeLogger) SetFileLogLevel(level logs.LogLevel)                    {}
   func (fakeLogger) GetErrorLogger() *stdlog.Logger                         { return stdlog.Default() }
   func (fakeLogger) SetDir(dir string)                                      {}
   func (fakeLogger) Mute()                                                  {}
   func (fakeLogger) Unmute()                                                {}
   ```

Then register it in DI for tests instead of using `LogsModule`.

---

## Summary

* Use `logs.NewLogger()` when you just need a logger.
* Use `logs/ioc.NewLogsModule()` when you want logging integrated into the DI container.
* Configure behavior with `SetConsoleLevel`, `SetFileLogLevel`, and `SetDir`.
* Use `TryError`, `TryWarning`, etc., to log `error` values safely.
* Combine with the `errors` module for centralized error capture and logging.

