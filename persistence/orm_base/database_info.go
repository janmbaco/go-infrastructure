package orm_base

import "fmt"

type DbEngine uint8

const (
	_SqlServerDB string   = "SqlServerDB"
	_PostgresDB  string   = "PostgresDB"
	_MySqlDB     string   = "MySqlDB"
	_SqliteDB    string   = "SqliteDB"
	SqlServer    DbEngine = iota
	Postgres
	MySql
	Sqlite
)

type DatabaseInfo struct {
	Engine       DbEngine `json:"engine"`
	Host         string   `json:"host"`
	Port         string   `json:"port"`
	Name         string   `json:"name"`
	UserName     string   `json:"user_name"`
	UserPassword string   `json:"user_password"`
}

func (engine DbEngine) ToString() (string, error) {
	switch engine {
	case SqlServer:
		return _SqlServerDB, nil
	case Postgres:
		return _PostgresDB, nil
	case MySql:
		return _MySqlDB, nil
	case Sqlite:
		return _SqliteDB, nil
	}
	return "", fmt.Errorf("unknown database engine: %d", engine)
}
