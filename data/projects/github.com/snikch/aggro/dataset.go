package aggro

import (
	"fmt"
)

// Dataset holds our *Table representation of *Fields and its matching Cell data.
type Dataset struct {
	Table *Table
	Rows  []map[string]Cell
}

// Run executes the query against the dataset.
func (set *Dataset) Run(query *Query) (*Resultset, error) {
	return (&queryProcessor{
		dataset: set,
		query:   query,
	}).Run()
}

// AddRows creates a Cell{} for each of our Table.Fields and ensures the cells data
// meets the cells defined format.
func (set *Dataset) AddRows(rows ...map[string]interface{}) error {
	var err error
	// Add each row.
	for i, data := range rows {
		row := map[string]Cell{}
		// For each row, we need to create a cell for each field definition and
		// ensure that we have received data that conforms to the supposed format.
		for j, field := range set.Table.Fields {
			datum, ok := data[field.Name]
			if !ok {
				return fmt.Errorf("Error adding row %d, cell %d: Data key %s not present", i, j, field.Name)
			}
			if datum == nil {
				// Skip if datum is nil
				continue
			}

			row[field.Name], err = newCell(data, datum, &field)
			if err != nil {
				return fmt.Errorf("Error adding row %d, cell %d: %s", i, j, err.Error())
			}
		}
		set.Rows = append(set.Rows, row)
	}
	return nil
}
