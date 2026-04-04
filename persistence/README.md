# Persistence

`github.com/janmbaco/go-infrastructure/v2/persistence`

The `persistence` package provides database wiring and typed CRUD helpers on top of GORM. It supports PostgreSQL, MySQL, SQLite and SQL Server, and it is designed to work both with the DI container in this repository and with plain `*gorm.DB` usage.

## What It Includes

- `DatabaseInfo` and `DbEngine` to describe database connections
- `NewDB` to build a `*gorm.DB` from a dialector resolver
- `ioc.ConfigureDatabaseModule` for the simplest DI-based setup
- `dataaccess.NewTypedDataAccess[T]` and generic CRUD helpers
- `DB()` access for advanced GORM queries when CRUD helpers are not enough

## Install

```bash
go get github.com/janmbaco/go-infrastructure/v2/persistence
```

## Recommended Setup

For most applications, the simplest path is:

1. register a database module with `ioc.ConfigureDatabaseModule`
2. resolve `*gorm.DB`
3. build typed accessors with `dataaccess.NewTypedDataAccess[T]`

```go
package main

import (
    di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
    "github.com/janmbaco/go-infrastructure/v2/persistence"
    persistenceioc "github.com/janmbaco/go-infrastructure/v2/persistence/ioc"
    "github.com/janmbaco/go-infrastructure/v2/persistence/dataaccess"
    "gorm.io/gorm"
)

type User struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"size:100;not null"`
    Email string `gorm:"size:100;unique;not null"`
}

func main() {
    container := di.NewBuilder().
        AddModule(persistenceioc.ConfigureDatabaseModule(
            "", "", "", "", "./app.db", persistence.Sqlite,
        )).
        MustBuild()

    resolver := container.Resolver()
    db := resolver.Type(new(*gorm.DB), nil).(*gorm.DB)

    userAccess := dataaccess.NewTypedDataAccess[User](db)
    _, _ = dataaccess.SelectRows(userAccess, &User{})
}
```

For PostgreSQL, MySQL and SQL Server, pass host, port, user, password and database name to `ConfigureDatabaseModule`. For SQLite, pass the database path as `dbName`.

## Direct `NewDB` Usage

If you want to construct the GORM connection directly, use `NewDB` with a dialector resolver:

```go
package main

import (
    di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
    "github.com/janmbaco/go-infrastructure/v2/persistence"
    dialectorsioc "github.com/janmbaco/go-infrastructure/v2/persistence/dialectors/ioc"
    "gorm.io/gorm"
)

type User struct {
    ID    uint   `gorm:"primaryKey"`
    Email string `gorm:"size:100;unique;not null"`
}

func main() {
    container := di.NewBuilder().
        AddModule(dialectorsioc.NewDialectorsModule()).
        MustBuild()

    resolver := persistence.NewDialectorResolver(container.Resolver())

    db, err := persistence.NewDB(
        resolver,
        &persistence.DatabaseInfo{
            Name:   "./app.db",
            Engine: persistence.Sqlite,
        },
        &gorm.Config{},
        []interface{}{&User{}},
    )
    if err != nil {
        panic(err)
    }

    _ = db
}
```

This path is useful when you want explicit control over `gorm.Config` and startup migrations.

## Core Types

Connection metadata:

```go
type DatabaseInfo struct {
    Host         string
    Port         string
    Name         string
    UserName     string
    UserPassword string
    Engine       DbEngine
}

const (
    SQLServer DbEngine = iota
    Postgres
    MySQL
    Sqlite
)
```

CRUD abstraction:

```go
type DataAccess interface {
    Insert(datarow interface{}) error
    Select(datafilter interface{}, preloads ...string) (interface{}, error)
    Update(datafilter interface{}, datarow interface{}) error
    Delete(datafilter interface{}, associateds ...string) error
    DB() interface{}
}
```

`Select` uses exact-match struct filters. For range queries, joins or engine-specific SQL, use the underlying `*gorm.DB` through `DB()`.

## CRUD Example

```go
userAccess := dataaccess.NewTypedDataAccess[User](db)

user := &User{Name: "Alice", Email: "alice@example.com"}
if err := dataaccess.InsertRow(userAccess, user); err != nil {
    panic(err)
}

users, err := dataaccess.SelectRows(userAccess, &User{Email: "alice@example.com"})
if err != nil {
    panic(err)
}

user.Name = "Alice Johnson"
if err := dataaccess.UpdateRow(userAccess, &User{ID: user.ID}, user); err != nil {
    panic(err)
}

if err := dataaccess.DeleteRows(userAccess, &User{ID: user.ID}); err != nil {
    panic(err)
}

_ = users
```

## Associations and Preloads

Preloads are passed directly to `SelectRows`:

```go
type Profile struct {
    ID     uint   `gorm:"primaryKey"`
    UserID uint   `gorm:"not null"`
    Bio    string `gorm:"size:500"`
    User   User   `gorm:"foreignKey:UserID"`
}

profileAccess := dataaccess.NewTypedDataAccess[Profile](db)

profiles, err := dataaccess.SelectRows(profileAccess, &Profile{}, "User")
if err != nil {
    panic(err)
}

_ = profiles
```

Associated deletions can be requested through `DeleteRows`:

```go
if err := dataaccess.DeleteRows(profileAccess, &Profile{UserID: user.ID}, "User"); err != nil {
    panic(err)
}
```

## Advanced Queries

`DB()` exposes the wrapped `*gorm.DB` for queries that do not fit the generic CRUD helpers:

```go
gormDB := userAccess.DB().(*gorm.DB)

var stats struct {
    TotalUsers  int
    ActiveUsers int
}

err := gormDB.Raw(`
    SELECT
        COUNT(*) as total_users,
        SUM(CASE WHEN active = 1 THEN 1 ELSE 0 END) as active_users
    FROM users
`).Scan(&stats).Error
if err != nil {
    panic(err)
}
```

## Database Support

- PostgreSQL
- MySQL
- SQLite
- SQL Server

`ioc.ConfigureDatabaseModule` selects the right dialector based on `DbEngine`.

## Testing

Unit tests:

```bash
go test ./persistence/...
```

SQLite integration coverage:

```bash
go test -tags=integration -v ./persistence/integration_test -run SQLite
```

Docker-backed integration tests for PostgreSQL, MySQL and SQL Server:

```bash
cd persistence/integration_test
.\run-integration-tests.ps1
```

Additional notes for the integration suite are documented in [INTEGRATION_TESTS.md](./integration_test/INTEGRATION_TESTS.md).

## Related Packages

- `persistence/dataaccess`: typed CRUD helpers
- `persistence/ioc`: DI modules and convenience helpers
- `persistence/dialectors`: per-engine dialector getters
- `persistence/integration_test`: integration coverage and scripts
