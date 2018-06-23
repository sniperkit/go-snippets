package sql_statement

import (
	"github.com/go-zero-boilerplate/databases"
	"github.com/go-zero-boilerplate/databases/sql"
)

type DeleteBuilder interface {
	Where(condition string, args ...interface{}) DeleteBuilder
	RowsAffectedDest(rowsAffectedDest *int64) DeleteBuilder
	Build() DeleteStatement
}

func NewDeleteBuilder(
	dialect databases.Dialect,
	db databases.Database,
	tableName string) DeleteBuilder {

	return &deleteBuilder{
		d: &deleteStatement{
			dialect:   dialect,
			db:        db,
			TableName: tableName,
		},
	}
}

type deleteBuilder struct {
	d *deleteStatement
}

func (d *deleteBuilder) Where(condition string, args ...interface{}) DeleteBuilder {
	d.d.Wheres = append(d.d.Wheres, &sql.WhereCondition{Condition: condition, Args: args})
	return d
}

func (d *deleteBuilder) RowsAffectedDest(rowsAffectedDest *int64) DeleteBuilder {
	d.d.RowsAffectedDest = rowsAffectedDest
	return d
}

func (d *deleteBuilder) Build() DeleteStatement {
	return d.d
}
