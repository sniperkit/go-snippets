package merge

import (
	"database/sql"
	"errors"
	"reflect"
)

var (
	ErrorMergePointers   error = errors.New("merge: must be pointer")
	ErrorMergeSameStruct error = errors.New("merge: must be same struct")
)

func MergeStruct(l interface{}, r interface{}) error {
	lVal := reflect.ValueOf(l)
	rVal := reflect.ValueOf(r)

	if lVal.Kind() != reflect.Ptr || rVal.Kind() != reflect.Ptr {
		return ErrorMergePointers
	}

	if lVal.Type() != rVal.Type() {
		return ErrorMergeSameStruct
	}

	lEleVal := lVal.Elem()
	rEleVal := rVal.Elem()

	for i := 0; i < lEleVal.NumField(); i++ {
		lFieldVal := lEleVal.Field(i)
		rFieldVal := rEleVal.Field(i)

		if !lFieldVal.CanSet() || !rFieldVal.IsValid() || (rFieldVal.Kind() == reflect.Ptr && rFieldVal.IsNil()) {
			// is nil skip it
			continue
		}
		lFieldVal.Set(rFieldVal)
	}

	return nil
}

func indirect(value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value
}

func indirectType(typ reflect.Type) reflect.Type {
	for typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}
	return typ
}

func set(to, from reflect.Value) bool {
	if from.IsValid() {
		if to.Kind() == reflect.Ptr {
			if to.IsNil() {
				to.Set(reflect.New(to.Type().Elem()))
			}
			to = to.Elem()
		}
		if from.Type().ConvertibleTo(to.Type()) {
			to.Set(from.Convert(to.Type()))
		} else if scanner, ok := to.Addr().Interface().(sql.Scanner); ok {
			scanner.Scan(from.Interface())
		} else if from.Kind() == reflect.Ptr {
			return set(to, from.Elem())
		} else {
			return false
		}
	}
	return true
}
