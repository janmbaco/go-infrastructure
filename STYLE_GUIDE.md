# Go Style Guide for go-infrastructure

This style guide defines the coding standards and conventions for the `go-infrastructure` project. Following these guidelines ensures consistency, readability, and maintainability across the codebase.

## Table of Contents

- [General Principles](#general-principles)
- [Code Formatting](#code-formatting)
- [Naming Conventions](#naming-conventions)
- [Package Design](#package-design)
- [Error Handling](#error-handling)
- [Interfaces](#interfaces)
- [Context Usage](#context-usage)
- [Concurrency](#concurrency)
- [Testing](#testing)
- [Documentation](#documentation)
- [Project-Specific Conventions](#project-specific-conventions)

## General Principles

### Simplicity First

```go
// ✅ Good - simple and clear
func GetUser(id string) (*User, error) {
    return db.FindByID(id)
}

// ❌ Bad - unnecessary abstraction
func GetUser(id string) (*User, error) {
    strategy := NewUserRetrievalStrategy()
    context := NewRetrievalContext(id)
    return strategy.Execute(context)
}
```

### Go Idioms

Follow established Go idioms and patterns. When in doubt, refer to:
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Proverbs](https://go-proverbs.github.io/)

### Readability Over Cleverness

```go
// ✅ Good - clear intent
if err != nil {
    return fmt.Errorf("failed to connect: %w", err)
}

// ❌ Bad - clever but obscure
if err != nil {
    return errors.Wrap(err, "failed to connect").If(prod).Else(err)
}
```

## Code Formatting

### Use gofmt

Always run `gofmt` before committing. Configure your editor to run it on save.

```bash
# Format all files
go fmt ./...

# Format specific package
go fmt ./dependencyinjection/...
```

### Line Length

- Keep lines under **120 characters** when possible
- Break long lines at logical points
- Use multiple lines for long function signatures

```go
// ✅ Good
func RegisterServiceWithDependencies(
    register Register,
    service ServiceProvider,
    dependencies []Dependency,
    options ...Option,
) error {
    // implementation
}

// ❌ Bad - too long
func RegisterServiceWithDependencies(register Register, service ServiceProvider, dependencies []Dependency, options ...Option) error {
```

### Import Organization

Group imports in three sections:

```go
import (
    // 1. Standard library
    "context"
    "fmt"
    "time"
    
    // 2. Third-party packages
    "github.com/stretchr/testify/assert"
    "gorm.io/gorm"
    
    // 3. Local packages
    "github.com/janmbaco/go-infrastructure/v2/errors"
    "github.com/janmbaco/go-infrastructure/v2/logs"
)
```

### Struct Initialization

Use field names for clarity:

```go
// ✅ Good
user := User{
    ID:        "123",
    Name:      "John",
    Email:     "john@example.com",
    CreatedAt: time.Now(),
}

// ❌ Bad - positional arguments are fragile
user := User{"123", "John", "john@example.com", time.Now()}
```

## Naming Conventions

### General Rules

- Use **camelCase** for unexported names: `localVariable`, `privateMethod`
- Use **PascalCase** for exported names: `PublicFunction`, `ExportedType`
- Use **ALL_CAPS** for constants: `MaxRetries`, `DefaultTimeout`
- Avoid abbreviations unless widely known: `DB`, `HTTP`, `URL` are okay

### Package Names

```go
// ✅ Good - short, lowercase, no underscores
package dependencyinjection
package errors
package persistence

// ❌ Bad
package dependency_injection
package Errors
package persistenceLayer
```

### Interface Names

Single-method interfaces should end in `-er`:

```go
// ✅ Good
type Logger interface {
    Log(message string)
}

type Resolver interface {
    Resolve(key string) interface{}
}

// Multi-method interfaces use descriptive names
type ConfigHandler interface {
    GetConfig() Config
    SetConfig(config Config) error
    Refresh() error
}
```

### Variable Names

```go
// ✅ Good - short names in small scopes
for i := 0; i < len(items); i++ {
    // i is clear in this context
}

// ✅ Good - descriptive names in larger scopes
func ProcessUserRequest(requestID string) error {
    userRepository := NewUserRepository()
    // ...
}

// ❌ Bad - too verbose
for indexOfCurrentItem := 0; indexOfCurrentItem < len(items); indexOfCurrentItem++ {
    // ...
}
```

### Receiver Names

Use short, consistent receiver names:

```go
// ✅ Good - consistent, short
func (c *Container) Register() Register {
    return c.register
}

func (c *Container) Resolver() Resolver {
    return c.resolver
}

// ❌ Bad - inconsistent or too verbose
func (container *Container) Register() Register {
    return container.register
}

func (this *Container) Resolver() Resolver {
    return this.resolver
}
```

### Test Function Names

Follow the pattern: `Test<Function>_When<Condition>_Then<Result>`

```go
// ✅ Good
func TestResolver_TypeCtx_WhenContextCanceled_ThenPanics(t *testing.T) {
    // ...
}

func TestContainer_Register_WhenValidInput_ThenRegistersSuccessfully(t *testing.T) {
    // ...
}

// ❌ Bad
func TestResolver(t *testing.T) {
    // unclear what is being tested
}
```

## Package Design

### Package Organization

```
mypackage/
├── mypackage.go          # Main types and public API
├── mypackage_test.go     # Tests for public API
├── internal.go           # Internal implementation
├── internal_test.go      # Tests for internal code
├── errors.go             # Package-specific errors
├── ioc/                  # IoC module (if applicable)
│   └── module.go
└── README.md             # Package documentation
```

### Avoid Circular Dependencies

```go
// ✅ Good - clear dependency direction
package persistence
import "github.com/janmbaco/go-infrastructure/v2/errors"

// ❌ Bad - circular dependency
package errors
import "github.com/janmbaco/go-infrastructure/v2/persistence"
```

### Internal Packages

Use `internal/` for code that should not be imported by external packages:

```
persistence/
├── dataaccess.go
└── internal/
    └── query_builder.go  # Cannot be imported outside persistence
```

## Error Handling

### Always Handle Errors

```go
// ✅ Good
result, err := DoSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// ❌ Bad - ignoring errors
result, _ := DoSomething()
```

### Error Wrapping

Use `%w` to wrap errors for better error chains:

```go
// ✅ Good
if err := validateInput(data); err != nil {
    return fmt.Errorf("validation failed: %w", err)
}

// ❌ Bad - loses error context
if err := validateInput(data); err != nil {
    return errors.New("validation failed")
}
```

### Error Types

Define custom error types for package-specific errors:

```go
// ✅ Good
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

// Usage
if !isValid(input) {
    return &ValidationError{Field: "email", Message: "invalid format"}
}
```

### Panic vs Error

```go
// ✅ Good - return errors for expected failures
func OpenFile(path string) (*File, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    return file, nil
}

// ✅ Acceptable - panic for programmer errors
func MustCompile(pattern string) *regexp.Regexp {
    re, err := regexp.Compile(pattern)
    if err != nil {
        panic(fmt.Sprintf("invalid regex pattern: %v", err))
    }
    return re
}

// ❌ Bad - panic for expected failures
func OpenFile(path string) *File {
    file, err := os.Open(path)
    if err != nil {
        panic(err) // Don't do this!
    }
    return file
}
```

## Interfaces

### Accept Interfaces, Return Structs

```go
// ✅ Good
func NewService(logger Logger, db Database) *Service {
    return &Service{logger: logger, db: db}
}

// ❌ Bad - returning interface
func NewService(logger Logger, db Database) ServiceInterface {
    return &Service{logger: logger, db: db}
}
```

### Small Interfaces

Prefer small, focused interfaces:

```go
// ✅ Good - small, focused
type Logger interface {
    Log(message string)
}

type ErrorLogger interface {
    LogError(err error)
}

// ❌ Bad - too many responsibilities
type Logger interface {
    Log(message string)
    LogError(err error)
    LogWithLevel(level Level, message string)
    SetOutput(w io.Writer)
    GetLevel() Level
    SetLevel(level Level)
}
```

### Interface Segregation

Define interfaces where they are used:

```go
// ✅ Good - defined in consumer package
package service

type Repository interface {
    GetUser(id string) (*User, error)
}

type UserService struct {
    repo Repository
}

// ❌ Bad - defined in provider package
package repository

type Repository interface {
    GetUser(id string) (*User, error)
    GetPost(id string) (*Post, error)
    // ... many methods
}
```

## Context Usage

### Context as First Parameter

```go
// ✅ Good
func DoSomething(ctx context.Context, arg string) error {
    // ...
}

// ❌ Bad
func DoSomething(arg string, ctx context.Context) error {
    // ...
}
```

### Don't Store Context

```go
// ✅ Good
func (s *Service) Process(ctx context.Context) error {
    // use ctx within this function
}

// ❌ Bad - storing context in struct
type Service struct {
    ctx context.Context // Don't do this!
}
```

### Context Values

Use typed keys for context values:

```go
// ✅ Good
type contextKey string

const (
    requestIDKey contextKey = "requestID"
    userIDKey    contextKey = "userID"
)

func WithRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, requestIDKey, id)
}

func GetRequestID(ctx context.Context) string {
    if id, ok := ctx.Value(requestIDKey).(string); ok {
        return id
    }
    return ""
}

// ❌ Bad - string keys are public
ctx = context.WithValue(ctx, "requestID", id)
```

## Concurrency

### Use Channels for Ownership Transfer

```go
// ✅ Good - clear ownership
func worker(jobs <-chan Job, results chan<- Result) {
    for job := range jobs {
        results <- process(job)
    }
}
```

### Use Mutexes for Shared State

```go
// ✅ Good
type Cache struct {
    mu    sync.RWMutex
    items map[string]interface{}
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    item, exists := c.items[key]
    return item, exists
}
```

### Don't Start Goroutines in Libraries

```go
// ✅ Good - let caller control goroutines
func ProcessItems(items []Item) []Result {
    results := make([]Result, len(items))
    for i, item := range items {
        results[i] = process(item)
    }
    return results
}

// Caller decides concurrency
go ProcessItems(items)

// ❌ Bad - library starts goroutines
func ProcessItems(items []Item) []Result {
    go func() {
        // processing
    }()
}
```

## Testing

### Table-Driven Tests

```go
// ✅ Good
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"with zero", 5, 0, 5},
        {"negative numbers", -2, -3, -5},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Use testify/assert

```go
// ✅ Good
assert.Equal(t, expected, actual)
assert.NoError(t, err)
assert.True(t, condition)

// ❌ Bad
if actual != expected {
    t.Errorf("expected %v, got %v", expected, actual)
}
```

### Test Both Success and Failure

```go
func TestGetUser_WhenUserExists_ThenReturnsUser(t *testing.T) {
    // test success case
}

func TestGetUser_WhenUserNotFound_ThenReturnsError(t *testing.T) {
    // test failure case
}
```

## Documentation

### Package Comments

```go
// Package dependencyinjection provides a type-safe dependency injection
// container for Go applications with support for lifetimes, generics,
// and context-aware resolution.
//
// Basic usage:
//
//     container := di.NewBuilder().
//         Register(func(r di.Register) {
//             di.RegisterSingleton[*Service](r, NewService)
//         }).
//         MustBuild()
//
//     service := di.Resolve[*Service](container.Resolver())
package dependencyinjection
```

### Function Comments

```go
// NewContainer creates a new dependency injection container.
// The container is initialized with default registrations for
// the Container, Register, and Resolver interfaces.
//
// Example:
//
//     container := NewContainer()
//     container.Register().AsSingleton(new(*Service), NewService, nil)
//     service := container.Resolver().Type(new(*Service), nil)
func NewContainer() Container {
    // implementation
}
```

### Exported Names Must Be Documented

```go
// ✅ Good
// Logger defines the interface for logging operations.
type Logger interface {
    Log(message string)
}

// ❌ Bad - no documentation
type Logger interface {
    Log(message string)
}
```

## Project-Specific Conventions

### Dependency Injection

```go
// ✅ Good - use generics for type safety
di.RegisterSingleton[*UserService](r, func() *UserService {
    return NewUserService()
})

service := di.Resolve[*UserService](resolver)

// ❌ Bad - type assertion required
r.AsSingleton(new(*UserService), func() interface{} {
    return NewUserService()
}, nil)

service := resolver.Type(new(*UserService), nil).(*UserService)
```

### Error Handling Pattern

```go
// ✅ Good - use ErrorCatcher for consistent error handling
func (s *Service) Process() error {
    return s.errorCatcher.TryCatchError(
        func() error {
            return s.doProcess()
        },
        func(err error) error {
            s.logger.Error("processing failed", err)
            return err
        },
    )
}
```

### IoC Module Pattern

```go
// ✅ Good - consistent module structure
package ioc

type Module struct{}

func NewModule() di.Module {
    return &Module{}
}

func (m *Module) RegisterServices(register di.Register) error {
    // registrations
    return nil
}
```

### Context-Aware Providers

```go
// ✅ Good - context as first parameter
r.AsType(new(*Service), func(ctx context.Context, logger Logger) *Service {
    return &Service{
        logger:    logger,
        requestID: ctx.Value(requestIDKey).(string),
    }
}, nil)
```

## Tools and Automation

### Recommended Tools

- **gofmt** - Code formatting
- **go vet** - Static analysis
- **golangci-lint** - Comprehensive linting
- **staticcheck** - Advanced static analysis

### Pre-commit Checks

```bash
#!/bin/bash
# .git/hooks/pre-commit

go fmt ./...
go vet ./...
golangci-lint run
go test ./...
```

## References

- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Google Go Style Guide](https://google.github.io/styleguide/go/)

---

When in doubt, prioritize **readability**, **simplicity**, and **consistency** with the existing codebase.
