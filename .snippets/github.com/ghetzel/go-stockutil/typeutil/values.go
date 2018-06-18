// Helpers for type inflection and simplifying working with Golang generic interface types
package typeutil

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/ghetzel/go-stockutil/utils"
)

var scs = spew.ConfigState{
	Indent:            `    `,
	DisableCapacities: true,
	SortKeys:          true,
}

// Returns whether the given value represents the underlying type's zero value
func IsZero(value interface{}) bool {
	return utils.IsZero(value)
}

// Returns whether the given value is "empty" in the semantic sense. Zero values
// are considered empty, as are arrays, slices, and maps containing only empty
// values (called recursively). Finally, strings are trimmed of whitespace and
// considered empty if the result is zero-length.
//
func IsEmpty(value interface{}) bool {
	valueV := reflect.ValueOf(value)

	if valueV.Kind() == reflect.Ptr {
		valueV = valueV.Elem()
	}

	// short circuit for zero values of certain types
	switch valueV.Kind() {
	case reflect.Struct:
		if IsZero(value) {
			return true
		}
	}

	switch valueV.Kind() {
	case reflect.Array, reflect.Slice:
		if valueV.Len() == 0 {
			return true
		} else {
			for i := 0; i < valueV.Len(); i++ {
				if indexV := valueV.Index(i); indexV.IsValid() && !IsEmpty(indexV.Interface()) {
					return false
				}
			}

			return true
		}

	case reflect.Map:
		if valueV.Len() == 0 {
			return true
		} else {
			for _, keyV := range valueV.MapKeys() {
				if indexV := valueV.MapIndex(keyV); indexV.IsValid() && !IsEmpty(indexV.Interface()) {
					return false
				}
			}

			return true
		}

	case reflect.Chan:
		if valueV.Len() == 0 {
			return true
		}

	case reflect.String:
		if len(strings.TrimSpace(fmt.Sprintf("%v", value))) == 0 {
			return true
		}
	}

	return false
}

// Return the concrete value pointed to by a pointer type, or within an
// interface type.  Allows functions receiving pointers to supported types
// to work with those types without doing reflection.
//
func ResolveValue(in interface{}) interface{} {
	if inV, ok := in.(Variant); ok {
		in = inV.Value
	}

	return utils.ResolveValue(in)
}

// Dectect whether the concrete underlying value of the given input is one or more
// Kinds of value.
func IsKind(in interface{}, kinds ...reflect.Kind) bool {
	return utils.IsKind(in, kinds...)
}

// Return whether the given input is a discrete scalar value (ints, floats, bools,
// strings), otherwise known as "primitive types" in some other languages.
//
func IsScalar(in interface{}) bool {
	if !IsKind(in, utils.CompoundTypes...) {
		return true
	}

	return false
}

// Returns whether the given value is a slice or array.
func IsArray(in interface{}) bool {
	return IsKind(in, utils.SliceTypes...)
}

// Returns whether the given value is a map.
func IsMap(in interface{}) bool {
	return IsKind(in, reflect.Map)
}

// Returns whether the given value is a struct.
func IsStruct(in interface{}) bool {
	return IsKind(in, reflect.Struct)
}

// Returns whether the given value is a function of any kind
func IsFunction(in interface{}) bool {
	return IsKind(in, reflect.Func)
}

// Returns whether the given value is a function.  If inParams is not -1, the function must
// accept that number of arguments.  If outParams is not -1, the function must return that
// number of values.
func IsFunctionArity(in interface{}, inParams int, outParams int) bool {
	if IsKind(in, reflect.Func) {
		inT := reflect.TypeOf(in)

		if inParams < 0 || inParams >= 0 && inT.NumIn() == inParams {
			if outParams < 0 || outParams >= 0 && inT.NumOut() == outParams {
				return true
			}
		}
	}

	return false
}

// Returns the length of the given value that could have a length (strings, slices, arrays,
// maps, and channels).  If the value is not a type that has a length, -1 is returned.
func Len(in interface{}) int {
	if IsKind(in, reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String) {
		return reflect.ValueOf(in).Len()
	} else {
		return -1
	}
}

// Returns a pretty-printed string representation of the given values.
func Dump(in1 interface{}, in ...interface{}) string {
	return scs.Sdump(append([]interface{}{in1}, in...))
}

// Returns a pretty-printed string representation of the given values.
func Dumpf(format string, in ...interface{}) string {
	return fmt.Sprintf(format, scs.Sdump(in...))
}

// Attempts to set the given reflect.Value to the given interface value
func SetValue(target interface{}, value interface{}) error {
	var targetV, valueV, originalV reflect.Value

	// if we were given a reflect.Value target, then we shouldn't take the reflect.ValueOf that
	if tV, ok := target.(reflect.Value); ok {
		targetV = tV
	} else {
		targetV = reflect.ValueOf(target)

		if targetV.Kind() == reflect.Struct {
			return fmt.Errorf("Must pass a pointer to a struct instance, got %T", target)
		} else if targetV.Kind() == reflect.Ptr {
			// dereference pointers to get at the real destination
			targetV = targetV.Elem()
		}
	}

	// targets must be valid in order to set them to values
	if !targetV.IsValid() {
		return fmt.Errorf("Target %T is not valid", target)
	}

	// if the value we were given was a reflect.Value, just use that
	if vV, ok := value.(reflect.Value); ok {
		originalV = vV
		valueV = vV
	} else {
		originalV = reflect.ValueOf(value)
		valueV = reflect.ValueOf(ResolveValue(value))
	}

	// get the underlying value that was passed in
	if valueV.IsValid() {
		targetT := targetV.Type()
		valueT := valueV.Type()

		// if the target value is a string-a-like, stringify whatever we got in
		if targetT.Kind() == reflect.String && valueV.CanInterface() {
			valueV = reflect.ValueOf(fmt.Sprintf("%v", valueV.Interface()))
			valueT = valueV.Type()

			if !valueV.IsValid() {
				return fmt.Errorf(
					"Converting %T to %v produced an invalid value",
					value,
					targetT,
				)
			}
		}

		if valueT.AssignableTo(targetT) {
			// attempt direct assignment
			targetV.Set(valueV)
		} else if valueT.ConvertibleTo(targetT) {
			// attempt type conversion
			targetV.Set(valueV.Convert(targetT))
		} else if targetV.Kind() == reflect.Ptr {
			if originalV.Kind() == reflect.Ptr {
				return SetValue(targetV, originalV)
			} else {
				return fmt.Errorf(
					"Unable to set target: value for target %v must be given as a pointer",
					targetT,
				)
			}
		} else {
			if targetV.Kind() == reflect.Struct {
				if embeddedV := targetV.FieldByName(valueT.Name()); embeddedV.IsValid() {
					if err := SetValue(embeddedV, value); err == nil {
						return nil
					}
				}
			}

			// no dice.
			return fmt.Errorf(
				"Unable to set target: %T has no path to becoming %v",
				value,
				targetT,
			)
		}
	} else {
		return fmt.Errorf("Unable to set target to the given %T value", value)
	}

	return nil
}

func IsKindOfInteger(in interface{}) bool {
	var kind reflect.Kind

	if k, ok := in.(reflect.Kind); ok {
		kind = k
	} else {
		kind = reflect.TypeOf(in).Kind()
	}

	return IsKind(
		kind,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
	)
}

func IsKindOfFloat(in interface{}) bool {
	var kind reflect.Kind

	if k, ok := in.(reflect.Kind); ok {
		kind = k
	} else {
		kind = reflect.TypeOf(in).Kind()
	}

	return IsKind(
		kind,
		reflect.Float32,
		reflect.Float64,
	)
}

func IsKindOfBool(in interface{}) bool {
	var kind reflect.Kind

	if k, ok := in.(reflect.Kind); ok {
		kind = k
	} else {
		kind = reflect.TypeOf(in).Kind()
	}

	return IsKind(kind, reflect.Bool)
}

func IsKindOfString(in interface{}) bool {
	var kind reflect.Kind

	if k, ok := in.(reflect.Kind); ok {
		kind = k
	} else {
		kind = reflect.TypeOf(in).Kind()
	}

	return IsKind(kind, reflect.String)
}
