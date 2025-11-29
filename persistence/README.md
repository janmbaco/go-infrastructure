# 🚀 Persistence Package - Simple & Powerful Data Access

> **One-line setup, zero boilerplate** - Connect to any database and start querying in minutes!

The `persistence` package gives you a clean, type-safe way to work with databases in Go. Whether you're building a simple app or a complex enterprise system, this package makes database operations feel natural and enjoyable.

## ✨ What's So Great About It?

- **🔧 Zero Configuration**: Just add one module and you're ready to query
- **🎯 Type-Safe**: Full Go generics support with compile-time safety
- **🗄️ Multi-Database**: PostgreSQL, MySQL, SQLite, SQL Server - all supported
- **🔄 Dependency Injection**: Seamlessly integrates with your DI container
- **🧪 Battle-Tested**: Comprehensive tests with real databases
- **📚 Simple API**: Just 4 methods: Insert, Select, Update, Delete

## 🚀 Quick Start (5 Minutes!)

### 1. Choose Your Database

```go
// PostgreSQL
dbInfo := &persistence.DatabaseInfo{
    Host: "localhost", Port: "5432", Name: "myapp",
    UserName: "myuser", UserPassword: "mypass",
    Engine: persistence.Postgres,
}

// MySQL
dbInfo := &persistence.DatabaseInfo{
    Host: "localhost", Port: "3306", Name: "myapp",
    UserName: "myuser", UserPassword: "mypass",
    Engine: persistence.MySQL,
}

// SQLite (super simple!)
dbInfo := &persistence.DatabaseInfo{
    Name: "/path/to/myapp.db",
    Engine: persistence.Sqlite,
}
```

### 2. Setup Your Container (One Line!)

```go
container := dependencyinjection.NewBuilder().
    AddModule(persistenceioc.ConfigureDatabaseModule(
        dbInfo.Host, dbInfo.Port, dbInfo.UserName,
        dbInfo.UserPassword, dbInfo.Name, dbInfo.Engine,
    )).
    MustBuild()

resolver := container.Resolver()
```

### 3. Define Your Models

```go
type User struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `gorm:"size:100;not null"`
    Email string `gorm:"size:100;unique;not null"`
    Age   int    `gorm:"not null"`
}

type Profile struct {
    ID     uint   `gorm:"primaryKey"`
    UserID uint   `gorm:"not null"`
    Bio    string `gorm:"size:500"`
    User   User   `gorm:"foreignKey:UserID"`
}
```

### 4. Start Querying! 🎉

```go
// Get typed data access (that's it!)
userAccess := dataaccess.NewTypedDataAccess[User](resolver.Type(new(*gorm.DB), nil).(*gorm.DB))

// Create a user
user := &User{Name: "Alice", Email: "alice@example.com", Age: 28}
err := dataaccess.InsertRow(userAccess, user)

// Find users
users, err := dataaccess.SelectRows(userAccess, &User{Email: "alice@example.com"})

// Update user
user.Name = "Alice Johnson"
err = dataaccess.UpdateRow(userAccess, &User{ID: user.ID}, user)

// Delete user
err = dataaccess.DeleteRows(userAccess, &User{ID: user.ID})
```

**That's it!** You're now doing type-safe database operations with full CRUD support.

## 🔧 Advanced Queries with DB() Method

For complex queries that go beyond the basic CRUD operations, you can access the underlying GORM database directly:

```go
// Get raw GORM instance for advanced operations
gormDB := userAccess.DB().(*gorm.DB)

// Complex queries with joins, aggregations, etc.
var result []User
err := gormDB.
    Select("users.name, COUNT(orders.id) as order_count").
    Joins("LEFT JOIN orders ON users.id = orders.user_id").
    Group("users.id").
    Having("order_count > ?", 5).
    Find(&result).Error

// Raw SQL when needed (still parameterized!)
var stats struct {
    TotalUsers int
    ActiveUsers int
}
err := gormDB.Raw(`
    SELECT 
        COUNT(*) as total_users,
        COUNT(CASE WHEN last_login > ? THEN 1 END) as active_users
    FROM users
`, time.Now().Add(-30*24*time.Hour)).Scan(&stats).Error
```

> **Security Note**: Even with raw access, GORM automatically parameterizes queries to prevent SQL injection.

## 📖 Core Concepts

### DataAccess Interface

Everything revolves around this simple interface:

```go
type DataAccess interface {
    Insert(datarow interface{}) error           // Create
    Select(datafilter interface{}, preloads ...string) (interface{}, error)  // Read
    Update(datafilter interface{}, datarow interface{}) error   // Update
    Delete(datafilter interface{}, associateds ...string) error // Delete
}
```

### Type-Safe Operations

No more `interface{}` casting or reflection headaches:

```go
// ✅ Type-safe: Compiler catches errors
userAccess := dataaccess.NewTypedDataAccess[User](db)
users, err := dataaccess.SelectRows[User](userAccess, &User{Name: "Alice"})

// ❌ Old way: Runtime errors possible
genericAccess := someGenericDataAccess(db)
users, err := genericAccess.Select(&User{Name: "Alice"})
```

## 🗄️ Database Support

### PostgreSQL
```go
container := dependencyinjection.NewBuilder().
    AddModule(persistenceioc.ConfigureDatabaseModule(
        "localhost", "5432", "myuser", "mypass", "myapp", persistence.Postgres,
    )).
    MustBuild()
```

### MySQL
```go
container := dependencyinjection.NewBuilder().
    AddModule(persistenceioc.ConfigureDatabaseModule(
        "localhost", "3306", "myuser", "mypass", "myapp", persistence.MySQL,
    )).
    MustBuild()
```

### SQLite
```go
container := dependencyinjection.NewBuilder().
    AddModule(persistenceioc.ConfigureDatabaseModule(
        "", "", "", "", "/path/to/app.db", persistence.Sqlite,
    )).
    MustBuild()
```

### SQL Server
```go
container := dependencyinjection.NewBuilder().
    AddModule(persistenceioc.ConfigureDatabaseModule(
        "localhost", "1433", "sa", "MyPass123!", "myapp", persistence.SQLServer,
    )).
    MustBuild()
```

## 🎯 Advanced Usage

### Working with Relationships

```go
// Get user with profile
profileAccess := dataaccess.NewTypedDataAccess[Profile](db)
profiles, err := dataaccess.SelectRows[Profile](profileAccess, &Profile{}, "User")

// Delete user and profile together
err = dataaccess.DeleteRows(profileAccess, &Profile{UserID: userID}, "User")
```

### Custom Queries with Filters

```go
// Find users exactly 18 years old
eighteenYearOlds, err := dataaccess.SelectRows(userAccess, &User{Age: 18})

// Find by multiple fields (all must match exactly)
activeUser := &User{Email: "alice@example.com", Age: 28}
user, err := dataaccess.SelectRows(userAccess, activeUser)

// Find all users (empty filter)
allUsers, err := dataaccess.SelectRows(userAccess, &User{})
```

**Note**: Filters use exact matching by default. For range queries or complex conditions, you'll need to use GORM's query builder directly.

### Error Handling

The package provides clear, actionable errors:

```go
user := &User{Name: "Bob", Email: "alice@example.com"} // Duplicate email
err := dataaccess.InsertRow(userAccess, user)
if err != nil {
    // Handle constraint violation
    log.Printf("Failed to create user: %v", err)
}
```

## 🧪 Testing

### Unit Tests
```bash
go test ./persistence/dataaccess/...
```

### Integration Tests (Real Databases!)
```bash
# Uses Docker containers - no local DB setup needed!
cd persistence/integration_test
.\run-integration-tests.ps1
```

The integration tests automatically:
- 🐳 Spin up PostgreSQL, MySQL, and SQL Server containers
- 🔍 Test all CRUD operations
- 🧹 Clean up everything when done

## 📁 Package Structure

```
persistence/
├── dataaccess/           # Core CRUD operations
├── dialectors/           # Database drivers
├── ioc/                  # Dependency injection modules
├── integration_test/     # Real database testing
└── README.md            # You are here! 👋
```

## 🔄 Migration Guide

Coming from the old `orm_base` package? Here's what changed:

- ✅ **Same API**: Your existing code works unchanged
- ✅ **Better Organization**: Cleaner package structure
- ✅ **New Convenience Functions**: `ConfigureDatabaseModule()` for easier setup
- ✅ **Enhanced Testing**: Docker-based integration tests

## 🤝 Contributing

We love contributions! Here's how to get started:

1. **Fork & Clone**: Get the code locally
2. **Test**: `go test ./...` (unit) and `.\run-integration-tests.ps1` (integration)
3. **Add**: Your awesome feature or fix
4. **Test Again**: Make sure everything still works
5. **PR**: Submit your changes!

## 📚 Examples in the Wild

Check out how this package is used in the codebase:

- **Server Setup**: `server/facades/singlepageapp_facade.go`
- **Integration Tests**: `persistence/integration_test/`
- **Module Configuration**: `persistence/ioc/`

---

**Ready to simplify your database code?** Just add one module and start querying! 🚀

*Built with ❤️ for the Go community*</content>
<parameter name="filePath">d:\go-infrastructure\persistence\README.md