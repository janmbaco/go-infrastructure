package persistence //nolint:revive // established package name, changing would break API

import "fmt"

type DbEngine uint8

const (
	_SQLServerDB string   = "SqlServerDB"
	_PostgresDB  string   = "PostgresDB"
	_MySQLDB     string   = "MySqlDB"
	_SqliteDB    string   = "SqliteDB"
	SQLServer    DbEngine = iota
	Postgres
	MySQL
	Sqlite
)

type DatabaseInfo struct {
	Host         string   `json:"host"`
	Port         string   `json:"port"`
	Name         string   `json:"name"`
	UserName     string   `json:"user_name"`
	UserPassword string   `json:"user_password"`
	Engine       DbEngine `json:"engine"`
}

func (engine DbEngine) ToString() (string, error) {
	switch engine {
	case SQLServer:
		return _SQLServerDB, nil
	case Postgres:
		return _PostgresDB, nil
	case MySQL:
		return _MySQLDB, nil
	case Sqlite:
		return _SqliteDB, nil
	}
	return "", fmt.Errorf("unknown database engine: %d", engine)
}
