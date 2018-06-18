package util

import "reflect"

// Collection of common util functions
//
// Made by Vegard Løkken <vegard@loekken.org>
// License MIT

// InSliceInterface checks wether or not a value is present in a slice of
// any type. Uses reflection which can be slow in some cases.
// It returns a boolean value of true or false
func InSliceInterface(search interface{}, haystack interface{}) bool {
	h := interfaceSlice(haystack)

	for _, v := range h {
		if v == search {
			return true
		}
	}
	return false
}

func interfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
