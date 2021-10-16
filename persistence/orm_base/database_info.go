package orm_base

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

func (engine DbEngine) ToString() string {
	switch engine {
	case SqlServer:
		return _SqlServerDB
	case Postgres:
		return _PostgresDB
	case MySql:
		return _MySqlDB
	case Sqlite:
		return _SqliteDB
	}
	panic("not found")
}
