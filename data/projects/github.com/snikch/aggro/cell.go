package aggro

import (
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const (
	fieldTypeString   = "string"
	fieldTypeNumber   = "number"
	fieldTypeDatetime = "datetime"
	fieldTypeBoolean  = "boolean"
)

// newCell constructs a Cell{} based on the given *Field.Type. Data being assigned to the *Field
// must be of the same datatype as *Field.Type.
func newCell(data, datum interface{}, field *Field) (Cell, error) {
	var cell Cell

	switch field.Type {
	case fieldTypeString:
		stringValue, ok := datum.(string)
		if !ok {
			return nil, fmt.Errorf("Expected string datatype, got %T", datum)
		}
		cell = &StringCell{
			value: stringValue,
			field: field,
			data:  data,
		}
	case fieldTypeDatetime:
		datetimeCell := &DatetimeCell{
			field: field,
			data:  data,
		}
		switch datumTyped := datum.(type) {
		case time.Time:
			datetimeCell.value = &datumTyped
		case *time.Time:
			if datumTyped == nil {
				return nil, errors.New("Got nil *time.Time for datetime field")
			}
			datetimeCell.value = datumTyped
		case string:
			t, err := time.Parse(time.RFC3339, datumTyped)
			if err != nil {
				return nil, errors.New("Invalid date string passed for datetime field. RFC3339 datetime string required")
			}
			datetimeCell.value = &t
		default:
			return nil, fmt.Errorf("Expected string or time.Time, got %T", datum)
		}
		cell = datetimeCell
	case fieldTypeNumber:
		numberCell := &NumberCell{
			field: field,
			data:  data,
		}
		var d decimal.Decimal
		switch datumTyped := datum.(type) {
		case int:
			d = decimal.NewFromFloat(float64(datumTyped))
		case int32:
			d = decimal.NewFromFloat(float64(datumTyped))
		case int64:
			d = decimal.NewFromFloat(float64(datumTyped))
		case float32:
			d = decimal.NewFromFloat(float64(datumTyped))
		case float64:
			d = decimal.NewFromFloat(datumTyped)
		default:
			return nil, fmt.Errorf("Expected number, got %T", datum)
		}
		numberCell.value = &d
		cell = numberCell
	case fieldTypeBoolean:
		boolValue, ok := datum.(bool)
		if !ok {
			return nil, fmt.Errorf("Expected boolean datatype, got %T", datum)
		}
		booleanCell := &BooleanCell{
			value: boolValue,
			field: field,
			data:  data,
		}
		cell = booleanCell
	default:
		return nil, fmt.Errorf("Unknown field type: %s", field.Type)
	}
	return cell, nil
}

// Cell represents data and configuration for each of our *Table.Fields.
type Cell interface {
	FieldDefinition() *Field
	IsMetricable(m measurer) bool
	MeasurableCell() MeasurableCell
}

// MeasurableCell returns the cells MeasurableCell{}.
type MeasurableCell interface {
	Value() interface{}
}

// NumberCell implements the Cell interface.
type NumberCell struct {
	value *decimal.Decimal
	field *Field
	data  interface{}
}

// FieldDefinition returns the field definition (name) representing the cell.
func (cell *NumberCell) FieldDefinition() *Field {
	return cell.field
}

// IsMetricable determines whether the measurer{} provided can be run by the cell.
func (cell *NumberCell) IsMetricable(m measurer) bool {
	return true
}

// Value returns the cell value.
func (cell *NumberCell) Value() interface{} {
	return cell.value
}

// MeasurableCell returns the cells MeasurableCell{}.
func (cell *NumberCell) MeasurableCell() MeasurableCell {
	return cell
}

// ValueForPeriod returns the start of a given period.
func (cell *NumberCell) ValueForPeriod(period []interface{}) (string, error) {
	value, err := rangeValueForPeriod(cell.value, period)
	return value.String(), err
}

// DatetimeCell implements the Cell interface.
type DatetimeCell struct {
	value *time.Time
	field *Field
	data  interface{}
}

// FieldDefinition returns the field definition (name) representing the cell.
func (cell *DatetimeCell) FieldDefinition() *Field {
	return cell.field
}

// IsMetricable determines whether the measurer{} provided can be run by the cell.
func (cell *DatetimeCell) IsMetricable(m measurer) bool {
	return false
}

// Value returns the cell value.
func (cell *DatetimeCell) Value() interface{} {
	return cell.value
}

// MeasurableCell returns the cells MeasurableCell{}.
func (cell *DatetimeCell) MeasurableCell() MeasurableCell {
	return nil
}

// ValueForPeriod returns the start of a given period.
func (cell *DatetimeCell) ValueForPeriod(period DatetimePeriod, location *time.Location) (string, error) {
	return datetimeValueForPeriod(cell.value, period, location)
}

// StringCell implements the Cell{} interface.
type StringCell struct {
	value string
	field *Field
	data  interface{}
}

// FieldDefinition returns the field definition (name) representing the cell.
func (cell *StringCell) FieldDefinition() *Field {
	return cell.field
}

// IsMetricable determines whether the measurer{} provided can be run by the cell.
func (cell *StringCell) IsMetricable(m measurer) bool {
	// We allow certain metrics to run on StringCells.
	switch m.(type) {
	case *cardinality, *valueCount:
		return true
	default:
		return false
	}
}

// Value returns the cell value.
func (cell *StringCell) Value() interface{} {
	return cell.value
}

// MeasurableCell returns the cells MeasurableCell{}.
func (cell *StringCell) MeasurableCell() MeasurableCell {
	return cell
}

// BooleanCell implements the Cell{} interface.
type BooleanCell struct {
	value bool
	field *Field
	data  interface{}
}

// FieldDefinition returns the field definition (name) representing the cell.
func (cell *BooleanCell) FieldDefinition() *Field {
	return cell.field
}

// IsMetricable determines whether the measurer{} provided can be run by the cell.
func (cell *BooleanCell) IsMetricable(m measurer) bool {
	// We allow certain metrics to run on BooleanCells.
	switch m.(type) {
	case *valueCount:
		return true
	default:
		return false
	}
}

// Value returns the cell value.
func (cell *BooleanCell) Value() interface{} {
	return cell.value
}

// MeasurableCell returns the cells MeasurableCell{}.
func (cell *BooleanCell) MeasurableCell() MeasurableCell {
	return cell
}
