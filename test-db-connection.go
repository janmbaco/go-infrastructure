package main

import (
	"fmt"
	"os"

	"github.com/janmbaco/go-infrastructure/v2/persistence"
	"github.com/janmbaco/go-infrastructure/v2/persistence/dialectors"
	"gorm.io/gorm"
)

func main() {
	if len(os.Args) != 7 {
		fmt.Println("Usage: go run test-db-connection.go <engine> <host> <port> <user> <password> <dbname>")
		os.Exit(1)
	}

	engine := os.Args[1]
	host := os.Args[2]
	port := os.Args[3]
	user := os.Args[4]
	password := os.Args[5]
	dbname := os.Args[6]

	var dbEngine persistence.DbEngine
	switch engine {
	case "postgres":
		dbEngine = persistence.Postgres
	case "mysql":
		dbEngine = persistence.MySQL
	case "sqlserver":
		dbEngine = persistence.SQLServer
	default:
		fmt.Printf("Unsupported engine: %s\n", engine)
		os.Exit(1)
	}

	dbInfo := &persistence.DatabaseInfo{
		Host:         host,
		Port:         port,
		Name:         dbname,
		UserName:     user,
		UserPassword: password,
		Engine:       dbEngine,
	}

	var dialectorGetter persistence.DialectorGetter
	switch dbEngine {
	case persistence.Postgres:
		dialectorGetter = dialectors.NewPostgresDialectorGetter()
	case persistence.MySQL:
		dialectorGetter = dialectors.NewMysqlDialectorGetter()
	case persistence.SQLServer:
		dialectorGetter = dialectors.NewSqlServerDialectorGetter()
	default:
		fmt.Printf("Unsupported database engine\n")
		os.Exit(1)
	}

	dialector, err := dialectorGetter.Get(dbInfo)
	if err != nil {
		fmt.Printf("Failed to get dialector: %v\n", err)
		os.Exit(1)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("Failed to get sql.DB: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			fmt.Printf("Warning: failed to close database connection: %v\n", err)
		}
	}()

	if err := sqlDB.Ping(); err != nil {
		fmt.Printf("Failed to ping database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database connection successful")
}
