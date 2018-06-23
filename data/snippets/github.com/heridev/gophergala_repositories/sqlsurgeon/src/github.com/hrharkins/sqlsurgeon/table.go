package sqlsurgeon;

type Table struct {
    name string
    columns ColumnNameMap
}

func NewTable(name string) AnyTable {
    return &Table{ name, make(ColumnNameMap) }
}

func (t *Table) Name() string {
    return t.name
}

func (t *Table) AddColumn(c AnyColumn) {
    t.columns[c.Name()] = c
}

func (t *Table) RemoveColumn(c AnyColumn) {
    delete(t.columns, c.Name())
}

func (t *Table)ColumnNamed(name string) AnyColumn {
    return t.columns[name]
}

func (t *Table)Columns() [] AnyColumn {
    columns := make([]AnyColumn, len(t.columns))
    for _, column := range t.columns {
        columns = append(columns, column)
    }
    return columns
}

