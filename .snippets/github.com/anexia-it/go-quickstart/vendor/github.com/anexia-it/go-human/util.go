package human

import "reflect"

// IsNilOrEmpty checks if a passed interface is either nil or of
// the type's zero value.
//
// The zero value depends on the passed type. For example, the zero value of a string is an empty string.
func IsNilOrEmpty(i interface{}, v reflect.Value) bool {
	// Simple case: interface is nil
	if i == nil {
		return true
	}

	// Hard case: check if interface has "zero" value (ie. empty string, zero integer, etc.)
	return reflect.DeepEqual(i, reflect.Zero(v.Type()).Interface())
}
