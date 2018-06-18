package sql_statement

import (
	"github.com/thcyron/sqlbuilder"

	"github.com/go-zero-boilerplate/databases"
	"github.com/go-zero-boilerplate/databases/sql"
)

type SelectStatement interface {
	BaseStatement
}

type selectStatement struct {
	dialect databases.Dialect
	db      databases.Database

	TableName   string
	Where       []*sql.WhereCondition
	OrderBy     string
	Limit       int
	Offset      int
	DestColumns []*sql.ColumnNameAndValue
	DestSlice   interface{}
}

func (s *selectStatement) Clone() *selectStatement {
	clonedWheres := []*sql.WhereCondition{}
	if s.Where != nil {
		for _, w := range s.Where {
			clonedWheres = append(clonedWheres, w.Clone())
		}
	}

	clonedDestCols := []*sql.ColumnNameAndValue{}
	if s.DestColumns != nil {
		for _, d := range s.DestColumns {
			clonedDestCols = append(clonedDestCols, d.Clone())
		}
	}

	return &selectStatement{
		dialect:     s.dialect,
		db:          s.db,
		TableName:   s.TableName,
		Where:       clonedWheres,
		OrderBy:     s.OrderBy,
		Limit:       s.Limit,
		Offset:      s.Offset,
		DestColumns: clonedDestCols,
		DestSlice:   s.DestSlice,
	}
}

func (s *selectStatement) buildQuery() (query string, args, dest []interface{}) {
	if len(s.DestColumns) == 0 && s.DestSlice == nil {
		panic("SelectStatement requires either DestColumns or DestSlice")
	}

	builder := sqlbuilder.Select().Dialect(s.dialect).From(s.TableName)

	for _, d := range s.DestColumns {
		builder = builder.Map(d.Name, d.Value)
	}

	if s.OrderBy != "" {
		builder = builder.Order(s.OrderBy)
	}

	if s.Limit > 0 {
		builder = builder.Limit(s.Limit)
	}

	if s.Offset > 0 {
		builder = builder.Offset(s.Offset)
	}

	for _, w := range s.Where {
		builder = builder.Where(w.Condition, w.Args...)
	}

	return builder.Build()
}

func (s *selectStatement) Execute() error {
	query, args, dest := s.buildQuery()

	if s.DestSlice != nil {
		err := s.db.Select(s.DestSlice, query, args...)
		if err != nil {
			return err
		}
	} else if len(s.DestColumns) > 0 {
		err := s.db.QueryRow(query, args...).Scan(dest...)
		if err != nil {
			return err
		}
	}

	return nil
}
