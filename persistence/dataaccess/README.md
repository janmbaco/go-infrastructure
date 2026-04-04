# DataAccess

`github.com/janmbaco/go-infrastructure/v2/persistence/dataaccess`

The `dataaccess` package wraps GORM with a small CRUD-oriented interface and generic helper functions.

## What It Includes

- `DataAccess`
- `NewDataAccess`
- `NewTypedDataAccess[T]`
- generic helpers: `InsertRow`, `SelectRows`, `UpdateRow`, `DeleteRows`
- `GormDBInterface` and `NewDataAccessWithInterface` for testing

## Install

```bash
go get github.com/janmbaco/go-infrastructure/v2/persistence/dataaccess
go get gorm.io/gorm
```

## Quick Start

```go
package main

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"

    "github.com/janmbaco/go-infrastructure/v2/persistence/dataaccess"
)

type User struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"size:100"`
    Email string `gorm:"size:100;unique"`
}

func main() {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic(err)
    }

    if err := db.AutoMigrate(&User{}); err != nil {
        panic(err)
    }

    userDA := dataaccess.NewTypedDataAccess[User](db)

    user := &User{Name: "Alice", Email: "alice@example.com"}
    if err := dataaccess.InsertRow(userDA, user); err != nil {
        panic(err)
    }

    users, err := dataaccess.SelectRows(userDA, &User{Email: "alice@example.com"})
    if err != nil {
        panic(err)
    }

    _ = users
}
```

## Core Interface

```go
type DataAccess interface {
    Insert(datarow interface{}) error
    Select(datafilter interface{}, preloads ...string) (interface{}, error)
    Update(datafilter interface{}, datarow interface{}) error
    Delete(datafilter interface{}, associateds ...string) error
    DB() interface{}
}
```

`DB()` exposes the wrapped `*gorm.DB` for more complex queries.

## Generic Helpers

```go
func NewTypedDataAccess[T any](db *gorm.DB) DataAccess
func NewTypedDataAccessWithInterface[T any](db GormDBInterface) DataAccess

func InsertRow[T any](da DataAccess, datarow *T) error
func SelectRows[T any](da DataAccess, datafilter *T, preloads ...string) ([]*T, error)
func UpdateRow[T any](da DataAccess, datafilter *T, datarow *T) error
func DeleteRows[T any](da DataAccess, datafilter *T, associateds ...string) error
```

These helpers keep application code strongly typed while still using the underlying `DataAccess` abstraction.

## Preloads and Associated Deletes

Preloads are passed to `SelectRows`:

```go
profiles, err := dataaccess.SelectRows(profileDA, &Profile{}, "User")
```

Associated deletions are passed to `DeleteRows`:

```go
err := dataaccess.DeleteRows(profileDA, &Profile{UserID: userID}, "User")
```

## Advanced Queries

For queries that do not fit exact-match struct filters, use `DB()`:

```go
gormDB := userDA.DB().(*gorm.DB)

var stats struct {
    TotalUsers int
}

if err := gormDB.Raw(`SELECT COUNT(*) as total_users FROM users`).Scan(&stats).Error; err != nil {
    panic(err)
}
```

## Testing

`NewDataAccessWithInterface` and `NewTypedDataAccessWithInterface[T]` exist to make unit tests easier by targeting `GormDBInterface` instead of a real `*gorm.DB`.

Unit tests for this package:

```bash
go test ./persistence/dataaccess/...
```

Integration coverage lives under `persistence/integration_test`.

## Related Packages

- `persistence`: database setup and engine selection
- `persistence/ioc`: DI helpers
- `persistence/integration_test`: integration coverage
