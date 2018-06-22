package sqlsurgeon;

type Column struct {
    table AnyTable
    name string
}

func NewColumn(table AnyTable, name string) AnyColumn {
    column := &Column{ table, name }
    table.AddColumn(column)
    return column
}

func (c *Column) Drop() {
    c.table.RemoveColumn(c)
}

func (c *Column) Table() AnyTable {
    return c.table
}

func (c *Column) Name() string {
    return c.name
}

