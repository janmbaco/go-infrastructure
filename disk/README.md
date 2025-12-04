# Disk Module

`github.com/janmbaco/go-infrastructure/v2/disk`

File system utilities with path normalization and file change notifications.

## Overview

The disk module provides essential file system operations:

- **Path normalization** - Cross-platform path handling
- **File change notifications** - Watch files for modifications
- **FD limit management** - Handle file descriptor limits on Unix systems
- **Safe path operations** - Prevent path traversal vulnerabilities

## Installation

```bash
go get github.com/janmbaco/go-infrastructure/v2/disk
```

## Quick Start

### Path Operations

```go
package main

import (
    "fmt"
    "github.com/janmbaco/go-infrastructure/v2/disk"
)

func main() {
    // Normalize paths across platforms
    path := disk.NormalizePath("./config/../data/file.txt")
    fmt.Println("Normalized:", path) // Output: data/file.txt

    // Clean and resolve paths
    absPath := disk.GetAbsolutePath("./configs/app.json")
    fmt.Println("Absolute:", absPath)
}
```

### File Change Notifications

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/janmbaco/go-infrastructure/v2/disk"
)

func main() {
    // Create notifier for a config file
    notifier := disk.NewFileChangedNotifier("./config.json")

    // Subscribe to changes
    notifier.Subscribe(func(event disk.FileChangedEvent) {
        fmt.Printf("File changed: %s at %s\n", event.Path, event.Timestamp)
        
        // Reload configuration or perform other actions
        reloadConfig(event.Path)
    })

    // Start watching
    if err := notifier.Start(); err != nil {
        panic(err)
    }

    // Keep running
    time.Sleep(time.Hour)
    
    // Cleanup
    notifier.Stop()
}

func reloadConfig(path string) {
    // Your config reload logic
}
```

## API Reference

### Path Operations

```go
// NormalizePath normalizes a file path to use forward slashes
// and removes .. and . components
func NormalizePath(path string) string

// GetAbsolutePath returns the absolute path of a file
func GetAbsolutePath(path string) string

// EnsureDirectory creates a directory if it doesn't exist
func EnsureDirectory(path string) error

// IsPathSafe checks if a path is safe (doesn't escape base directory)
func IsPathSafe(basePath, targetPath string) bool
```

### FileChangedNotifier

```go
type FileChangedNotifier interface {
    // Subscribe registers a callback for file change events
    Subscribe(handler func(FileChangedEvent))
    
    // Start begins watching the file for changes
    Start() error
    
    // Stop stops watching the file
    Stop()
}

type FileChangedEvent struct {
    Path      string    // Path to the changed file
    Timestamp time.Time // When the change occurred
}

// NewFileChangedNotifier creates a new file change notifier
func NewFileChangedNotifier(filePath string) FileChangedNotifier
```

## Usage Examples

### Configuration File Watcher

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "sync"
    
    "github.com/janmbaco/go-infrastructure/v2/disk"
)

type Config struct {
    mu       sync.RWMutex
    data     map[string]interface{}
    notifier disk.FileChangedNotifier
}

func NewConfig(path string) (*Config, error) {
    cfg := &Config{
        data:     make(map[string]interface{}),
        notifier: disk.NewFileChangedNotifier(path),
    }

    // Load initial config
    if err := cfg.load(path); err != nil {
        return nil, err
    }

    // Watch for changes
    cfg.notifier.Subscribe(func(event disk.FileChangedEvent) {
        fmt.Println("Config changed, reloading...")
        if err := cfg.load(event.Path); err != nil {
            fmt.Println("Error reloading:", err)
        }
    })

    if err := cfg.notifier.Start(); err != nil {
        return nil, err
    }

    return cfg, nil
}

func (c *Config) load(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return err
    }

    var newData map[string]interface{}
    if err := json.Unmarshal(data, &newData); err != nil {
        return err
    }

    c.mu.Lock()
    c.data = newData
    c.mu.Unlock()

    return nil
}

func (c *Config) Get(key string) interface{} {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.data[key]
}

func (c *Config) Close() {
    c.notifier.Stop()
}
```

### Safe File Access

```go
package main

import (
    "errors"
    "os"
    "path/filepath"
    
    "github.com/janmbaco/go-infrastructure/v2/disk"
)

type SafeFileAccess struct {
    baseDir string
}

func NewSafeFileAccess(baseDir string) (*SafeFileAccess, error) {
    absBase := disk.GetAbsolutePath(baseDir)
    
    // Ensure base directory exists
    if err := disk.EnsureDirectory(absBase); err != nil {
        return nil, err
    }

    return &SafeFileAccess{baseDir: absBase}, nil
}

func (s *SafeFileAccess) ReadFile(relativePath string) ([]byte, error) {
    // Normalize the path
    normalizedPath := disk.NormalizePath(relativePath)
    
    // Build full path
    fullPath := filepath.Join(s.baseDir, normalizedPath)
    
    // Check if path is safe (doesn't escape base directory)
    if !disk.IsPathSafe(s.baseDir, fullPath) {
        return nil, errors.New("path traversal detected")
    }

    // Read file
    return os.ReadFile(fullPath)
}

func (s *SafeFileAccess) WriteFile(relativePath string, data []byte) error {
    normalizedPath := disk.NormalizePath(relativePath)
    fullPath := filepath.Join(s.baseDir, normalizedPath)

    if !disk.IsPathSafe(s.baseDir, fullPath) {
        return errors.New("path traversal detected")
    }

    // Ensure directory exists
    dir := filepath.Dir(fullPath)
    if err := disk.EnsureDirectory(dir); err != nil {
        return err
    }

    return os.WriteFile(fullPath, data, 0644)
}
```

### Multi-File Watcher

```go
package main

import (
    "fmt"
    "sync"
    
    "github.com/janmbaco/go-infrastructure/v2/disk"
)

type MultiFileWatcher struct {
    notifiers []disk.FileChangedNotifier
    mu        sync.Mutex
}

func NewMultiFileWatcher(paths ...string) *MultiFileWatcher {
    watcher := &MultiFileWatcher{
        notifiers: make([]disk.FileChangedNotifier, 0, len(paths)),
    }

    for _, path := range paths {
        notifier := disk.NewFileChangedNotifier(path)
        notifier.Subscribe(func(event disk.FileChangedEvent) {
            watcher.handleChange(event)
        })
        watcher.notifiers = append(watcher.notifiers, notifier)
    }

    return watcher
}

func (w *MultiFileWatcher) Start() error {
    for _, notifier := range w.notifiers {
        if err := notifier.Start(); err != nil {
            w.StopAll()
            return err
        }
    }
    return nil
}

func (w *MultiFileWatcher) StopAll() {
    for _, notifier := range w.notifiers {
        notifier.Stop()
    }
}

func (w *MultiFileWatcher) handleChange(event disk.FileChangedEvent) {
    w.mu.Lock()
    defer w.mu.Unlock()
    
    fmt.Printf("File changed: %s\n", event.Path)
    // Handle the change
}
```

### Integration with DI Container

```go
package main

import (
    di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
    "github.com/janmbaco/go-infrastructure/v2/disk"
    "github.com/janmbaco/go-infrastructure/v2/disk/ioc"
)

func main() {
    container := di.NewBuilder().
        AddModule(ioc.NewDiskModule()).
        MustBuild()

    // Path operations are typically used directly
    // without dependency injection, but you can access
    // them if needed through the module

    // File notifiers are usually created as needed
    notifier := disk.NewFileChangedNotifier("./config.json")
    // ... use notifier
}
```

## Platform-Specific Features

### File Descriptor Limits (Unix)

The `fdlimit` package helps manage file descriptor limits on Unix systems:

```go
package main

import (
    "fmt"
    "runtime"
    
    "github.com/janmbaco/go-infrastructure/v2/disk/fdlimit"
)

func main() {
    if runtime.GOOS != "windows" {
        // Get current limit
        current, max, err := fdlimit.GetLimit()
        if err != nil {
            panic(err)
        }
        fmt.Printf("Current FD limit: %d, Max: %d\n", current, max)

        // Raise limit
        newLimit := uint64(10000)
        if err := fdlimit.RaiseLimit(newLimit); err != nil {
            fmt.Println("Failed to raise limit:", err)
        } else {
            fmt.Println("FD limit raised to:", newLimit)
        }
    }
}
```

## Security Best Practices

### Path Traversal Prevention

```go
// BAD - Vulnerable to path traversal
func readFile(userPath string) ([]byte, error) {
    return os.ReadFile(userPath)
}

// GOOD - Protected against path traversal
func readFile(baseDir, userPath string) ([]byte, error) {
    normalizedPath := disk.NormalizePath(userPath)
    fullPath := filepath.Join(baseDir, normalizedPath)
    
    if !disk.IsPathSafe(baseDir, fullPath) {
        return nil, errors.New("invalid path")
    }
    
    return os.ReadFile(fullPath)
}
```

### Safe Path Construction

```go
// Always normalize user input
userInput := "../../../etc/passwd"
safePath := disk.NormalizePath(userInput)

// Validate against base directory
baseDir := "/app/data"
fullPath := filepath.Join(baseDir, safePath)

if !disk.IsPathSafe(baseDir, fullPath) {
    return errors.New("path traversal attempt detected")
}
```

## Performance Considerations

### File Change Notifications

- **Polling interval:** Default 1 second (configurable)
- **Resource usage:** One goroutine per watched file
- **Debouncing:** Built-in to prevent duplicate events

### Best Practices

```go
// For many files, consider batching
type BatchWatcher struct {
    files map[string]time.Time
    mu    sync.Mutex
}

func (b *BatchWatcher) CheckChanges() {
    b.mu.Lock()
    defer b.mu.Unlock()
    
    for path, lastMod := range b.files {
        info, err := os.Stat(path)
        if err != nil {
            continue
        }
        
        if info.ModTime().After(lastMod) {
            b.files[path] = info.ModTime()
            // Handle change
        }
    }
}
```

## Error Handling

```go
// Check if file exists before watching
if _, err := os.Stat(filePath); os.IsNotExist(err) {
    return fmt.Errorf("file does not exist: %s", filePath)
}

// Handle notifier errors
notifier := disk.NewFileChangedNotifier(filePath)
if err := notifier.Start(); err != nil {
    return fmt.Errorf("failed to start watching: %w", err)
}

// Always cleanup
defer notifier.Stop()
```

## Testing

Example tests:

```go
func TestNormalizePath(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {"./foo/../bar", "bar"},
        {"/foo/./bar", "/foo/bar"},
        {"foo\\bar", "foo/bar"}, // Windows path
    }

    for _, tt := range tests {
        result := disk.NormalizePath(tt.input)
        assert.Equal(t, tt.expected, result)
    }
}

func TestFileChangedNotifier(t *testing.T) {
    tmpFile := createTempFile(t)
    defer os.Remove(tmpFile)

    notifier := disk.NewFileChangedNotifier(tmpFile)
    
    changed := make(chan struct{})
    notifier.Subscribe(func(event disk.FileChangedEvent) {
        changed <- struct{}{}
    })

    assert.NoError(t, notifier.Start())
    defer notifier.Stop()

    // Modify file
    os.WriteFile(tmpFile, []byte("modified"), 0644)

    // Wait for notification
    select {
    case <-changed:
        // Success
    case <-time.After(5 * time.Second):
        t.Fatal("No change notification received")
    }
}
```

## Common Use Cases

### Hot-Reload Configuration

```go
type App struct {
    config   *Config
    watcher  disk.FileChangedNotifier
}

func (a *App) Start(configPath string) error {
    // Load initial config
    a.config = loadConfig(configPath)

    // Watch for changes
    a.watcher = disk.NewFileChangedNotifier(configPath)
    a.watcher.Subscribe(func(event disk.FileChangedEvent) {
        newConfig := loadConfig(event.Path)
        a.config = newConfig
        log.Println("Configuration reloaded")
    })

    return a.watcher.Start()
}
```

### Template Reloading

```go
type TemplateRenderer struct {
    templates *template.Template
    watcher   disk.FileChangedNotifier
}

func (r *TemplateRenderer) WatchTemplates(dir string) error {
    r.watcher = disk.NewFileChangedNotifier(dir)
    r.watcher.Subscribe(func(event disk.FileChangedEvent) {
        r.templates = template.Must(template.ParseGlob(dir + "/*.html"))
        log.Println("Templates reloaded")
    })

    return r.watcher.Start()
}
```

### Plugin System

```go
type PluginManager struct {
    plugins  map[string]*Plugin
    watcher  disk.FileChangedNotifier
}

func (pm *PluginManager) WatchPluginDir(dir string) error {
    pm.watcher = disk.NewFileChangedNotifier(dir)
    pm.watcher.Subscribe(func(event disk.FileChangedEvent) {
        // Reload plugin if .so file changed
        if filepath.Ext(event.Path) == ".so" {
            pm.reloadPlugin(event.Path)
        }
    })

    return pm.watcher.Start()
}
```

## Troubleshooting

### File watcher not triggering

**Cause:** File system doesn't support proper modification time updates.

**Solution:** Check file system capabilities, consider increasing polling interval:
```go
// If using custom implementation, adjust polling
ticker := time.NewTicker(2 * time.Second) // Increase from 1s
```

### Too many open files error

**Cause:** File descriptor limit reached.

**Solution:** Use fdlimit package on Unix:
```go
if runtime.GOOS != "windows" {
    fdlimit.RaiseLimit(10000)
}
```

### Path traversal vulnerabilities

**Cause:** Not validating user-provided paths.

**Solution:** Always use `IsPathSafe`:
```go
if !disk.IsPathSafe(baseDir, userPath) {
    return errors.New("invalid path")
}
```

## Migration Tips

### From fsnotify

```go
// Before (fsnotify)
watcher, _ := fsnotify.NewWatcher()
watcher.Add(filePath)
for event := range watcher.Events {
    if event.Op&fsnotify.Write == fsnotify.Write {
        handleChange()
    }
}

// After (disk module)
notifier := disk.NewFileChangedNotifier(filePath)
notifier.Subscribe(func(event disk.FileChangedEvent) {
    handleChange()
})
notifier.Start()
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md)

## License

Apache License 2.0 - see [LICENSE](../LICENSE)
