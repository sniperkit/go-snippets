package sql_statement

import (
	"github.com/go-zero-boilerplate/databases"
	"github.com/go-zero-boilerplate/databases/sql"
)

type SelectBuilder interface {
	Where(condition string, args ...interface{}) SelectBuilder
	OrderBy(orderBy string) SelectBuilder
	Limit(limit int) SelectBuilder
	Offset(offset int) SelectBuilder
	Map(colName string, destValue interface{}) SelectBuilder
	MapCountStar(destValue *int) SelectBuilder
	MapSlice(destSlice interface{}) SelectBuilder

	ClearMappings() SelectBuilder
	ClearLimit() SelectBuilder
	ClearOffset() SelectBuilder

	Clone() SelectBuilder

	Build() SelectStatement
}

func NewSelectBuilder(
	dialect databases.Dialect,
	db databases.Database,
	tableName string) SelectBuilder {

	return &selectBuilder{
		s: &selectStatement{
			dialect:   dialect,
			db:        db,
			TableName: tableName,
		},
	}
}

type selectBuilder struct {
	s *selectStatement
}

func (s *selectBuilder) Where(condition string, args ...interface{}) SelectBuilder {
	s.s.Where = append(s.s.Where, &sql.WhereCondition{Condition: condition, Args: args})
	return s
}

func (s *selectBuilder) OrderBy(orderBy string) SelectBuilder {
	s.s.OrderBy = orderBy
	return s
}

func (s *selectBuilder) Limit(limit int) SelectBuilder {
	s.s.Limit = limit
	return s
}

func (s *selectBuilder) Offset(offset int) SelectBuilder {
	s.s.Offset = offset
	return s
}

func (s *selectBuilder) Map(colName string, destValue interface{}) SelectBuilder {
	s.s.DestSlice = nil
	s.s.DestColumns = append(s.s.DestColumns, &sql.ColumnNameAndValue{Name: colName, Value: destValue})
	return s
}

func (s *selectBuilder) MapCountStar(destValue *int) SelectBuilder {
	//TODO: This syntax might differ based on the Dialect?
	return s.Map("COUNT(*)", destValue)
}

func (s *selectBuilder) MapSlice(destSlice interface{}) SelectBuilder {
	s.s.DestColumns = nil
	s.s.DestSlice = destSlice
	s.s.DestColumns = append(s.s.DestColumns, &sql.ColumnNameAndValue{Name: "*", Value: nil})
	return s
}

func (s *selectBuilder) ClearMappings() SelectBuilder {
	s.s.DestColumns = nil
	s.s.DestSlice = nil
	return s
}

func (s *selectBuilder) ClearLimit() SelectBuilder {
	s.s.Limit = 0
	return s
}

func (s *selectBuilder) ClearOffset() SelectBuilder {
	s.s.Offset = 0
	return s
}

func (s *selectBuilder) Clone() SelectBuilder {
	return &selectBuilder{
		s: s.s.Clone(),
	}
}

func (s *selectBuilder) Build() SelectStatement {
	return s.s
}
