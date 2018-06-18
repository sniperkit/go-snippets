package inputs

import (
	"fmt"
	"strconv"
	"strings"
)

//NewValue creates a new instance of Value
func NewValue(typeName, name, value string) Value {
	return Value{
		TypeName: typeName,
		Name:     name,
		Value:    value,
	}
}

//NewValueError creates a value with just an error
func NewValueError(err error, typeName, name string) Value {
	return Value{
		Error:    err,
		TypeName: typeName,
		Name:     name,
	}
}

//Value contains a Name and Value used for basic operations
type Value struct {
	Error    error
	TypeName string
	Name     string
	Value    string
}

func (v Value) isValueSpecified() bool {
	return strings.TrimSpace(v.Value) != ""
}

func (v Value) confirmNonEmptyString() error {
	if !v.isValueSpecified() {
		return fmt.Errorf("Input %s '%s' is required", v.TypeName, v.Name)
	}
	return nil
}

//RequiredString will check that the string is not empty
func (v Value) RequiredString() (string, error) {
	if v.Error != nil {
		return "", v.Error
	}
	if err := v.confirmNonEmptyString(); err != nil {
		return "", err
	}
	return v.Value, nil
}

//OptionalString will check that the string is not empty
func (v Value) OptionalString() (string, bool, error) {
	if v.Error != nil {
		return "", false, v.Error
	}
	if !v.isValueSpecified() {
		return "", false, nil
	}
	return v.Value, true, nil
}

//RequiredInt64 will parse it as an Int64
func (v Value) RequiredInt64() (int64, error) {
	if v.Error != nil {
		return 0, v.Error
	}
	if err := v.confirmNonEmptyString(); err != nil {
		return 0, err
	}

	intVal, err := strconv.ParseInt(string(v.Value), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Cannot parse input %s '%s' as Int64, error: %s", v.TypeName, v.Name, err.Error())
	}
	return intVal, nil
}

//OptionalInt64 will parse it as an Int64
func (v Value) OptionalInt64() (int64, bool, error) {
	if v.Error != nil {
		return 0, false, v.Error
	}
	if !v.isValueSpecified() {
		return 0, false, nil
	}

	intVal, err := strconv.ParseInt(string(v.Value), 10, 64)
	if err != nil {
		return 0, true, fmt.Errorf("Cannot parse input %s '%s' as Int64, error: %s", v.TypeName, v.Name, err.Error())
	}
	return intVal, true, nil
}
