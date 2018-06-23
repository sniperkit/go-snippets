package sql_statement

import (
	"github.com/go-zero-boilerplate/databases"
)

func NewDeleteBuilderFactory(db databases.Database) DeleteBuilderFactory {
	return &deleteBuilderFactory{
		db: db,
	}
}

type DeleteBuilderFactory interface {
	FromDialect(dialect databases.Dialect, tableName string) DeleteBuilder
	Mysql(tableName string) DeleteBuilder
	Sqlite(tableName string) DeleteBuilder
	Postgres(tableName string) DeleteBuilder
}

type deleteBuilderFactory struct {
	db databases.Database
}

func (d *deleteBuilderFactory) FromDialect(dialect databases.Dialect, tableName string) DeleteBuilder {
	return NewDeleteBuilder(dialect, d.db, tableName)
}

func (d *deleteBuilderFactory) Mysql(tableName string) DeleteBuilder {
	return d.FromDialect(databases.MysqlDialect, tableName)
}

func (d *deleteBuilderFactory) Sqlite(tableName string) DeleteBuilder {
	return d.FromDialect(databases.SqliteDialect, tableName)
}

func (d *deleteBuilderFactory) Postgres(tableName string) DeleteBuilder {
	return d.FromDialect(databases.PostgresDialect, tableName)
}
