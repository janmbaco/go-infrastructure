# DataAccess Module

`github.com/janmbaco/go-infrastructure/v2/persistence/dataaccess`

Type-safe GORM wrapper with generics support for simplified database operations.

## Overview

The dataaccess module provides a clean, type-safe abstraction over GORM with:

- **Generic functions** - Type-safe operations with zero casting
- **Simple CRUD** - Insert, Select, Update, Delete operations
- **Preload support** - Eager loading of relationships
- **Cascade delete** - Delete with associated records
- **Testable design** - Interface-based for easy mocking

## Installation

```bash
go get github.com/janmbaco/go-infrastructure/v2/persistence/dataaccess
go get gorm.io/gorm
```

## Quick Start

### Basic Usage with Generics

```go
package main

import (
    "fmt"
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
    // Connect to database
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic(err)
    }

    // Auto migrate
    db.AutoMigrate(&User{})

    // Create typed data access
    userDA := dataaccess.NewTypedDataAccess[User](db)

    // Insert
    user := &User{Name: "Alice", Email: "alice@example.com"}
    if err := dataaccess.InsertRow(userDA, user); err != nil {
        panic(err)
    }

    // Select
    filter := &User{Email: "alice@example.com"}
    users, err := dataaccess.SelectRows(userDA, filter)
    if err != nil {
        panic(err)
    }
    fmt.Println("Found:", users[0].Name)

    // Update
    update := &User{Name: "Alice Updated"}
    if err := dataaccess.UpdateRow(userDA, filter, update); err != nil {
        panic(err)
    }

    // Delete
    if err := dataaccess.DeleteRows(userDA, filter); err != nil {
        panic(err)
    }
}
```

## API Reference

### Generic Functions

```go
// NewTypedDataAccess creates a type-safe DataAccess instance
func NewTypedDataAccess[T any](db *gorm.DB) DataAccess

// InsertRow inserts a row with type safety
func InsertRow[T any](da DataAccess, datarow *T) error

// SelectRows selects rows with type safety
func SelectRows[T any](da DataAccess, datafilter *T, preloads ...string) ([]*T, error)

// UpdateRow updates rows with type safety
func UpdateRow[T any](da DataAccess, datafilter *T, datarow *T) error

// DeleteRows deletes rows with type safety
func DeleteRows[T any](da DataAccess, datafilter *T, associateds ...string) error
```

### DataAccess Interface

```go
type DataAccess interface {
    // Insert creates a new record
    Insert(datarow interface{}) error
    
    // Select retrieves records matching the filter
    Select(datafilter interface{}, preloads ...string) (interface{}, error)
    
    // Update modifies records matching the filter
    Update(datafilter interface{}, datarow interface{}) error
    
    // Delete removes records matching the filter
    Delete(datafilter interface{}, associateds ...string) error
    
    // DB returns the underlying database instance
    DB() interface{}
}
```

### Non-Generic Constructors

```go
// NewDataAccess creates DataAccess with reflection
func NewDataAccess(db *gorm.DB, modelType reflect.Type) DataAccess

// NewDataAccessWithInterface creates DataAccess with mockable interface
func NewDataAccessWithInterface(db GormDBInterface, modelType reflect.Type) DataAccess

// NewTypedDataAccessWithInterface creates typed DataAccess for testing
func NewTypedDataAccessWithInterface[T any](db GormDBInterface) DataAccess
```

## Usage Examples

### Complete CRUD Operations

```go
package main

import (
    "fmt"
    "time"
    
    "gorm.io/gorm"
    "github.com/janmbaco/go-infrastructure/v2/persistence/dataaccess"
)

type Post struct {
    ID        uint      `gorm:"primaryKey"`
    Title     string    `gorm:"size:200;not null"`
    Content   string    `gorm:"type:text"`
    AuthorID  uint      `gorm:"not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func main() {
    db := setupDatabase() // Your DB setup
    postDA := dataaccess.NewTypedDataAccess[Post](db)

    // CREATE
    newPost := &Post{
        Title:    "Getting Started with Go",
        Content:  "Go is an amazing language...",
        AuthorID: 1,
    }
    if err := dataaccess.InsertRow(postDA, newPost); err != nil {
        panic(err)
    }
    fmt.Println("Created post with ID:", newPost.ID)

    // READ - Find by ID
    filter := &Post{ID: newPost.ID}
    posts, err := dataaccess.SelectRows(postDA, filter)
    if err != nil {
        panic(err)
    }
    if len(posts) > 0 {
        fmt.Println("Post title:", posts[0].Title)
    }

    // READ - Find by author
    authorFilter := &Post{AuthorID: 1}
    authorPosts, err := dataaccess.SelectRows(postDA, authorFilter)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Found %d posts by author\n", len(authorPosts))

    // UPDATE
    updateData := &Post{
        Title:   "Updated: Getting Started with Go",
        Content: "Updated content...",
    }
    if err := dataaccess.UpdateRow(postDA, filter, updateData); err != nil {
        panic(err)
    }
    fmt.Println("Post updated")

    // DELETE
    if err := dataaccess.DeleteRows(postDA, filter); err != nil {
        panic(err)
    }
    fmt.Println("Post deleted")
}
```

### Working with Relationships

```go
type User struct {
    ID      uint      `gorm:"primaryKey"`
    Name    string    `gorm:"size:100"`
    Email   string    `gorm:"size:100;unique"`
    Profile Profile   `gorm:"foreignKey:UserID"`
    Posts   []Post    `gorm:"foreignKey:AuthorID"`
}

type Profile struct {
    ID     uint   `gorm:"primaryKey"`
    UserID uint   `gorm:"not null;unique"`
    Bio    string `gorm:"size:500"`
    Avatar string `gorm:"size:255"`
}

type Post struct {
    ID       uint   `gorm:"primaryKey"`
    AuthorID uint   `gorm:"not null"`
    Title    string `gorm:"size:200"`
    Content  string `gorm:"type:text"`
}

func main() {
    db := setupDatabase()
    userDA := dataaccess.NewTypedDataAccess[User](db)

    // Create user with profile
    user := &User{
        Name:  "Bob",
        Email: "bob@example.com",
        Profile: Profile{
            Bio:    "Software developer",
            Avatar: "avatar.jpg",
        },
    }
    dataaccess.InsertRow(userDA, user)

    // Select with preloads
    filter := &User{Email: "bob@example.com"}
    users, err := dataaccess.SelectRows(userDA, filter, "Profile", "Posts")
    if err != nil {
        panic(err)
    }

    if len(users) > 0 {
        fmt.Println("User:", users[0].Name)
        fmt.Println("Bio:", users[0].Profile.Bio)
        fmt.Printf("Posts count: %d\n", len(users[0].Posts))
    }

    // Delete with associated records
    // This will delete the user AND their profile
    if err := dataaccess.DeleteRows(userDA, filter, "Profile"); err != nil {
        panic(err)
    }
}
```

### Partial Updates

```go
func updateUserEmail(userDA dataaccess.DataAccess, userID uint, newEmail string) error {
    // Filter: which records to update
    filter := &User{ID: userID}
    
    // Update: only the fields you specify
    update := &User{Email: newEmail}
    
    return dataaccess.UpdateRow(userDA, filter, update)
}

func updateUserProfile(userDA dataaccess.DataAccess, userID uint, name, email string) error {
    filter := &User{ID: userID}
    update := &User{
        Name:  name,
        Email: email,
    }
    return dataaccess.UpdateRow(userDA, filter, update)
}
```

### Complex Queries with Raw DB Access

```go
func findActiveUsers(userDA dataaccess.DataAccess) ([]*User, error) {
    // Get the underlying GORM DB
    db := userDA.DB().(*gorm.DB)
    
    var users []*User
    err := db.Where("last_login > ?", time.Now().AddDate(0, -1, 0)).
        Order("name ASC").
        Limit(100).
        Find(&users).Error
    
    return users, err
}

func countUsersByDomain(userDA dataaccess.DataAccess, domain string) (int64, error) {
    db := userDA.DB().(*gorm.DB)
    
    var count int64
    err := db.Model(&User{}).
        Where("email LIKE ?", "%@"+domain).
        Count(&count).Error
    
    return count, err
}
```

### Batch Operations

```go
func insertMultipleUsers(userDA dataaccess.DataAccess, users []*User) error {
    db := userDA.DB().(*gorm.DB)
    
    // GORM batch insert
    return db.Create(users).Error
}

func updateMultipleUsers(userDA dataaccess.DataAccess, ids []uint, status string) error {
    db := userDA.DB().(*gorm.DB)
    
    return db.Model(&User{}).
        Where("id IN ?", ids).
        Update("status", status).Error
}
```

### Testing with Mocks

```go
package mypackage_test

import (
    "testing"
    
    "github.com/janmbaco/go-infrastructure/v2/persistence/dataaccess"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// MockGormDB implements GormDBInterface for testing
type MockGormDB struct {
    mock.Mock
}

func (m *MockGormDB) Create(value interface{}) dataaccess.GormDBInterface {
    args := m.Called(value)
    return args.Get(0).(dataaccess.GormDBInterface)
}

func (m *MockGormDB) Find(dest interface{}) dataaccess.GormDBInterface {
    args := m.Called(dest)
    return args.Get(0).(dataaccess.GormDBInterface)
}

func (m *MockGormDB) GetError() error {
    args := m.Called()
    return args.Error(0)
}

// ... implement other methods

func TestUserService_CreateUser(t *testing.T) {
    mockDB := new(MockGormDB)
    userDA := dataaccess.NewTypedDataAccessWithInterface[User](mockDB)

    // Setup expectations
    mockDB.On("Model", mock.Anything).Return(mockDB)
    mockDB.On("Create", mock.Anything).Return(mockDB)
    mockDB.On("GetError").Return(nil)

    // Test
    user := &User{Name: "Test", Email: "test@example.com"}
    err := dataaccess.InsertRow(userDA, user)

    assert.NoError(t, err)
    mockDB.AssertExpectations(t)
}
```

## Integration with DI Container

```go
package main

import (
    "gorm.io/gorm"
    
    di "github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
    "github.com/janmbaco/go-infrastructure/v2/persistence/dataaccess"
    "github.com/janmbaco/go-infrastructure/v2/persistence/ioc"
)

func main() {
    // Setup container with database module
    container := di.NewBuilder().
        AddModule(ioc.ConfigureDatabaseModule(
            "localhost", "5432", "user", "pass", "mydb", "postgres",
        )).
        MustBuild()

    resolver := container.Resolver()

    // Resolve GORM DB
    db := resolver.Type(new(*gorm.DB), nil).(*gorm.DB)

    // Create data access instances
    userDA := dataaccess.NewTypedDataAccess[User](db)
    postDA := dataaccess.NewTypedDataAccess[Post](db)

    // Use them
    user := &User{Name: "Alice", Email: "alice@example.com"}
    dataaccess.InsertRow(userDA, user)
}
```

## Error Handling

The module defines specific error types:

```go
// Check for specific errors
if err := dataaccess.InsertRow(userDA, user); err != nil {
    if errors.Is(err, gorm.ErrDuplicatedKey) {
        return fmt.Errorf("user already exists: %w", err)
    }
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return fmt.Errorf("user not found: %w", err)
    }
    return fmt.Errorf("database error: %w", err)
}
```

### Custom Error Types

```go
const (
    DataFilterUnexpected = "DATA_FILTER_UNEXPECTED"
)

// DataBaseError wraps database errors
type DataBaseError struct {
    Code    string
    Message string
    Err     error
}
```

## Best Practices

### 1. Use Generics for Type Safety

```go
// GOOD - Type-safe, compile-time checking
userDA := dataaccess.NewTypedDataAccess[User](db)
users, err := dataaccess.SelectRows(userDA, filter)
// users is []*User, no casting needed

// AVOID - Requires type assertion
userDA := dataaccess.NewDataAccess(db, reflect.TypeOf(&User{}))
result, err := userDA.Select(filter)
users := result.([]*User) // Runtime error if wrong type
```

### 2. Always Use Pointers for Models

```go
// GOOD
type User struct {
    ID   uint
    Name string
}
user := &User{Name: "Alice"}
dataaccess.InsertRow(userDA, user)

// BAD - Won't populate ID after insert
user := User{Name: "Alice"}
dataaccess.InsertRow(userDA, &user)
```

### 3. Leverage Preloads for Relationships

```go
// Efficient - One query with joins
filter := &User{ID: 1}
users, _ := dataaccess.SelectRows(userDA, filter, "Profile", "Posts")

// Inefficient - Multiple queries (N+1 problem)
users, _ := dataaccess.SelectRows(userDA, filter)
for _, user := range users {
    // Separate queries for each user's profile and posts
}
```

### 4. Use Cascade Delete Carefully

```go
// Delete user and their profile
filter := &User{ID: 1}
dataaccess.DeleteRows(userDA, filter, "Profile")

// Delete user but keep orphaned records
dataaccess.DeleteRows(userDA, filter)
```

### 5. Access Raw DB for Complex Queries

```go
// For complex queries, use the underlying GORM DB
db := userDA.DB().(*gorm.DB)

db.Where("age > ?", 18).
    Where("status = ?", "active").
    Order("created_at DESC").
    Limit(10).
    Find(&users)
```

## Performance Considerations

### Batch Inserts

```go
// Efficient batch insert
users := []*User{
    {Name: "Alice", Email: "alice@example.com"},
    {Name: "Bob", Email: "bob@example.com"},
    // ... more users
}

db := userDA.DB().(*gorm.DB)
db.CreateInBatches(users, 100) // Insert in batches of 100
```

### Select Only Required Fields

```go
db := userDA.DB().(*gorm.DB)

// Only select specific columns
var users []User
db.Select("id", "name", "email").Find(&users)

// Omit large columns
db.Omit("content", "metadata").Find(&posts)
```

### Use Transactions for Multiple Operations

```go
db := userDA.DB().(*gorm.DB)

err := db.Transaction(func(tx *gorm.DB) error {
    // Create user
    if err := tx.Create(&user).Error; err != nil {
        return err
    }
    
    // Create profile
    profile.UserID = user.ID
    if err := tx.Create(&profile).Error; err != nil {
        return err
    }
    
    return nil
})
```

## Common Patterns

### Repository Pattern

```go
type UserRepository struct {
    da dataaccess.DataAccess
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{
        da: dataaccess.NewTypedDataAccess[User](db),
    }
}

func (r *UserRepository) Create(user *User) error {
    return dataaccess.InsertRow(r.da, user)
}

func (r *UserRepository) FindByEmail(email string) (*User, error) {
    filter := &User{Email: email}
    users, err := dataaccess.SelectRows(r.da, filter)
    if err != nil {
        return nil, err
    }
    if len(users) == 0 {
        return nil, gorm.ErrRecordNotFound
    }
    return users[0], nil
}

func (r *UserRepository) Update(user *User) error {
    filter := &User{ID: user.ID}
    return dataaccess.UpdateRow(r.da, filter, user)
}

func (r *UserRepository) Delete(id uint) error {
    filter := &User{ID: id}
    return dataaccess.DeleteRows(r.da, filter)
}
```

### Service Layer

```go
type UserService struct {
    repo *UserRepository
}

func (s *UserService) RegisterUser(name, email string) (*User, error) {
    // Validate
    if !isValidEmail(email) {
        return nil, errors.New("invalid email")
    }

    // Check if exists
    if _, err := s.repo.FindByEmail(email); err == nil {
        return nil, errors.New("user already exists")
    }

    // Create
    user := &User{Name: name, Email: email}
    if err := s.repo.Create(user); err != nil {
        return nil, err
    }

    return user, nil
}
```

## Troubleshooting

### "datafilter does not belong to this dataAccess" error

**Cause:** Using wrong model type in filter.

**Solution:** Ensure filter matches the DataAccess type:
```go
userDA := dataaccess.NewTypedDataAccess[User](db)

// CORRECT
filter := &User{ID: 1}
users, _ := dataaccess.SelectRows(userDA, filter)

// WRONG - Different type
filter := &Post{ID: 1}
users, _ := dataaccess.SelectRows(userDA, filter) // ERROR!
```

### Empty results when data exists

**Cause:** Filter fields don't match database values.

**Solution:** Use exact values or raw DB for flexible queries:
```go
// Exact match required
filter := &User{Name: "Alice"} // Must match exactly

// For LIKE or other operators, use raw DB
db := userDA.DB().(*gorm.DB)
db.Where("name LIKE ?", "%alice%").Find(&users)
```

### Updates not working

**Cause:** Zero values are ignored by GORM.

**Solution:** Use `Select` or raw DB:
```go
// Zero value ignored
update := &User{Age: 0} // Won't update to 0

// Force update with Select
db := userDA.DB().(*gorm.DB)
db.Model(&User{}).Where("id = ?", 1).Select("Age").Updates(update)
```

## Migration from Raw GORM

```go
// Before (Raw GORM)
var users []User
db.Where(&User{Status: "active"}).Find(&users)

// After (DataAccess)
userDA := dataaccess.NewTypedDataAccess[User](db)
filter := &User{Status: "active"}
users, err := dataaccess.SelectRows(userDA, filter)
```

## Contributing

See [CONTRIBUTING.md](../../CONTRIBUTING.md)

## License

Apache License 2.0 - see [LICENSE](../../LICENSE)
