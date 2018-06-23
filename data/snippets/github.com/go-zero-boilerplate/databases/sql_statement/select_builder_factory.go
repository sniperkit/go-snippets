package sql_statement

import (
	"github.com/go-zero-boilerplate/databases"
)

func NewSelectBuilderFactory(db databases.Database) SelectBuilderFactory {
	return &selectBuilderFactory{
		db: db,
	}
}

type SelectBuilderFactory interface {
	FromDialect(dialect databases.Dialect, tableName string) SelectBuilder
	Mysql(tableName string) SelectBuilder
	Sqlite(tableName string) SelectBuilder
	Postgres(tableName string) SelectBuilder
}

type selectBuilderFactory struct {
	db databases.Database
}

func (s *selectBuilderFactory) FromDialect(dialect databases.Dialect, tableName string) SelectBuilder {
	return NewSelectBuilder(dialect, s.db, tableName)
}

func (s *selectBuilderFactory) Mysql(tableName string) SelectBuilder {
	return s.FromDialect(databases.MysqlDialect, tableName)
}

func (s *selectBuilderFactory) Sqlite(tableName string) SelectBuilder {
	return s.FromDialect(databases.SqliteDialect, tableName)
}

func (s *selectBuilderFactory) Postgres(tableName string) SelectBuilder {
	return s.FromDialect(databases.PostgresDialect, tableName)
}
