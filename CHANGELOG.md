# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2025-11-06

### 🚨 Breaking Changes

#### Removed Components
- **`errors/errorschecker` package** - Eliminated panic-based validation wrappers
  - Removed `TryPanic()` function
  - Removed `CheckNilParameter()` function
  - **Migration:** Use direct error returns and `validation.ValidateNotNil()`

- **`errors/errordefer.go`** - Eliminated panic recovery wrapper
  - Removed `ErrorDefer` interface
  - Removed `NewErrorDefer()` constructor
  - **Migration:** Use `ErrorHandler` or inline defer+recover pattern

- **`errors/validation.RequireNotNil()`** - Eliminated panic-based validation
  - **Migration:** Use `validation.ValidateNotNil()` which returns error instead of panicking

#### Changed Signatures

- **`configuration.ConfigHandler.SetRefreshTime()`**
  ```go
  // OLD v1.x
  SetRefreshTime(period Period)
  
  // NEW v2.0
  SetRefreshTime(period Period) error
  ```
  - **Reason:** Enables proper validation and error propagation
  - **Migration:** Add error handling to all calls

- **`eventsmanager.NewSubscriptions()`**
  ```go
  // OLD v1.x
  NewSubscriptions(errorDefer errors.ErrorDefer) Subscriptions
  
  // NEW v2.0
  NewSubscriptions(errorHandler errors.ErrorHandler) Subscriptions
  ```
  - **Migration:** Replace `errorDefer` with `errorHandler`

- **`fileconfig.NewFileConfigHandler()`**
  ```go
  // OLD v1.x - 7 parameters
  NewFileConfigHandler(filePath, configType, errorCatcher, errorDefer, subscriptions, publisher, filechangeNotifier)
  
  // NEW v2.0 - 6 parameters (removed errorDefer)
  NewFileConfigHandler(filePath, configType, errorCatcher, subscriptions, publisher, filechangeNotifier)
  ```

- **`server.NewListenerBuilder()`**
  ```go
  // OLD v1.x - 4 parameters
  NewListenerBuilder(configHandler, logger, errorCatcher, errorDefer)
  
  // NEW v2.0 - 3 parameters (removed errorDefer)
  NewListenerBuilder(configHandler, logger, errorCatcher)
  ```

- **`crypto.NewCipher()`**
  ```go
  // OLD v1.x
  NewCipher(key, errorCatcher, errorDefer)
  
  // NEW v2.0
  NewCipher(key, errorCatcher)
  ```

- **`orm_base.NewDataAccess()`**
  ```go
  // OLD v1.x
  NewDataAccess(errorDefer, db, modelType)
  
  // NEW v2.0
  NewDataAccess(errorHandler, db, modelType)
  ```

### ✨ Added

- **`errors/handler.go`** - New error handler with pipeline pattern
  - `ErrorHandler` interface for clean error handling
  - `NewErrorHandler(thrower ErrorThrower)` constructor
  - `Handle(err error, pipe func(error) error) error` method
  - Replaces `ErrorDefer` with cleaner, more idiomatic Go pattern

- **`errors/validation.go`** - Non-panic validation utilities
  - `ValidateNotNil(params map[string]interface{}) error`
  - Returns error instead of panicking
  - Multiple parameter validation support

- **`server_test/listener_test.go`** - Integration test with panic recovery
  - Tests listener auto-restart on port conflict
  - Validates config rollback pattern (Kubernetes/Nginx style)
  - Demonstrates proper use of defer+recover for long-running services

- **`MIGRATION-V2.md`** - Complete migration guide
  - Step-by-step migration instructions
  - Code examples (before/after)
  - Common issues and solutions
  - Migration checklist

- **IOC Module System** - New dependency injection pattern
  - `module.go` files for each package (replaces `register.go`)
  - Cleaner separation of concerns
  - Better resolver patterns

### 🔧 Changed

- **Error Handling Philosophy**
  - Moved from panic-based to error-based approach
  - Better alignment with idiomatic Go
  - Improved testability and predictability

- **Panic Recovery Strategy**
  - Removed generic panic wrapper (`ErrorDefer`)
  - Added inline defer+recover in `server.listener.startLoop()`
  - Implements industry-standard config rollback pattern
  - Only used where necessary (long-running services)

- **Constructor Signatures**
  - Simplified by removing `errorDefer` parameter
  - More consistent error handling across packages
  - Better type safety

- **README.md**
  - Updated all examples to v2.0 patterns
  - Removed references to deprecated components
  - Added v2.0 badges and migration warning
  - Updated error handling examples

### 🗑️ Removed

- All panic-based error handling utilities
- `errorschecker` package
- `errordefer.go` file
- `RequireNotNil()` function
- Unused `register.go` IOC files
- Test files using deprecated patterns

### 📝 Documentation

- **README.md** - Fully updated for v2.0
- **MIGRATION-V2.md** - Complete migration guide
- **CHANGELOG.md** - This file
- All code examples updated to v2.0 patterns

### ✅ Testing

- All tests passing (`go test ./...`)
- 33 unit tests in `errors_test`
- 1 integration test in `server_test`
- Panic recovery validated in listener test

### 🔄 Migration Path

See [MIGRATION-V2.md](./MIGRATION-V2.md) for detailed migration instructions.

**Quick migration summary:**
1. Replace `errorschecker.TryPanic()` → direct error returns
2. Replace `ErrorDefer` → `ErrorHandler` or inline defer+recover
3. Replace `RequireNotNil()` → `ValidateNotNil()`
4. Add error handling to `SetRefreshTime()` calls
5. Update constructor calls (remove `errorDefer` parameters)

### ⚠️ Important Notes

- **No backward compatibility** with v1.x
- **Go 1.20+** required
- **Thorough testing required** after migration
- **Panic recovery** should only be used for long-running services
- **Error propagation** preferred over panic/recover

---

## [1.x] - Previous Versions

See git history for changes in v1.x releases.

---

**For migration help, see:** [MIGRATION-V2.md](./MIGRATION-V2.md)  
**For usage examples, see:** [README.md](./README.md)

**Version 2.0.0 Release Date:** November 6, 2025  
**Breaking Changes:** Yes (major version bump)  
**Go Version:** 1.20+
