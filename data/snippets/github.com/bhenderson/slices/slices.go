// Package slices provides a "generic" Convert function which takes a nil
// function pointer and converts all slice types in the arguments to
// []interface{}
package slices

import "reflect"

// Convert takes a nil function. Panics if fptr arguments don't match fptr
// return values. Exceptoins are that all slice types (including variadic)
// should be converted to []interface{}
func Convert(fptr interface{}) {
	// makeSwap expects fptr to be a pointer to a nil function.
	// It sets that pointer to a new function created with MakeFunc.
	// When the function is invoked, reflect turns the arguments
	// into Values, calls swap, and then turns swap's result slice
	// into the values returned by the new function.

	// fptr is a pointer to a function.
	// Obtain the function value itself (likely nil) as a reflect.Value
	// so that we can query its type and then set the value.
	fn := reflect.ValueOf(fptr).Elem()

	// Make a function of the right type.
	v := reflect.MakeFunc(fn.Type(), swap)

	// Assign it to the value fn represents.
	fn.Set(v)

}

// swap is the implementation passed to MakeFunc.
// It must work in terms of reflect.Values so that it is possible
// to write code without knowing beforehand what the types
// will be.
func swap(in []reflect.Value) []reflect.Value {
	for i := range in {
		v := in[i]
		if v.Kind() == reflect.Slice {
			mid := make([]interface{}, v.Len())
			for j := range mid {
				mid[j] = v.Index(j).Interface()
			}
			in[i] = reflect.ValueOf(mid)
		}
	}
	return in
}
