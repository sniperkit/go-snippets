package sql_statement

import (
	"github.com/go-zero-boilerplate/databases"
	"github.com/go-zero-boilerplate/databases/sql"
)

type InsertBuilder interface {
	Set(name string, value interface{}) InsertBuilder
	LastInsertIdDest(lastInsertIdDest *int64) InsertBuilder
	Build() InsertStatement
}

func NewInsertBuilder(
	dialect databases.Dialect,
	db databases.Database,
	tableName string) InsertBuilder {

	return &insertBuilder{
		i: &insertStatement{
			dialect:   dialect,
			db:        db,
			TableName: tableName,
		},
	}
}

type insertBuilder struct {
	i *insertStatement
}

func (i *insertBuilder) Set(name string, value interface{}) InsertBuilder {
	i.i.Columns = append(i.i.Columns, &sql.ColumnNameAndValue{Name: name, Value: value})
	return i
}

func (i *insertBuilder) LastInsertIdDest(lastInsertIdDest *int64) InsertBuilder {
	i.i.LastInsertIdDest = lastInsertIdDest
	return i
}

func (i *insertBuilder) Build() InsertStatement {
	return i.i
}
