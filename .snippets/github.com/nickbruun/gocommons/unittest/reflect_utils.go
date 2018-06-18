package unittest

import (
	"reflect"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func funcTakesSelfReturns0(fun reflect.Value) bool {
	funT := fun.Type()
	return funT.NumIn() == 1 && funT.NumOut() == 0
}
