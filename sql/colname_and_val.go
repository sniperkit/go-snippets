package sql

type ColumnNameAndValue struct {
	Name  string
	Value interface{}
}

func (c *ColumnNameAndValue) Clone() *ColumnNameAndValue {
	return &ColumnNameAndValue{
		Name:  c.Name,
		Value: c.Value,
	}
}
