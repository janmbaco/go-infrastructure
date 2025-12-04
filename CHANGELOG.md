# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.1.1] - 2025-12-04

### 📚 Added

#### Comprehensive Documentation

Added complete, production-ready documentation for all infrastructure modules:

- **configuration** - Complete README covering file-based configuration with hot-reload, rollback, typed events, DI integration, and Period API
- **crypto** - Full documentation for AES-256 encryption utilities, cipher API, security best practices, and key management
- **dependencyinjection** - Extensive guide covering lifetimes, generic helpers, context support, modules, and advanced DI patterns
- **disk** - Documentation for path operations, file change notifications, FD limit management, and safe file access patterns
- **errors** - Guide to CustomError, ErrorCatcher patterns, ValidateNotNil, and centralized error handling
- **eventsmanager** - Complete event bus documentation with type-safe events, subscriptions, publishers, and DI integration
- **logs** - Full logging documentation covering log levels, file rotation, ErrorLogger shortcuts, and DI integration
- **persistence/dataaccess** - Comprehensive GORM wrapper guide with generic functions, CRUD operations, and repository patterns
- **server** - Production-ready HTTP/HTTPS server documentation with SPA support, graceful shutdown, and configuration integration

All documentation follows consistent style:
- Practical code examples
- Security best practices
- Integration patterns
- Troubleshooting guides
- Testing strategies

### ✨ Enhanced

#### Dependency Injection - Full Context Support

The `dependencyinjection` package now provides **complete context.Context support** throughout the dependency resolution chain:

- **Context propagation**: Context flows through the entire dependency graph during resolution
- **Automatic detection**: Providers can optionally accept `context.Context` as their first parameter
- **Cancellation support**: Resolution respects `ctx.Done()` and stops immediately on cancellation
- **Timeout support**: Honors context deadlines during dependency creation
- **Value propagation**: Pass request-scoped data (request IDs, tenant info, tracing spans) through context
- **Nested propagation**: Context automatically propagates to all nested dependencies

**What changed:**
- `*Ctx` methods (`TypeCtx`, `TenantCtx`, `ResolveCtx`, etc.) now **actually use** the provided context instead of ignoring it
- `DependencyObject.Create()` now accepts `context.Context` as first parameter
- Provider functions can optionally accept `context.Context` as first parameter (automatically detected via reflection)
- Context cancellation is checked at each resolution step
- Methods without `Ctx` suffix use `context.Background()` internally (maintaining backward compatibility)

**Backward compatibility:**
- ✅ **Non-breaking change**: All existing code continues to work
- ✅ Methods without context (`Type()`, `Resolve()`, etc.) still work exactly as before
- ✅ Existing providers without context parameter continue to work
- ✅ The `*Ctx` methods existed before but now provide real functionality

**Example usage:**

```go
// Provider with context - automatically detected
di.RegisterType[*Service](r, func(ctx context.Context, logger logs.Logger) *Service {
    requestID := ctx.Value("requestID").(string)
    return &Service{
        Logger:    logger,
        RequestID: requestID,
    }
})

// Resolve with context containing request-scoped data
ctx := context.WithValue(r.Context(), "requestID", "req-12345")
service := di.ResolveCtx[*Service](ctx, resolver)

// Context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
service := di.ResolveCtx[*Service](ctx, resolver) // Respects timeout
```

**New test coverage:**
- 8 comprehensive tests for context support in `dependencyinjection/context_test.go`
- Example code demonstrating context usage in `dependencyinjection/examples/context_example.go`

**Documentation:**
- Full context support section added to `dependencyinjection/README.md`
- Best practices and usage patterns included

### 🧹 Removed

#### Cleanup of Empty Files and Artifacts

- Removed empty `eventsmanager/subscriptions_test.go`
- Removed empty `eventsmanager/examples/subscription_remove_example.go`
- Removed empty `eventsmanager/examples/` directory
- Cleaned up temporary test artifacts (`integration_test.test.exe`, `coverage.out`)

### 🔧 Fixed

#### Documentation Consistency

- Standardized `server/README.md` to match project documentation style
  - Removed decorative emojis from headers
  - Restructured with consistent section ordering
  - Changed "Installation" to "Install" for consistency
  - Added proper `---` separators between sections
  - Applied technical, professional tone throughout

### 🛠️ Internal

#### Testing and Quality

- Added 8 new tests for context support in dependency injection
- Added context usage examples
- Added database connection utility for integration tests (`persistence/integration_test/test-db-connection.go`)
- All existing tests continue to pass (15+ modules verified)

---

## [2.1.0] - 2025-11-29

### 🔄 Changed

#### Persistence Layer Refactoring

- **Breaking Change**: Restructured persistence module architecture
  - Moved from `persistence/orm_base` to `persistence/dataaccess` structure
  - Renamed core package for better clarity and organization
  - All database access now through `dataaccess` package

#### Module Structure Improvements

- Reorganized `persistence/` package structure:
  - `dataaccess/` - Core data access layer with GORM wrapper
  - `dialectors/` - Database dialector implementations (moved from `orm_base/dialectors`)
  - Root level now contains shared types (`DatabaseInfo`, `DatabaseProvider`, etc.)
  - Improved IoC module organization with typed resolvers

### 📚 Added

#### Comprehensive Test Coverage

Added extensive test suites across all modules:

- **configuration**: 173 tests for events, 456 for file config handler, 73 for Period
- **crypto**: 90 tests for cipher operations
- **dependencyinjection**: 277 builder tests, 159 dependencies tests, 187 register tests, 99 resolver tests
- **disk**: 176 tests for path operations
- **errors**: 55 CustomError tests, 243 ErrorCatcher tests
- **eventsmanager**: 85 tests for event manager
- **logs**: 227 tests for logger functionality
- **persistence**: 124 general tests, 501 dataaccess tests
- **server**: 128 tests for single-page app functionality

Total: **3000+ new tests** ensuring production-ready quality

#### Integration Tests Infrastructure

- Added Docker-based integration tests for database operations
  - PostgreSQL integration tests
  - MySQL integration tests
  - SQL Server integration tests
- PowerShell script (`run-integration-tests.ps1`) to orchestrate Docker containers
- `docker-compose.test.yml` for test database setup
- Database connection utility (`test-db-connection.go`) for verification
- Comprehensive integration test documentation (`INTEGRATION_TESTS.md`)

#### New Persistence Features

- **Generic data access helpers** in `dataaccess/dataaccess_generics.go`
  - Type-safe `InsertRow[T]`, `SelectRows[T]`, `UpdateRow[T]`, `DeleteRows[T]`
- **GORM interface abstraction** for better testability
- **SQLite dialector** support added
- **SQL Server dialector** implementation completed
- **Typed DataAccess IoC resolver** for dependency injection
- **Comprehensive persistence README** with usage examples

### 🔧 Updated

#### Dependencies and Tooling

- **Go version**: Updated to **1.24.0** (minimum required)
- **golang.org/x/crypto**: Bumped from 0.43.0 to 0.45.0
- **google.golang.org/grpc**: Updated from 1.76.0 to 1.77.0
- **actions/checkout**: Updated from v5 to v6 in CI workflow
- **codecov/codecov-action**: Updated from v4 to v5

#### CI/CD Improvements

- Updated `.github/workflows/ci.yml` to support Go 1.24+
- Fixed codecov action parameters for better coverage reporting
- Added integration test support in CI pipeline
- Enhanced Makefile with test and coverage targets

#### Documentation

- Expanded main `README.md` with better examples and module overview
- Added `persistence/README.md` with architecture guide
- Improved documentation for test utilities
- Added coverage reporting configuration

### 🛠️ Internal

#### Code Quality

- Improved error handling across persistence layer
- Enhanced type safety with generic helpers
- Better separation of concerns in IoC modules
- Standardized test patterns across all modules
- Added coverage tracking (43% baseline established)


**See also:**
- New comprehensive tests in `dependencyinjection/context_test.go`
- Extended documentation in `dependencyinjection/README.md` (Context Support section)

---

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
