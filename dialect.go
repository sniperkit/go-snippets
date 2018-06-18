package databases

import (
	"github.com/thcyron/sqlbuilder"
	"github.com/thcyron/sqlbuilder/mysql"
	"github.com/thcyron/sqlbuilder/postgres"
)

type Dialect interface {
	sqlbuilder.Dialect
}

var (
	MysqlDialect    Dialect = &mysql.Dialect{}
	SqliteDialect   Dialect = &mysql.Dialect{} //Currently same as mysql
	PostgresDialect Dialect = &postgres.Dialect{}
)
