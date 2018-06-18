package sql

type WhereCondition struct {
	Condition string
	Args      []interface{}
}

func (w *WhereCondition) Clone() *WhereCondition {
	clonedArgsSlice := []interface{}{}
	clonedArgsSlice = append(clonedArgsSlice, w.Args...)
	return &WhereCondition{
		Condition: w.Condition,
		Args:      clonedArgsSlice,
	}
}
