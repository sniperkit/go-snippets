package unittest

import (
	"reflect"
)

func isNil(v interface{}) bool {
	if v == nil {
		return true
	}

	value := reflect.ValueOf(v)
	kind := value.Kind()
	nilable := kind == reflect.Slice || kind == reflect.Chan || kind == reflect.Func || kind == reflect.Ptr || kind == reflect.Map
	return nilable && value.IsNil()
}
