package sql_statement

import (
	"fmt"

	"github.com/thcyron/sqlbuilder"

	"github.com/go-zero-boilerplate/databases"
	"github.com/go-zero-boilerplate/databases/sql"
)

type InsertStatement interface {
	BaseStatement
}

type insertStatement struct {
	dialect databases.Dialect
	db      databases.Database

	TableName string
	Columns   []*sql.ColumnNameAndValue

	LastInsertIdDest *int64
}

func (i *insertStatement) Execute() error {
	builder := sqlbuilder.Insert().Dialect(i.dialect).Into(i.TableName)

	for _, c := range i.Columns {
		builder = builder.Set(c.Name, c.Value)
	}

	query, args, _ := builder.Build()

	result, err := i.db.Exec(query, args...)
	if err != nil {
		return err
	}

	if i.LastInsertIdDest != nil {
		lastInsertId, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("Unable to get LastInsertId, error: %s", err.Error())
		}
		*i.LastInsertIdDest = lastInsertId
	}

	return nil
}
