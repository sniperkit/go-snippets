package sql_statement

import (
	"fmt"

	"github.com/thcyron/sqlbuilder"

	"github.com/go-zero-boilerplate/databases"
	"github.com/go-zero-boilerplate/databases/sql"
)

type DeleteStatement interface {
	BaseStatement
}

type deleteStatement struct {
	dialect databases.Dialect
	db      databases.Database

	TableName string
	Wheres    []*sql.WhereCondition

	RowsAffectedDest *int64
}

func (d *deleteStatement) Execute() error {
	builder := sqlbuilder.Delete().Dialect(d.dialect).From(d.TableName)

	for _, w := range d.Wheres {
		builder = builder.Where(w.Condition, w.Args...)
	}

	query, args := builder.Build()

	result, err := d.db.Exec(query, args...)
	if err != nil {
		return err
	}

	if d.RowsAffectedDest != nil {
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("Unable to get RowsAffected, error: %s", err.Error())
		}
		*d.RowsAffectedDest = rowsAffected
	}

	return nil
}
