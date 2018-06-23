package sql_statement

import (
	"github.com/go-zero-boilerplate/databases"
	"github.com/go-zero-boilerplate/databases/sql"
)

type UpdateBuilder interface {
	Where(condition string, args ...interface{}) UpdateBuilder
	Set(name string, value interface{}) UpdateBuilder
	RowsAffectedDest(rowsAffectedDest *int64) UpdateBuilder
	Build() UpdateStatement
}

func NewUpdateBuilder(
	dialect databases.Dialect,
	db databases.Database,
	tableName string) UpdateBuilder {

	return &updateBuilder{
		u: &updateStatement{
			dialect:   dialect,
			db:        db,
			TableName: tableName,
		},
	}
}

type updateBuilder struct {
	u *updateStatement
}

func (u *updateBuilder) Where(condition string, args ...interface{}) UpdateBuilder {
	u.u.Wheres = append(u.u.Wheres, &sql.WhereCondition{Condition: condition, Args: args})
	return u
}

func (u *updateBuilder) Set(name string, value interface{}) UpdateBuilder {
	u.u.Sets = append(u.u.Sets, &sql.ColumnNameAndValue{Name: name, Value: value})
	return u
}

func (u *updateBuilder) RowsAffectedDest(rowsAffectedDest *int64) UpdateBuilder {
	u.u.RowsAffectedDest = rowsAffectedDest
	return u
}

func (u *updateBuilder) Build() UpdateStatement {
	return u.u
}
