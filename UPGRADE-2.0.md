# Upgrade Guide: v1.x → v2.0.0

> **⚠️ MAJOR VERSION UPGRADE**  
> Version 2.0 is a major breaking release that removes panic-based patterns and introduces cleaner, error-based handling. This guide provides step-by-step migration instructions.

---

## Overview

**go-infrastructure v2.0** removes deprecated panic recovery patterns and introduces cleaner error handling. This is a **major breaking release** requiring code updates.

### Key Changes Summary

| Area | v1.x | v2.0 | Impact |
|------|------|------|--------|
| **Module Path** | `github.com/janmbaco/go-infrastructure` | `github.com/janmbaco/go-infrastructure/v2` | 🔴 Breaking |
| **Error Validation** | `RequireNotNil` (panics) | `ValidateNotNil` (returns error) | 🔴 Breaking |
| **Error Recovery** | `ErrorDefer` | `ErrorHandler` | 🔴 Breaking |
| **Dependency Injection** | Manual container | Module-based `Builder` | 🟡 Improved |
| **Event System** | Scattered handling | `EventManager` centralized | 🟢 Addition |
| **ConfigHandler** | `SetRefreshTime(period)` | `SetRefreshTime(period) error` | 🔴 Breaking |

---

## Breaking Changes Overview

---

## Prerequisites

Before upgrading, ensure:

1. **Backup your code**: Commit all changes to version control
2. **Review this guide completely** before starting
3. **Test in a non-production environment** first
4. **Go 1.23+** installed
5. **Read about module path changes** (critical for Go modules)

---

## Module Path Change (Go Modules)

**CRITICAL**: Major version v2 requires `/v2` suffix in module path per Go modules specification.

### v1.x Module & Imports

```go
// go.mod
module myapp

require (
    github.com/janmbaco/go-infrastructure v1.2.5
)
```

```go
// your code
import (
    "github.com/janmbaco/go-infrastructure/logs"
    "github.com/janmbaco/go-infrastructure/errors"
)
```

### v2.0 Module & Imports

```go
// go.mod
module myapp

require (
    github.com/janmbaco/go-infrastructure/v2 v2.0.0
)
```

```go
// your code
import (
    "github.com/janmbaco/go-infrastructure/v2/logs"
    "github.com/janmbaco/go-infrastructure/v2/errors"
)
```

**Migration Action**:
```bash
# Update go.mod
go get github.com/janmbaco/go-infrastructure/v2

# Update all imports
find . -name '*.go' -exec sed -i 's|github.com/janmbaco/go-infrastructure/|github.com/janmbaco/go-infrastructure/v2/|g' {} +

# Tidy dependencies
go mod tidy
```

---

## 🚨 Removed Components

### 1. `errors/errorschecker` Package (DELETED)

**What was removed:**
```go
// ❌ v1.x - NO LONGER AVAILABLE
import "github.com/janmbaco/go-infrastructure/errors/errorschecker"

errorschecker.TryPanic(func() { /* code */ })
errorschecker.CheckNilParameter(param, "paramName")
```

**Migration:**
```go
// ✅ v2.0 - Use direct error handling
import "github.com/janmbaco/go-infrastructure/v2/errors/validation"

func MyFunction() error {
    // Instead of TryPanic wrapper, return errors directly
    if err := someOperation(); err != nil {
        return fmt.Errorf("operation failed: %w", err)
    }
    return nil
}

// Instead of CheckNilParameter, use ValidateNotNil
func MyFunction(param interface{}) error {
    if err := validation.ValidateNotNil(map[string]interface{}{"param": param}); err != nil {
        return err
    }
    // ... rest of logic
    return nil
}
```

---

### 2. `errors/errordefer.go` (DELETED)

**What was removed:**
```go
// ❌ v1.x - NO LONGER AVAILABLE
errorDefer := errors.NewErrorDefer(errorThrower)
defer errorDefer.TryThrowError(pipeError)
```

**Migration Strategy:**

#### Option A: Use `ErrorHandler` (Recommended)
```go
// ✅ v2.0 - Recommended pattern
errorManager := errors.NewErrorManager()
errorThrower := errors.NewErrorThrower(errorManager)
errorHandler := errors.NewErrorHandler(errorThrower)

err := errorHandler.Handle(someError, func(e error) error {
    return fmt.Errorf("wrapped: %w", e)
})
```

#### Option B: Inline defer+recover (When panic recovery is required)
```go
// ✅ v2.0 - Inline panic recovery
func riskyOperation() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
            // Optional: execute callbacks via ErrorThrower
            if thrower != nil {
                err = thrower.Throw(err)
            }
        }
    }()
    
    // ... risky code that might panic
    return nil
}
```

**When to use each option:**
- **ErrorHandler**: Normal error flow without panic
- **Inline defer+recover**: When you need panic recovery (e.g., long-running services)

---

### 3. `errors/validation.RequireNotNil` (DELETED)

**What was removed:**
```go
// ❌ v1.x - NO LONGER AVAILABLE
import "github.com/janmbaco/go-infrastructure/errors/validation"

validation.RequireNotNil(param, "param")
```

**Migration:**
```go
// ✅ v2.0 - Use ValidateNotNil (returns error instead of panic)
import "github.com/janmbaco/go-infrastructure/v2/errors/validation"

if err := validation.ValidateNotNil(map[string]interface{}{"param": param}); err != nil {
    return err
}

// Or use standard error handling
if param == nil {
    return fmt.Errorf("param cannot be nil")
}
```

---

## 🔧 Modified Components

### 4. `configuration.ConfigHandler.SetRefreshTime()` (Breaking Change)

**v1.x Signature:**
```go
// ❌ v1.x
SetRefreshTime(period Period)  // void return
```

**v2.0 Signature:**
```go
// ✅ v2.0
SetRefreshTime(period Period) error  // returns error
```

**Migration:**
```go
// v1.x
configHandler.SetRefreshTime(period)

// v2.0
if err := configHandler.SetRefreshTime(period); err != nil {
    return fmt.Errorf("failed to set refresh time: %w", err)
}
```

**Reason:** Allows proper validation and error propagation when `period` is nil or invalid.

---

## 📦 Component Updates

### Updated Constructor Signatures

#### `eventsmanager.NewSubscriptions`
```go
// v1.x
NewSubscriptions(errorDefer errors.ErrorDefer) Subscriptions

// ✅ v2.0
NewSubscriptions(errorHandler errors.ErrorHandler) Subscriptions
```

#### `fileconfig.NewFileConfigHandler`
```go
// v1.x
NewFileConfigHandler(
    filePath string,
    configType interface{},
    errorCatcher errors.ErrorCatcher,
    errorDefer errors.ErrorDefer,
    subscriptions eventsmanager.Subscriptions,
    publisher eventsmanager.Publisher,
    filechangeNotifier disk.FileChangedNotifier,
) (ConfigHandler, error)

// ✅ v2.0 - Removed errorDefer parameter
NewFileConfigHandler(
    filePath string,
    configType interface{},
    errorCatcher errors.ErrorCatcher,
    subscriptions eventsmanager.Subscriptions,
    publisher eventsmanager.Publisher,
    filechangeNotifier disk.FileChangedNotifier,
) (ConfigHandler, error)
```

#### `server.NewListenerBuilder`
```go
// v1.x
NewListenerBuilder(
    configHandler configuration.ConfigHandler,
    logger logs.Logger,
    errorCatcher errors.ErrorCatcher,
    errorDefer errors.ErrorDefer,
) ListenerBuilder

// ✅ v2.0 - Removed errorDefer parameter
NewListenerBuilder(
    configHandler configuration.ConfigHandler,
    logger logs.Logger,
    errorCatcher errors.ErrorCatcher,
) ListenerBuilder
```

---

## Step-by-Step Migration Guide

### Step 1: Backup Everything

```bash
# Commit all changes
git add .
git commit -m "Pre-migration backup before upgrading to go-infrastructure v2.0"
git tag v1-backup
```

---

### Step 2: Update Module Path

```bash
# Update to v2
go get github.com/janmbaco/go-infrastructure/v2@v2.0.0

# Remove old version
go mod edit -droprequire github.com/janmbaco/go-infrastructure

# Clean up
go mod tidy
```

---

### Step 3: Update All Import Paths

```bash
# Find all Go files importing old paths
grep -r "github.com/janmbaco/go-infrastructure" . --include="*.go" --exclude-dir=vendor

# Replace v1 imports with v2
find . -name '*.go' -not -path "*/vendor/*" -exec sed -i 's|github.com/janmbaco/go-infrastructure/|github.com/janmbaco/go-infrastructure/v2/|g' {} +

# Verify changes
git diff
```

---

### Step 4: Update Deprecated APIs

Follow the migration patterns in this guide for:
- [ ] `errorschecker` → `validation.ValidateNotNil`
- [ ] `ErrorDefer` → `ErrorHandler`
- [ ] `RequireNotNil` → `ValidateNotNil`
- [ ] `ConfigHandler.SetRefreshTime` → Add error handling
- [ ] Constructor signatures (remove `errorDefer` parameters)

---

### Step 5: Build and Test

```bash
# Build to catch compilation errors
go build ./...

# Run tests
go test ./...

# Check for runtime issues
go run . # or your main entry point
```

---

### Step 6: Fix Compilation Errors

Common errors and fixes:

**Error**: `cannot find package "github.com/janmbaco/go-infrastructure/errors/errorschecker"`
```bash
# Solution: Update import
import "github.com/janmbaco/go-infrastructure/v2/errors/validation"
```

**Error**: `too many arguments in call to NewFileConfigHandler`
```bash
# Solution: Remove errorDefer parameter from constructor
```

**Error**: `SetRefreshTime used as value`
```bash
# Solution: Add error handling
if err := configHandler.SetRefreshTime(period); err != nil {
    return err
}
```

---

### Step 7: Verify Migration

```bash
# Build succeeds
go build ./...

# Tests pass
go test ./... -v

# Application starts without errors
go run . # Check logs for startup errors
```

---

## 🎯 Migration Checklist

### Step 1: Remove Deprecated Imports
```bash
# Find all usages of deprecated packages
grep -r "errors/errorschecker" .
grep -r "errordefer" .
grep -r "RequireNotNil" .
```

### Step 2: Replace `errorschecker` Usages
- [ ] Replace `TryPanic()` with direct error returns
- [ ] Replace `CheckNilParameter()` with `ValidateNotNil()` or manual checks

### Step 3: Replace `ErrorDefer` Pattern
- [ ] Identify all `NewErrorDefer()` calls
- [ ] Replace with `ErrorHandler` or inline defer+recover
- [ ] Update constructor calls (remove `errorDefer` parameters)

### Step 4: Update `SetRefreshTime` Calls
- [ ] Add error handling to all `SetRefreshTime()` calls

### Step 5: Update Constructor Signatures
- [ ] `NewSubscriptions` - pass `ErrorHandler` instead of `ErrorDefer`
- [ ] `NewFileConfigHandler` - remove `errorDefer` parameter
- [ ] `NewListenerBuilder` - remove `errorDefer` parameter

### Step 6: Test
```bash
go build ./...
go test ./...
```

---

## 📋 Complete Migration Example

### v1.x Code
```go
package main

import (
    "github.com/janmbaco/go-infrastructure/configuration/fileconfig"
    "github.com/janmbaco/go-infrastructure/errors"
    "github.com/janmbaco/go-infrastructure/errors/errorschecker"
    "github.com/janmbaco/go-infrastructure/eventsmanager"
    "github.com/janmbaco/go-infrastructure/logs"
)

func setupListener() {
    logger := logs.NewLogger()
    errorCatcher := errors.NewErrorCatcher(logger)
    errorThrower := errors.NewErrorThrower(nil)
    errorDefer := errors.NewErrorDefer(errorThrower)
    
    subscriptions := eventsmanager.NewSubscriptions(errorDefer)
    
    errorschecker.CheckNilParameter(logger, "logger")
    
    configHandler.SetRefreshTime(period)
}
```

### v2.0 Code
```go
package main

import (
    "github.com/janmbaco/go-infrastructure/v2/configuration/fileconfig"
    "github.com/janmbaco/go-infrastructure/v2/errors"
    "github.com/janmbaco/go-infrastructure/v2/errors/validation"
    "github.com/janmbaco/go-infrastructure/v2/eventsmanager"
    "github.com/janmbaco/go-infrastructure/v2/logs"
)

func setupListener() error {
    logger := logs.NewLogger()
    errorCatcher := errors.NewErrorCatcher(logger)
    errorManager := errors.NewErrorManager()
    errorThrower := errors.NewErrorThrower(errorManager)
    errorHandler := errors.NewErrorHandler(errorThrower)
    
    subscriptions := eventsmanager.NewSubscriptions(errorHandler)
    
    if err := validation.ValidateNotNil(map[string]interface{}{"logger": logger}); err != nil {
        return err
    }
    
    if err := configHandler.SetRefreshTime(period); err != nil {
        return fmt.Errorf("failed to set refresh time: %w", err)
    }
    
    return nil
}
```

---

## 🔄 Error Handling Pattern Changes

### Old Pattern (v1.x)
```go
// Panic-based with recovery
defer errorDefer.TryThrowError(func(err error) error {
    return fmt.Errorf("wrapped: %w", err)
})

errors.RequireNotNil(value, "value")  // panics if nil
errorschecker.TryPanic(func() {       // catches panics
    riskyOperation()
})
```

### New Pattern (v2.0)
```go
// Error-based (no panic)
if err := errorHandler.Handle(err, func(e error) error {
    return fmt.Errorf("wrapped: %w", e)
}); err != nil {
    return err
}

if err := validation.ValidateNotNil(value, "value"); err != nil {
    return err
}

if err := riskyOperation(); err != nil {
    return err
}
```

---

## 🆘 Need Help?

### Common Migration Issues

**Issue 1: Missing ErrorHandler**
```
Error: cannot use errorDefer as ErrorHandler
```
**Solution:** Replace `errorDefer` with `errorHandler`:
```go
errorManager := errors.NewErrorManager()
errorThrower := errors.NewErrorThrower(errorManager)
errorHandler := errors.NewErrorHandler(errorThrower)
```

**Issue 2: Too Many Arguments**
```
Error: too many arguments in call to NewFileConfigHandler
```
**Solution:** Remove `errorDefer` parameter from constructor call.

**Issue 3: SetRefreshTime Type Mismatch**
```
Error: SetRefreshTime used as value
```
**Solution:** Add error handling:
```go
if err := configHandler.SetRefreshTime(period); err != nil {
    return err
}
```

---

## 📚 Additional Resources

- **README.md** - Updated usage examples for v2.0
- **server_test/listener_test.go** - Integration test example with new patterns
- **errors/handler.go** - ErrorHandler implementation

---

## ⚠️ Important Notes

1. **No Backward Compatibility**: v2.0 is a breaking change release
2. **Semantic Versioning**: Update imports to `v2` module path (if using modules)
3. **Testing Required**: Thoroughly test your code after migration
4. **Panic Recovery**: Only use inline defer+recover when absolutely necessary (e.g., long-running servers)
5. **Error Propagation**: Prefer returning errors over panic/recover

---

**Migration completed?** Run full test suite:
```bash
go test ./... -v
go build ./...
```

**Version:** 2.0.0  
**Migration Date:** November 6, 2025  
**Go Version Required:** 1.20+
