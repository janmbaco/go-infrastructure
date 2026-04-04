# Logs

`github.com/janmbaco/go-infrastructure/v2/logs`

The `logs` package provides a small logger with level-based console output, file output and a companion `ErrorLogger` interface used by other modules in this repository.

## What It Includes

- `Logger` with `Trace`, `Info`, `Warning`, `Error` and `Fatal`
- configurable console and file thresholds
- file logging under a configurable directory
- `ErrorLogger` helpers such as `TryError`
- `logs/ioc` for DI registration

## Install

```bash
go get github.com/janmbaco/go-infrastructure/v2/logs
```

## Log Levels

```go
const (
    Trace LogLevel = iota
    Info
    Warning
    Error
    Fatal
)
```

## Quick Start

```go
package main

import "github.com/janmbaco/go-infrastructure/v2/logs"

func main() {
    logger := logs.NewLogger()

    logger.SetDir("./logs")
    logger.SetConsoleLevel(logs.Info)
    logger.SetFileLogLevel(logs.Trace)

    logger.Info("application started")
    logger.Trace("debug details")
}
```

If a log directory is configured, log files are written with the executable name as prefix and a date-based suffix.

## Logger Interface

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
    SetDir(string)
    Mute()
    Unmute()
}
```

`Fatal` and `Fatalf` call the underlying logger fatal path and terminate the process.

## ErrorLogger

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

Example:

```go
if err := doWork(); err != nil {
    logger.TryError(err)
}
```

## Output Control

Console and file thresholds are configured independently:

```go
logger.SetConsoleLevel(logs.Warning)
logger.SetFileLogLevel(logs.Trace)
```

`Mute()` and `Unmute()` disable and restore all output temporarily.

## DI Integration

`logs/ioc` registers `logs.Logger` as a singleton and binds `logs.ErrorLogger` to the same implementation:

```go
container := di.NewBuilder().
    AddModule(logsioc.NewLogsModule()).
    MustBuild()

logger := di.Resolve[logs.Logger](container.Resolver())
```

That binding is what allows modules such as `errors` to depend on `logs.ErrorLogger`.

## Related Files

- `logger.go`: core implementation
- `ioc/module.go`: DI registration
