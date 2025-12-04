# Contributing to go-infrastructure

Thank you for your interest in contributing to `go-infrastructure`! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Pull Request Process](#pull-request-process)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Documentation](#documentation)
- [Community](#community)

## Code of Conduct

This project adheres to the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/go-infrastructure.git
   cd go-infrastructure
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/janmbaco/go-infrastructure.git
   ```
4. **Create a branch** for your work:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional, for using Makefile commands)

### Install Dependencies

```bash
go mod download
```

### Run Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run tests for a specific package
go test ./dependencyinjection/...

# Run tests with verbose output
go test -v ./...
```

### Build

```bash
# Build all packages
go build ./...

# Build specific command
go build ./cmd/singlepageapp
```

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title and description**
- **Go version** (`go version`)
- **Operating system** and version
- **Steps to reproduce** the issue
- **Expected behavior** vs **actual behavior**
- **Code samples** or test cases that demonstrate the issue
- **Error messages** or stack traces

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

- **Clear title and description** of the enhancement
- **Use case** - why is this enhancement useful?
- **Proposed solution** or implementation approach
- **Alternatives considered**
- **Examples** of how the feature would be used

### Your First Code Contribution

Unsure where to begin? Look for issues labeled:

- `good first issue` - simpler issues for newcomers
- `help wanted` - issues where we'd appreciate community help
- `documentation` - documentation improvements

## Coding Standards

### Go Style Guide

Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines and the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

### Key Principles

1. **Keep it simple** - Prefer clarity over cleverness
2. **Follow Go idioms** - Write idiomatic Go code
3. **Error handling** - Always handle errors appropriately
4. **Documentation** - Document all exported functions, types, and packages
5. **Naming conventions**:
   - Use `camelCase` for unexported names
   - Use `PascalCase` for exported names
   - Use meaningful, descriptive names
   - Avoid abbreviations unless widely known

### Code Organization

```go
// Package documentation comes first
package mypackage

import (
    // Standard library imports
    "context"
    "fmt"
    
    // Third-party imports
    "github.com/some/package"
    
    // Local imports
    "github.com/janmbaco/go-infrastructure/v2/errors"
)

// Constants
const (
    DefaultTimeout = 30 * time.Second
)

// Types
type MyService struct {
    // fields
}

// Constructor
func NewMyService() *MyService {
    return &MyService{}
}

// Methods
func (s *MyService) DoSomething() error {
    return nil
}
```

### Code Formatting

- Run `go fmt` on your code before committing
- Run `go vet` to catch common mistakes
- Consider using `golangci-lint` for additional checks:
  ```bash
  golangci-lint run
  ```

## Testing Guidelines

### Writing Tests

1. **Test file naming**: `*_test.go`
2. **Test function naming**: `TestFunctionName_WhenCondition_ThenExpectedBehavior`
3. **Use table-driven tests** when testing multiple scenarios
4. **Use testify/assert** for assertions
5. **Mock dependencies** appropriately

### Test Structure

```go
func TestMyFunction_WhenValidInput_ThenReturnsExpected(t *testing.T) {
    // Arrange
    input := "test"
    expected := "result"
    
    // Act
    result := MyFunction(input)
    
    // Assert
    assert.Equal(t, expected, result)
}
```

### Table-Driven Tests

```go
func TestMyFunction_VariousInputs(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"empty string", "", ""},
        {"valid input", "test", "result"},
        {"special chars", "!@#", "escaped"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := MyFunction(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Coverage Requirements

- Aim for **80%+ code coverage** on new code
- All public APIs must have tests
- Test both success and error paths
- Include integration tests where appropriate

## Pull Request Process

### Before Submitting

1. **Update your branch** with upstream changes:
   ```bash
   git fetch upstream
   git rebase upstream/master
   ```

2. **Run all tests** and ensure they pass:
   ```bash
   go test ./...
   ```

3. **Run linters**:
   ```bash
   go fmt ./...
   go vet ./...
   ```

4. **Update documentation** if you've changed APIs

5. **Add/update tests** for your changes

### Creating the Pull Request

1. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create a pull request** from your fork to `janmbaco/go-infrastructure:master`

3. **Fill out the PR template** with:
   - Description of changes
   - Related issue numbers (if applicable)
   - Testing performed
   - Breaking changes (if any)

4. **Wait for review** - maintainers will review your PR and may request changes

### PR Review Criteria

- Code follows project style guidelines
- Tests pass and provide adequate coverage
- Documentation is updated
- Commits are clear and follow commit message guidelines
- No merge conflicts
- Changes are focused and purposeful

## Commit Message Guidelines

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples

```
feat(dependencyinjection): add full context.Context support

- Update DependencyObject.Create to accept context
- Implement automatic detection of context-aware providers
- Add cancellation checks during resolution
- Propagate context through nested dependencies

Closes #123
```

```
fix(errors): handle nil logger in ErrorCatcher constructor

Previously, NewErrorCatcher would panic if logger was nil.
Now it uses a no-op logger as fallback.

Fixes #456
```

### Breaking Changes

If your commit introduces breaking changes, add `BREAKING CHANGE:` in the footer:

```
feat(persistence)!: change DatabaseInfo constructor signature

BREAKING CHANGE: NewDB now requires context.Context as first parameter.
Migration: Add ctx parameter to all NewDB calls.
```

## Documentation

### Code Documentation

- **All exported functions, types, and packages** must have documentation comments
- Start comments with the name of the item being documented
- Use complete sentences
- Explain *what* and *why*, not just *how*

```go
// NewContainer creates a new dependency injection container.
// The container is initialized with default registrations for
// Container, Register, and Resolver interfaces.
func NewContainer() Container {
    // implementation
}
```

### README Updates

If your changes affect:
- **Public APIs** - update relevant README sections
- **Examples** - add or update code examples
- **Configuration** - document new options

### Package Documentation

Each package should have a `doc.go` file or package-level comment explaining:
- Purpose of the package
- Main concepts
- Basic usage example
- Links to more detailed documentation

## Community

### Getting Help

- **GitHub Discussions** - for questions and discussions
- **GitHub Issues** - for bug reports and feature requests
- **Pull Requests** - for code reviews and contributions

### Recognition

Contributors are recognized in:
- Git commit history
- Release notes
- CHANGELOG.md
- GitHub contributors page

## Release Process

Maintainers handle releases, but here's the process:

1. Update `CHANGELOG.md`
2. Update version numbers
3. Create a release tag (`v2.x.x`)
4. Push tag to trigger release workflow
5. Publish release notes on GitHub

## Questions?

If you have questions about contributing, feel free to:
- Open a GitHub Discussion
- Comment on relevant issues
- Reach out to maintainers

---

Thank you for contributing to `go-infrastructure`! 🎉
