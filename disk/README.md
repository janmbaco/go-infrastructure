# Disk

`github.com/janmbaco/go-infrastructure/v2/disk`

The `disk` package provides small file utilities and a file change notifier used by the configuration system.

## What It Includes

- `ExistsPath`
- `CreateFile`
- `Copy`
- `DeleteFile`
- `FileChangedNotifier`
- `disk/ioc` for DI registration

## Install

```bash
go get github.com/janmbaco/go-infrastructure/v2/disk
```

## File Utilities

```go
exists := disk.ExistsPath("config.json")

err := disk.CreateFile("config.json", []byte(`{"port":":8080"}`))
err = disk.Copy("config.json", "config.backup.json")
err = disk.DeleteFile("config.backup.json")

_, _ = exists, err
```

These helpers are thin wrappers around common `os` and `io` file operations.

## FileChangedNotifier

`FileChangedNotifier` is a callback-based watcher for write events on a single file:

```go
type FileChangedNotifier interface {
    Subscribe(subscribeFunc func()) error
}
```

Create it directly:

```go
eventManager := eventsmanager.NewEventManager()
logger := logs.NewLogger()

notifier, err := disk.NewFileChangedNotifier("config.json", eventManager, logger)
if err != nil {
    panic(err)
}

if err := notifier.Subscribe(func() {
    fmt.Println("config.json changed")
}); err != nil {
    panic(err)
}
```

The watcher starts lazily on the first subscription and publishes an internal `FileChangedEvent` for write operations.

## DI Integration

`disk/ioc` registers `disk.FileChangedNotifier` with `filePath` as a named parameter:

```go
container := di.NewBuilder().
    AddModule(logsioc.NewLogsModule()).
    AddModule(eventsioc.NewEventsModule()).
    AddModule(diskioc.NewDiskModule()).
    MustBuild()

notifier := di.ResolveWithParams[disk.FileChangedNotifier](
    container.Resolver(),
    map[string]interface{}{"filePath": "config.json"},
)
```

Because `NewFileChangedNotifier` depends on `*eventsmanager.EventManager` and `logs.Logger`, those modules should be registered too when using DI.

## Related Files

- `path.go`: file helpers
- `filechangednotifier.go`: fsnotify-backed notifier
- `ioc/module.go`: DI registration
