package sql_statement

import (
	"fmt"

	"github.com/thcyron/sqlbuilder"

	"github.com/go-zero-boilerplate/databases"
	"github.com/go-zero-boilerplate/databases/sql"
)

type UpdateStatement interface {
	BaseStatement
}

type updateStatement struct {
	dialect databases.Dialect
	db      databases.Database

	TableName string
	Wheres    []*sql.WhereCondition
	Sets      []*sql.ColumnNameAndValue

	RowsAffectedDest *int64
}

func (u *updateStatement) Execute() error {
	builder := sqlbuilder.Update().Dialect(u.dialect).Table(u.TableName)

	for _, c := range u.Sets {
		builder = builder.Set(c.Name, c.Value)
	}

	for _, w := range u.Wheres {
		builder = builder.Where(w.Condition, w.Args...)
	}

	query, args := builder.Build()

	result, err := u.db.Exec(query, args...)
	if err != nil {
		return err
	}

	if u.RowsAffectedDest != nil {
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("Unable to get RowsAffected, error: %s", err.Error())
		}
		*u.RowsAffectedDest = rowsAffected
	}

	return nil
}
