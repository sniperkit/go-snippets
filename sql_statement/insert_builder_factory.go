package sql_statement

import (
	"github.com/go-zero-boilerplate/databases"
)

func NewInsertBuilderFactory(db databases.Database) InsertBuilderFactory {
	return &insertBuilderFactory{
		db: db,
	}
}

type InsertBuilderFactory interface {
	FromDialect(dialect databases.Dialect, tableName string) InsertBuilder
	Mysql(tableName string) InsertBuilder
	Sqlite(tableName string) InsertBuilder
	Postgres(tableName string) InsertBuilder
}

type insertBuilderFactory struct {
	db databases.Database
}

func (i *insertBuilderFactory) FromDialect(dialect databases.Dialect, tableName string) InsertBuilder {
	return NewInsertBuilder(dialect, i.db, tableName)
}

func (i *insertBuilderFactory) Mysql(tableName string) InsertBuilder {
	return i.FromDialect(databases.MysqlDialect, tableName)
}

func (i *insertBuilderFactory) Sqlite(tableName string) InsertBuilder {
	return i.FromDialect(databases.SqliteDialect, tableName)
}

func (i *insertBuilderFactory) Postgres(tableName string) InsertBuilder {
	return i.FromDialect(databases.PostgresDialect, tableName)
}
