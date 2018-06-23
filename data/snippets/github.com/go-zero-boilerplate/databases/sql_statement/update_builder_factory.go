package sql_statement

import (
	"github.com/go-zero-boilerplate/databases"
)

func NewUpdateBuilderFactory(db databases.Database) UpdateBuilderFactory {
	return &updateBuilderFactory{
		db: db,
	}
}

type UpdateBuilderFactory interface {
	FromDialect(dialect databases.Dialect, tableName string) UpdateBuilder
	Mysql(tableName string) UpdateBuilder
	Sqlite(tableName string) UpdateBuilder
	Postgres(tableName string) UpdateBuilder
}

type updateBuilderFactory struct {
	db databases.Database
}

func (u *updateBuilderFactory) FromDialect(dialect databases.Dialect, tableName string) UpdateBuilder {
	return NewUpdateBuilder(dialect, u.db, tableName)
}

func (u *updateBuilderFactory) Mysql(tableName string) UpdateBuilder {
	return u.FromDialect(databases.MysqlDialect, tableName)
}

func (u *updateBuilderFactory) Sqlite(tableName string) UpdateBuilder {
	return u.FromDialect(databases.SqliteDialect, tableName)
}

func (u *updateBuilderFactory) Postgres(tableName string) UpdateBuilder {
	return u.FromDialect(databases.PostgresDialect, tableName)
}
