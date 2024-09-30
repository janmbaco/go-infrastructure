# Go Infrastructure
[![Go Report Card](https://goreportcard.com/badge/github.com/janmbaco/go-infrastructure)](https://goreportcard.com/report/github.com/janmbaco/go-infrastructure)

This is an infrastructure project in Go that serves the Go-ReverseProxy-SSL and Saprocate projects. It also aims to be a common base for projects in Go.

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Packages](#packages)
  - [Dependency Injection](#dependency-injection)
  - [Logs](#logs)
  - [Error Handler](#error-handler)
  - [Event](#event)
  - [Disk](#disk)
  - [Config](#config)
  - [Server](#server)
  - [Crypto](#crypto)
  - [Persistence](#persistence)
- [Contributing](#contributing)
- [License](#license)

## Installation

To install the package, use `go get`:

```bash
go get github.com/janmbaco/go-infrastructure
```

## Usage

### Basic Example

```go
package main

import (
    "github.com/janmbaco/go-infrastructure/logs"
)

func main() {
    logger := logs.NewLogger()
    logger.Info("This is an informational message")
}
```

## Packages

### Dependency Injection
Provides a container for managing dependencies within an application.

#### Definition
`container.go` defines the interface for a container responsible for managing dependencies.

#### Functions
- `Register() Register`: Gets the object responsible for registering dependencies.
- `Resolver() Resolver`: Gets the object responsible for resolving dependencies.
- `NewContainer() Container`: Returns a new container for managing dependencies.

#### Practical Example
**Example:**

```go
package main

import (
    "github.com/janmbaco/go-infrastructure/dependencyinjection"
)

func main() {
    container := dependencyinjection.NewContainer()
    register := container.Register()
    resolver := container.Resolver()

    register.AsSingleton(new(MyService), func() *MyService {
        return &MyService{}
    }, nil)

    myService := resolver.Get(new(MyService)).(*MyService)
    myService.DoSomething()
}

type MyService struct{}

func (s *MyService) DoSomething() {
    println("Service is doing something.")
}
```

### Logs
Provides a logging service that allows writing to a log file from a logging level (Trace, Info, Warning, Error, Fatal).

#### Example

```go
package main

import (
    "github.com/janmbaco/go-infrastructure/logs"
)

func main() {
    logger := logs.NewLogger()
    logger.Trace("This is a trace message")
    logger.Info("This is an informational message")
    logger.Warning("This is a warning message")
    logger.Error("This is an error message")
    logger.Fatal("This is a fatal message")
}
```

### Error Handler
Provides utilities for handling errors in Go that can be thrown or caught. If the log level is set below fatal, all errors are logged.

#### Definitions and Practical Examples

#### `errorchecker.go`

**Definition:**
`errorchecker.go` provides a function to check if parameters are nil.

**Example:**

```go
package main

import (
    "github.com/janmbaco/go-infrastructure/errors/errorchecker"
)

func main() {
    var param interface{} = nil
    errorchecker.CheckNilParameter(map[string]interface{}{"param": param})
}
```

#### `errorcatcher.go`

**Definition:**
`errorcatcher.go` provides an interface and implementation for catching and handling errors.

**Example:**

```go
package main

import (
    "errors"
    "fmt"
    "github.com/janmbaco/go-infrastructure/errors"
    "github.com/janmbaco/go-infrastructure/logs"
)

func main() {
    logger := logs.NewLogger()
    errorCatcher := errors.NewErrorCatcher(logger)

    // Example of TryCatchError
    errorCatcher.TryCatchError(func() {
        panic("an unexpected error")
    }, func(err error) {
        fmt.Println("Caught an error:", err)
    })

    // Example of CatchError
    err := errors.New("an example error")
    errorCatcher.CatchError(err, func(err error) {
        fmt.Println("Caught an error:", err)
    })
}
```

#### `errormanager.go`

**Definition:**
`errormanager.go` provides an interface and implementation to manage errors and their associated callbacks.

**Example:**

```go
package main

import (
    "errors"
    "fmt"
    "github.com/janmbaco/go-infrastructure/errors"
)

func main() {
    errorManager := errors.NewErrorManager()

    // Register a callback for a specific type of error
    errorManager.On(&MyError{}, func(err error) {
        fmt.Println("Handling MyError:", err)
    })

    // Trigger an error and execute the registered callback
    err := &MyError{Message: "An error occurred"}
    if callback := errorManager.GetCallback(err); callback != nil {
        callback(err)
    }
}

type MyError struct {
    Message string
}

func (e *MyError) Error() string {
    return e.Message
}
```

#### `errorthrower.go`

**Definition:**
`errorthrower.go` provides an interface and implementation to throw errors and execute the corresponding callbacks if registered.

**Example:**

```go
package main

import (
    "errors"
    "fmt"
    "github.com/janmbaco/go-infrastructure/errors"
)

func main() {
    errorManager := errors.NewErrorManager()
    errorThrower := errors.NewErrorThrower(errorManager)

    // Register a callback for a specific type of error
    errorManager.On(&MyError{}, func(err error) {
        fmt.Println("Handling MyError:", err)
    })

    // Throw an error and execute the registered callback
    err := &MyError{Message: "An error occurred"}
    errorThrower.Throw(err)
}

type MyError struct {
    Message string
}

func (e *MyError) Error() string {
    return e.Message
}
```

#### `errordefer.go`

**Definition:**
`errordefer.go` provides an interface and implementation to handle errors that may occur in deferred functions.

**Example:**

```go
package main

import (
    "errors"
    "fmt"
    "github.com/janmbaco/go-infrastructure/errors"
)

func main() {
    errorThrower := errors.NewErrorThrower(nil)
    errorDefer := errors.NewErrorDefer(errorThrower)

    defer errorDefer.TryThrowError(nil)

    // Simulate an error that causes a panic
    panic(errors.New("an unexpected error"))
}
```

#### Combined Practical Example

**Definition:**
This example demonstrates the usage of `errormanager.go`, `errorthrower.go`, and `errordefer.go` together.

**Example:**

```go
package main

import (
    "errors"
    "fmt"
    "github.com/janmbaco/go-infrastructure/errors"
)

type MyError struct {
    Message string
}

func (e *MyError) Error() string {
    return e.Message
}

func main() {
    errorManager := errors.NewErrorManager()
    errorThrower := errors.NewErrorThrower(errorManager)
    errorDefer := errors.NewErrorDefer(errorThrower)

    errorManager.On(&MyError{}, func(err error) {
        fmt.Println("Handling MyError:", err)
    })

    defer errorDefer.TryThrowError(nil)

    // Trigger an error and execute the registered callback
    panic(&MyError{Message: "An example error"})
}
```

### Event
Provides services for subscribing to and publishing events.

#### `eventsmanager/eventobject.go`

**Definition:**
`eventobject.go` defines the interface for an event object responsible for making an event.

#### Functions:
- `GetEventArgs() interface{}`: Returns the event arguments.
- `HasEventArgs() bool`: Checks if the event has arguments.
- `StopPropagation() bool`: Checks if the event propagation should be stopped.
- `IsParallelPropagation() bool`: Checks if the event propagation should be parallel.
- `GetTypeOfFunc() reflect.Type`: Gets the type of the function associated with the event.

**Example:**

```go
package main

import (
    "reflect"
    "github.com/janmbaco/go-infrastructure/eventsmanager"
)

type MyEvent struct {
    args interface{}
}

func (e *MyEvent) GetEventArgs() interface{} {
    return e.args
}

func (e *MyEvent) HasEventArgs() bool {
    return e.args != nil
}

func (e *MyEvent) StopPropagation() bool {
    return false
}

func (e *MyEvent) IsParallelPropagation() bool {
    return true
}

func (e *MyEvent) GetTypeOfFunc() reflect.Type {
    return reflect.TypeOf(func(interface{}) {})
}
```

#### `eventsmanager/publisher.go`

**Definition:**
`publisher.go` defines the interface for a publisher responsible for publishing events.

#### Functions:
- `Publish(event EventObject)`: Publishes an event.

**Example:**

```go
package main

import (
    "reflect"
    "github.com/janmbaco/go-infrastructure/errors"
    "github.com/janmbaco/go-infrastructure/eventsmanager"
)

func main() {
    subscriptions := eventsmanager.NewSubscriptions(errors.NewErrorDefer(nil))
    publisher := eventsmanager.NewPublisher(subscriptions, errors.NewErrorCatcher(nil))

    event := &MyEvent{args: "Event Data"}
    publisher.Publish(event)
}
```

#### `eventsmanager/subscriptions.go`

**Definition:**
`subscriptions.go` defines the interface for managing subscriptions to events.

#### Functions:
- `Add(event EventObject, subscribeFunc interface{})`: Adds a subscription to an event.
- `Remove(event EventObject, subscribeFunc interface{})`: Removes a subscription from an event.
- `GetAlls(event EventObject) []reflect.Value`: Gets all subscriptions for an event.

**Example:**

```go
package main

import (
    "reflect"
    "github.com/janmbaco/go-infrastructure/errors"
    "github.com/janmbaco/go-infrastructure/eventsmanager"
)

func main() {
    subscriptions := eventsmanager.NewSubscriptions(errors.NewErrorDefer(nil))

    subscribeFunc := func(args interface{}) {
        println(args.(string))
    }

    event := &MyEvent{args: "Event Data"}
    subscriptions.Add(event, subscribeFunc)
    subscriptions.Remove(event, subscribeFunc)
}
```

### Disk
Provides tools for writing and deleting files on disk. It also provides a service that listens for changes to a disk file.

#### Example

```go
package main

import (
    "fmt"
    "github.com/janmbaco/go-infrastructure/disk"
)

func main() {
    filePath := "example.txt"

    // Write to a file
    disk.WriteFile(filePath, []byte("Hello, world!"))

    // Read from a file
    data, err := disk.ReadFile(filePath)
    if err != nil {
        fmt.Println("Error reading file:", err)
    } else {
        fmt.Println("File contents:", string(data))
    }

    // Delete a file
    disk.DeleteFile(filePath)
}
```

### Config
Provides a configuration interface and an implementation for a file configuration.

#### Example

```go
package main

import (
    "fmt"
    "github.com/janmbaco/go-infrastructure/config"
)

type Config struct {
    Address string `json:"address"`
}

func main() {
    configHandler := config.NewFileConfigHandler("config.json", &Config{Address: ":8080"})
    config := configHandler.GetConfig().(*Config)
    fmt.Println("Server address:", config.Address)
}
```

### Server
Provides a service that starts an HTTP or gRPC server that automatically restarts when the configuration changes.

#### Example

```go
package main

import (
    "net/http"
    "github.com/janmbaco/go-infrastructure/server"
)

func main() {
    serverHandler := server.NewHTTPServer(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, world!"))
    }))
    serverHandler.Start()
}
```

### Crypto
Provides a service to encrypt and decrypt bytes.

#### `crypto/cipher.go`

**Definition:**
`cipher.go` defines the interface for a cipher responsible for encrypting and decrypting values by a key.

#### Functions:
- `Encrypt(value []byte) []byte`: Encrypts the value.
- `Decrypt(value []byte) []byte`: Decrypts the value.

**Example:**

```go
package main

import (
    "fmt"
    "github.com/janmbaco/go-infrastructure/crypto"
    "github.com/janmbaco/go-infrastructure/errors"
)

func main() {
    key := []byte("example key 1234")
    errorCatcher := errors.NewErrorCatcher(nil)
    errorDefer := errors.NewErrorDefer(nil)
    cipherService := crypto.NewCipher(key, errorCatcher, errorDefer)

    value := []byte("Hello, world!")
    encryptedValue := cipherService.Encrypt(value)
    decryptedValue := cipherService.Decrypt(encryptedValue)

    fmt.Println("Encrypted value:", encryptedValue)
    fmt.Println("Decrypted value:", string(decryptedValue))
}
```

### Persistence
Provides an interface and implementation for data access using GORM.

#### Definitions and Functions

##### `dataaccess.go`
**Definition:**
`dataaccess.go` defines the `DataAccess` interface for CRUD operations.

**Functions:**
- `Insert(datarow interface{})`: Inserts a record into the database.
- `Select(datafilter interface{}, preloads ...string) interface{}`: Selects records from the database.
- `Update(datafilter interface{}, datarow interface{})`: Updates records in the database.
- `Delete(datafilter interface{}, associateds ...string)`: Deletes records from the database.

**Practical Example:**
```go
package main

import (
    "reflect"
    "github.com/janmbaco/go-infrastructure/errors"
    "github.com/janmbaco/go-infrastructure/persistence/orm_base"
    "gorm.io/gorm"
    "fmt"
)

type User struct {
    ID   int
    Name string
}

func main() {
    var db *gorm.DB // Initialize your GORM database here
    errorDefer := errors.NewErrorDefer(nil)
    dataAccess := orm_base.NewDataAccess(errorDefer, db, reflect.TypeOf(&User{}))

    // Insert a user
    user := &User{Name: "John Doe"}
    dataAccess.Insert(user)

    // Select users
    users := dataAccess.Select(&User{Name: "John Doe"})
    fmt.Println(users)
}
```

##### `database_info.go`
**Definition:**
`database_info.go` defines the `DatabaseInfo` struct containing database connection information.

**Practical Example:**
```go
package main

import (
    "github.com/janmbaco/go-infrastructure/persistence/orm_base"
    "fmt"
)

func main() {
    dbInfo := &orm_base.DatabaseInfo{
        Engine:       orm_base.Postgres,
        Host:         "localhost",
        Port:         "5432",
        Name:         "exampledb",
        UserName:     "user",
        UserPassword: "password",
    }
    fmt.Println(dbInfo)
}
```

##### `database_provider.go`
**Definition:**
`database_provider.go` provides the `NewDB` function to create a new database instance using a `DialectorResolver`.

**Practical Example:**
```go
package main

import (
    "github.com/janmbaco/go-infrastructure/persistence/orm_base"
    "gorm.io/gorm"
    "fmt"
)

func main() {
    var dialectorResolver orm_base.DialectorResolver // Initialize your dialector resolver here
    var config *gorm.Config // GORM configuration
    var tables []interface{} // Tables to migrate

    db := orm_base.NewDB(dialectorResolver, &orm_base.DatabaseInfo{
        Engine:       orm_base.Postgres,
        Host:         "localhost",
        Port:         "5432",
        Name:         "exampledb",
        UserName:     "user",
        UserPassword: "password",
    }, config, tables)

    fmt.Println(db)
}
```

##### `dialector_getter.go`
**Definition:**
`dialector_getter.go` defines the `DialectorGetter` interface to obtain a `gorm.Dialector` based on the database information.

##### `dialector_resolver.go`
**Definition:**
`dialector_resolver.go` provides the implementation of `DialectorResolver` to resolve the `gorm.Dialector` using a `dependencyinjection.Resolver`.

**Practical Example:**
```go
package main

import (
    "github.com/janmbaco/go-infrastructure/dependencyinjection"
    "github.com/janmbaco/go-infrastructure/persistence/orm_base"
    "fmt"
)

func main() {
    var resolver dependencyinjection.Resolver // Initialize your resolver here
    dialectorResolver := orm_base.NewDialectorResolver(resolver)

    dbInfo := &orm_base.DatabaseInfo{
        Engine:       orm_base.Postgres,
        Host:         "localhost",
        Port:         "5432",
        Name:         "exampledb",
        UserName:     "user",
        UserPassword: "password",
    }

    dialector := dialectorResolver.Resolve(dbInfo)
    fmt.Println(dialector)
}
```

## Contributing

Contributions are welcome! Please follow these steps to contribute:

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/AmazingFeature`).
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4. Push to the branch (`git push origin feature/AmazingFeature`).
5. Open a pull request.

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
