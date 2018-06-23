package sqlsurgeon;

// AnySQLSource can generate SQL
type AnySQLSource interface {
    SQL() string
}

// AnyStatement is the basis for any SQL statement
type AnyStatement interface {
    AnySQLSource
}

// AnyQuery associates a set of products, the tables that generate the
// products, and the constraints that govern the production.
type AnyQuery interface {
    AnyStatement
    Products() [] AnyProduct
    Tables() [] AnyTable
    Constraints() [] AnyConstraint
}

// AnyProduct defines a label and the associated column that produces
// results.
type AnyProduct interface {
    Label() string
    Expression() AnyExpression
}

// AnyColumn relates a table with one field within.
type AnyColumn interface {
    Table() AnyTable
    Name() string
}

// AnyTable has a name defining the collection.
type AnyTable interface {
    Name() string
    AddColumn(AnyColumn)
    RemoveColumn(AnyColumn)
    ColumnNamed(string) AnyColumn
    Columns() [] AnyColumn
}

//  AnyConstraint can bind a column to some operation. 
type AnyConstraint interface {
    Column() AnyColumn
}

// AnyRelationship compares two columns between tables.
type AnyRelationship interface {
    LeftColumn() AnyColumn
    RightColumn() AnyColumn
}

// AnyExpression manages a node tree of columns, operators, and constatns.
// Most importantly, however, it manages the list of clumns involved in the
// expression.
type AnyExpression interface {
    Columns() [] AnyColumn
}

// ColumnNameMap collects columns by name
type ColumnNameMap map[string]AnyColumn

