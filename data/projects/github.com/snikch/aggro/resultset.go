package aggro

import (
	"fmt"
	"strings"
)

// Resultset represents a complete set of result buckets and any associated errors.
type Resultset struct {
	Errors      []error         `json:"errors"`
	Buckets     []*ResultBucket `json:"buckets"`
	Composition []interface{}   `json:"-"`
}

// ResultBucket represents recursively built metrics for our tablular data.
type ResultBucket struct {
	Value        string                 `json:"value"`
	Metrics      map[string]interface{} `json:"metrics"`
	Buckets      []*ResultBucket        `json:"buckets"`
	bucketLookup map[string]*ResultBucket
	sourceRows   []map[string]Cell
}

// ResultTable represents a Resultset split into row / columns at a depth.
type ResultTable struct {
	Rows         [][]map[string]interface{} `json:"rows"`
	RowTitles    [][]string                 `json:"row_titles"`
	ColumnTitles [][]string                 `json:"column_titles"`
}

// Concrete errors.
var (
	ErrTargetDepthTooLow     = fmt.Errorf("Tabulate: target depth should be 1 or above")
	ErrTargetDepthNotReached = fmt.Errorf("Tabulate: reached deepest bucket before hitting target depth")
)

// Tabulate takes a Resultset and converts it to tabular data.
func Tabulate(results *Resultset, depth int) (*ResultTable, error) {
	if depth < 1 {
		return nil, ErrTargetDepthTooLow
	}
	// Create our table.
	table := &ResultTable{
		Rows:         [][]map[string]interface{}{},
		RowTitles:    [][]string{},
		ColumnTitles: [][]string{},
	}

	// And a lookup helper instance.
	lookup := &resultLookup{
		cells:        map[string]map[string]interface{}{},
		rowLookup:    map[string]bool{},
		columnLookup: map[string]bool{},
	}

	// Recursively build the lookup for each of the root result buckets.
	for _, bucket := range results.Buckets {
		err := buildLookup([]string{}, 1, depth, table, lookup, bucket)
		if err != nil {
			return nil, err
		}
	}

	// Now build up the cells for each of the row / column tuples.
	for _, row := range table.RowTitles {
		tableRow := []map[string]interface{}{}
		for _, column := range table.ColumnTitles {
			tableRow = append(tableRow, lookup.cells[strings.Join(row, lookupKeyDelimiter)+lookupKeyDelimiter+strings.Join(column, lookupKeyDelimiter)])
		}
		table.Rows = append(table.Rows, tableRow)
	}
	// And we're done 👌.
	return table, nil
}

// resultLookup stores specific data as the result set is recursively iterated over.
type resultLookup struct {
	cells        map[string]map[string]interface{}
	rowLookup    map[string]bool
	columnLookup map[string]bool
}

// lookupKeyDelimiter is used to flatten a string array to a single key.
const lookupKeyDelimiter = "😡"

// buildLookup is a recursive function that breaks data into rows and columns
// at a specific depth.
func buildLookup(key []string, depth, targetDepth int, table *ResultTable, lookup *resultLookup, bucket *ResultBucket) error {
	// Add the new bucket value to the lookup key.
	key = append(key, bucket.Value)

	// If we have no buckets, we're at a metric point.
	if len(bucket.Buckets) == 0 {
		if depth <= targetDepth {
			return ErrTargetDepthNotReached
		}
		// The column key is made up of just the key parts from the target depth.
		columnKey := strings.Join(key[targetDepth:], lookupKeyDelimiter)
		// If we haven't seen this column tuple before, add it to the lookup.
		if _, ok := lookup.columnLookup[columnKey]; !ok {
			table.ColumnTitles = append(table.ColumnTitles, key[targetDepth:])
			lookup.columnLookup[columnKey] = true
		}
		m := bucket.Metrics
		lookup.cells[strings.Join(key, lookupKeyDelimiter)] = m
		return nil
	}

	// If we've reached target depth, add this key to the rows if it's not there.
	if depth == targetDepth {
		rowKey := strings.Join(key, lookupKeyDelimiter)
		if _, ok := lookup.rowLookup[rowKey]; !ok {
			table.RowTitles = append(table.RowTitles, key)
			lookup.rowLookup[rowKey] = true
		}
	}

	// Now continue down the 🐇 hole with the next depth of result buckets.
	for _, bucket := range bucket.Buckets {
		newKey := make([]string, len(key))
		copy(newKey, key)
		err := buildLookup(newKey, depth+1, targetDepth, table, lookup, bucket)
		if err != nil {
			return err
		}
	}
	return nil
}
