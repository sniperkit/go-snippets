// Utilities for converting, manipulating, and iterating over maps
package maputil

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/ghetzel/go-stockutil/sliceutil"
	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/ghetzel/go-stockutil/utils"
)

var UnmarshalStructTag string = `maputil`
var SkipDescendants = errors.New("skip descendants")

type WalkFunc func(value interface{}, path []string, isLeaf bool) error
type ApplyFunc func(key []string, value interface{}) (interface{}, bool)

func Keys(input interface{}) []interface{} {
	keys := make([]interface{}, 0)
	input = typeutil.ResolveValue(input)

	if input == nil {
		return keys
	}

	inputV := reflect.ValueOf(input)

	if inputV.Kind() == reflect.Map {
		keysV := inputV.MapKeys()

		for _, keyV := range keysV {
			keys = append(keys, keyV)
		}
	} else if syncMap, ok := input.(sync.Map); ok {
		syncMap.Range(func(key interface{}, _ interface{}) bool {
			keys = append(keys, key)
			return true
		})
	}

	return keys
}

func StringKeys(input interface{}) []string {
	keys := sliceutil.Stringify(Keys(input))
	sort.Strings(keys)

	return keys
}

func MapValues(input interface{}) []interface{} {
	values := make([]interface{}, 0)

	inputV := reflect.ValueOf(input)

	switch inputV.Kind() {
	case reflect.Map:
		for _, mapKeyV := range inputV.MapKeys() {
			if mapV := inputV.MapIndex(mapKeyV); mapV.IsValid() && mapV.CanInterface() {
				values = append(values, mapV.Interface())
			}
		}
	}

	return values
}

func TaggedStructFromMap(input map[string]interface{}, populate interface{}, tagname string) error {
	var populateV reflect.Value

	if pV, ok := populate.(reflect.Value); ok {
		populateV = pV
	} else {
		populateV = reflect.ValueOf(populate)
	}

	if ptrKind := populateV.Kind(); ptrKind != reflect.Ptr {
		return fmt.Errorf("Output value must be a pointer to a struct instance, got: %s", ptrKind.String())
	} else {
		if elem := populateV.Elem(); elem.Kind() != reflect.Struct {
			return fmt.Errorf("Value must point to a struct instance, got: %T (%s)", populate, elem.Kind().String())
		} else {
			elemType := elem.Type()

			for i := 0; i < elemType.NumField(); i++ {
				field := elemType.Field(i)
				fieldName := field.Name

				// do this so that we'll always consider the "maputil:" tag first
				tagValue := field.Tag.Get(UnmarshalStructTag)

				// no maputil tag, fallback to whatever we were given
				if tagValue == `` && tagname != `` {
					tagValue = field.Tag.Get(tagname)
				}

				// if we found a tag, parse it
				if tagValue != `` {
					tagParts := strings.Split(tagValue, `,`)
					fieldName = tagParts[0]
				}

				if v, ok := input[fieldName]; ok {
					fieldValue := elem.Field(i)

					if fieldValue.CanSet() {
						vValue := reflect.ValueOf(v)

						if vValue.IsValid() {
							// this is where we handle nested structs being populated with nested maps
							switch v.(type) {
							case map[string]interface{}:
								vMap := v.(map[string]interface{})

								// see if we can directly convert/assign the values
								if vValue.Type().ConvertibleTo(fieldValue.Type()) {
									fieldValue.Set(vValue.Convert(fieldValue.Type()))
									continue
								}

								// recursively populate a new instance of whatever type the destination field is
								// using this input map value
								if convertedValue, err := populateNewInstanceFromMap(vMap, fieldValue.Type()); err == nil {
									fieldValue.Set(convertedValue)
								} else {
									return err
								}

							case []interface{}:
								switch fieldValue.Kind() {
								case reflect.Array, reflect.Slice:
									vISlice := v.([]interface{})

									// for each element of the input array...
									for _, value := range vISlice {
										vIValue := reflect.ValueOf(value)

										switch value.(type) {
										case map[string]interface{}:
											nestedMap := value.(map[string]interface{})

											// recursively populate a new instance of whatever type the destination field is
											// using this input array element
											if convertedValue, err := populateNewInstanceFromMap(nestedMap, fieldValue.Type().Elem()); err == nil {
												vIValue = convertedValue
											} else {
												return err
											}
										}

										// make sure the types are compatible and append the new value to the output field
										if vIValue.Type().ConvertibleTo(fieldValue.Type().Elem()) {
											fieldValue.Set(reflect.Append(fieldValue, vIValue.Convert(fieldValue.Type().Elem())))
										}
									}
								}

							case []map[string]interface{}:
								switch fieldValue.Kind() {
								case reflect.Array, reflect.Slice:
									vISlice := v.([]map[string]interface{})

									// for each nested map element of the input array...
									for _, nestedMap := range vISlice {
										// recursively populate a new instance of whatever type the destination field is
										// using this input array element
										if convertedValue, err := populateNewInstanceFromMap(nestedMap, fieldValue.Type().Elem()); err == nil {
											// make sure the types are compatible and append the new value to the output field
											if convertedValue.Type().ConvertibleTo(fieldValue.Type().Elem()) {
												fieldValue.Set(reflect.Append(fieldValue, convertedValue.Convert(fieldValue.Type().Elem())))
											}
										} else {
											return err
										}
									}
								}
							default:
								// if not special cases from above were encountered, attempt a type conversion
								// and fill in the data
								if vValue.Type().ConvertibleTo(fieldValue.Type()) {
									fieldValue.Set(vValue.Convert(fieldValue.Type()))
								}
							}
						}
					} else {
						return fmt.Errorf("Field '%s' value cannot be changed", field.Name)
					}
				}
			}

		}
	}

	return nil
}

func StructFromMap(input map[string]interface{}, populate interface{}) error {
	return TaggedStructFromMap(input, populate, ``)
}

func populateNewInstanceFromMap(input map[string]interface{}, destination reflect.Type) (reflect.Value, error) {
	var newFieldInstance reflect.Value

	if destination.Kind() == reflect.Struct {
		// get a new instance of the type we want to populate
		newFieldInstance = reflect.New(destination)
	} else if destination.Kind() == reflect.Ptr && destination.Elem().Kind() == reflect.Struct {
		// get a new instance of the type this pointer is pointing at
		newFieldInstance = reflect.New(destination.Elem())
	}

	if newFieldInstance.IsValid() {
		// recursively call StructFromMap, passing the current map[s*]i* value and the new
		// instance we just created.
		//
		if err := StructFromMap(input, newFieldInstance.Interface()); err == nil {
			if newFieldInstance.Elem().Type().ConvertibleTo(destination) {
				// handle as-value
				return newFieldInstance.Elem().Convert(destination), nil
			} else if newFieldInstance.Type().ConvertibleTo(destination) {
				// handle as-ptr
				return newFieldInstance.Convert(destination), nil
			}
		} else {
			return reflect.ValueOf(nil), err
		}
	}

	return reflect.ValueOf(nil), fmt.Errorf("Could not instantiate type %v", destination)
}

func Join(input map[string]interface{}, innerJoiner string, outerJoiner string) string {
	parts := make([]string, 0)

	for key, value := range input {
		if v, err := stringutil.ToString(value); err == nil {
			parts = append(parts, key+innerJoiner+v)
		}
	}

	return strings.Join(parts, outerJoiner)
}

func Split(input string, innerJoiner string, outerJoiner string) map[string]interface{} {
	rv := make(map[string]interface{})
	pairs := strings.Split(input, outerJoiner)

	for _, pair := range pairs {
		kv := strings.SplitN(pair, innerJoiner, 2)

		if len(kv) == 2 {
			rv[kv[0]] = kv[1]
		}
	}

	return rv
}

// Take a flat (non-nested) map keyed with fields joined on fieldJoiner and return a
// deeply-nested map
//
func DiffuseMap(data map[string]interface{}, fieldJoiner string) (map[string]interface{}, error) {
	rv, _ := DiffuseMapTyped(data, fieldJoiner, "")
	return rv, nil
}

// Take a flat (non-nested) map keyed with fields joined on fieldJoiner and return a
// deeply-nested map
//
func DiffuseMapTyped(data map[string]interface{}, fieldJoiner string, typePrefixSeparator string) (map[string]interface{}, []error) {
	errs := make([]error, 0)
	output := make(map[string]interface{})

	//  get the list of keys and sort them because order in a map is undefined
	dataKeys := StringKeys(data)
	sort.Strings(dataKeys)

	//  for each data item
	for _, key := range dataKeys {
		var keyParts []string

		value, _ := data[key]

		//  handle "typed" maps in which the type information is embedded
		if typePrefixSeparator != "" {
			var typeName string

			typeName, key = stringutil.SplitPairTrailing(key, typePrefixSeparator)

			if typeName == `` {
				typeName = `str`
			}

			if v, err := coerceIntoType(value, typeName); err == nil {
				value = v
			} else {
				errs = append(errs, err)
			}
		}

		keyParts = strings.Split(key, fieldJoiner)

		output = DeepSet(output, keyParts, value).(map[string]interface{})
	}

	return output, errs
}

// Take a deeply-nested map and return a flat (non-nested) map with keys whose intermediate tiers are joined with fieldJoiner
//
func CoalesceMap(data map[string]interface{}, fieldJoiner string) (map[string]interface{}, error) {
	return deepGetValues([]string{}, fieldJoiner, data), nil
}

// Take a deeply-nested map and return a flat (non-nested) map with keys whose intermediate tiers are joined with fieldJoiner
// Additionally, values will be converted to strings and keys will be prefixed with the datatype of the value
//
func CoalesceMapTyped(data map[string]interface{}, fieldJoiner string, typePrefixSeparator string) (map[string]interface{}, []error) {
	errs := make([]error, 0)
	rv := make(map[string]interface{})

	for k, v := range deepGetValues([]string{}, fieldJoiner, data) {
		if stringVal, err := stringutil.ToString(v); err == nil {
			rv[prepareCoalescedKey(k, v, typePrefixSeparator)] = stringVal
		} else {
			errs = append(errs, err)
		}
	}

	return rv, errs
}

func deepGetValues(keys []string, joiner string, data interface{}) map[string]interface{} {
	rv := make(map[string]interface{})
	data = typeutil.ResolveValue(data)

	if data != nil {
		dType := reflect.TypeOf(data)

		switch dType.Kind() {
		case reflect.Map:
			for k, v := range data.(map[string]interface{}) {
				newKey := keys
				newKey = append(newKey, k)

				for kk, vv := range deepGetValues(newKey, joiner, v) {
					rv[kk] = vv
				}
			}

		case reflect.Slice, reflect.Array:
			for i, value := range sliceutil.Sliceify(data) {
				newKey := keys
				newKey = append(newKey, strconv.Itoa(i))

				for k, v := range deepGetValues(newKey, joiner, value) {
					rv[k] = v
				}
			}

		default:
			rv[strings.Join(keys, joiner)] = data
		}
	}

	return rv
}

func prepareCoalescedKey(key string, value interface{}, typePrefixSeparator string) string {
	if typePrefixSeparator == "" {
		return key
	} else {
		var datatype string

		if dtype := utils.DetectConvertType(value); dtype != utils.Invalid {
			datatype = dtype.String()
		}

		return datatype + typePrefixSeparator + key
	}
}

func coerceIntoType(in interface{}, typeName string) (interface{}, error) {
	if dtype := stringutil.ParseType(typeName); dtype != stringutil.Invalid {
		if v, err := stringutil.ConvertTo(dtype, in); err == nil {
			return v, nil
		}
	}

	if inStr, err := stringutil.ToString(in); err == nil {
		return inStr, nil
	} else {
		return in, nil
	}
}

func Get(data interface{}, key string, fallback ...interface{}) interface{} {
	data = typeutil.ResolveValue(data)

	if typeutil.IsKind(data, reflect.Map) {
		dataV := reflect.ValueOf(data)

		if valueV := dataV.MapIndex(reflect.ValueOf(key)); valueV.IsValid() {
			if valueI := valueV.Interface(); !typeutil.IsZero(valueI) {
				return valueI
			}
		}
	}

	if len(fallback) > 0 {
		return fallback[0]
	} else {
		return nil
	}
}

func DeepGet(data interface{}, path []string, fallbacks ...interface{}) interface{} {
	current := typeutil.ResolveValue(data)

	if len(fallbacks) == 0 {
		fallbacks = []interface{}{nil}
	}

	fallback := fallbacks[0]

	for i := 0; i < len(path); i++ {
		part := path[i]

		// fmt.Printf("dg:%s = %v (%T)\n", part, current, current)

		dValue := reflect.ValueOf(current)

		// if this value is not valid, return fallback here
		if !dValue.IsValid() {
			return fallback
		}

		dType := dValue.Type()

		// for pointers and interfaces, get the underlying type
		switch dType.Kind() {
		case reflect.Interface, reflect.Ptr:
			dType = dType.Elem()
		}

		switch dType.Kind() {
		//  arrays
		case reflect.Slice, reflect.Array:
			if stringutil.IsInteger(part) {
				if partIndex, err := strconv.Atoi(part); err == nil {
					if partIndex < dValue.Len() {
						if value := dValue.Index(partIndex).Interface(); value != nil {
							current = value
							continue
						}
					}
				}
			}

			return fallback

		//  maps
		case reflect.Map:
			if mapValue := dValue.MapIndex(
				reflect.ValueOf(part),
			); mapValue.IsValid() {
				current = mapValue.Interface()
			} else {
				return fallback
			}

		// structs
		case reflect.Struct:
			if dValue.Type().Kind() == reflect.Ptr {
				dValue = dValue.Elem()
			}

			if structField := dValue.FieldByName(part); structField.IsValid() && structField.CanInterface() {
				current = structField.Interface()
				continue
			}

		// attempting to retrieve nested data from a scalar value; return fallback
		default:
			return fallback
		}

	}

	return current
}

func DeepGetBool(data interface{}, path []string) bool {
	vI := DeepGet(data, path, false)

	if v, ok := vI.(bool); ok && v {
		return true
	}

	return false
}

func DeepGetString(data interface{}, path []string) string {

	if v, err := stringutil.ToString(DeepGet(data, path)); err == nil {
		return v
	}

	return ``
}

func DeepSet(data interface{}, path []string, value interface{}) interface{} {
	if len(path) == 0 {
		return data
	}

	var first = path[0]
	var rest = make([]string, 0)

	if len(path) > 1 {
		rest = path[1:]
	}

	//  Leaf Nodes
	//    this is where the value we're setting actually gets set/appended
	if len(rest) == 0 {
		switch data.(type) {
		//  parent element is an ARRAY
		case []interface{}:
			return append(data.([]interface{}), value)

			//  parent element is a MAP
		case map[string]interface{}:
			dataMap := data.(map[string]interface{})
			dataMap[first] = value

			return dataMap
		}
	} else {
		//  Array Embedding
		//    this is where keys that are actually array indices get processed
		//  ================================
		//  is `first' numeric (an array index)
		if stringutil.IsInteger(rest[0]) {
			switch data.(type) {
			case map[string]interface{}:
				dataMap := data.(map[string]interface{})

				//  is the value at `first' in the map isn't present or isn't an array, create it
				//  -------->
				curVal, _ := dataMap[first]

				switch curVal.(type) {
				case []interface{}:
				default:
					dataMap[first] = make([]interface{}, 0)
					curVal, _ = dataMap[first]
				}
				//  <--------|

				//  recurse into our cool array and do awesome stuff with it
				dataMap[first] = DeepSet(curVal.([]interface{}), rest, value).([]interface{})
				return dataMap
			default:
				// log.Printf("WHAT %s/%s", first, rest)
			}

			//  Intermediate Map Processing
			//    this is where branch nodes get created and populated via recursion
			//    depending on the data type of the input `data', non-existent maps
			//    will be created and either set to `data[first]' (the map)
			//    or appended to `data[first]' (the array)
			//  ================================
		} else {
			switch data.(type) {
			//  handle arrays of maps
			case []interface{}:
				dataArray := data.([]interface{})

				if curIndex, err := strconv.Atoi(first); err == nil {
					if curIndex >= len(dataArray) {
						for add := len(dataArray); add <= curIndex; add++ {
							dataArray = append(dataArray, make(map[string]interface{}))
						}
					}

					if curIndex < len(dataArray) {
						dataArray[curIndex] = DeepSet(dataArray[curIndex], rest, value)
						return dataArray
					}
				}

				//  handle good old fashioned maps-of-maps
			case map[string]interface{}:
				dataMap := data.(map[string]interface{})

				//  is the value at `first' in the map isn't present or isn't a map, create it
				//  -------->
				curVal, _ := dataMap[first]

				switch curVal.(type) {
				case map[string]interface{}:
				default:
					dataMap[first] = make(map[string]interface{})
					curVal, _ = dataMap[first]
				}
				//  <--------|

				dataMap[first] = DeepSet(dataMap[first], rest, value)
				return dataMap
			}
		}
	}

	return data
}

func Append(maps ...map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})

	for _, mapV := range maps {
		for k, v := range mapV {
			out[k] = v
		}
	}

	return out
}

func Pluck(sliceOfMaps interface{}, key []string) []interface{} {
	rv := make([]interface{}, 0)

	if sliceOfMaps == nil {
		return rv
	}

	inV := reflect.ValueOf(sliceOfMaps)

	if inV.IsValid() {
		if inV.Kind() == reflect.Interface {
			inV = inV.Elem()
		}

		if inV.IsValid() {
			if inV.Kind() == reflect.Ptr {
				inV = inV.Elem()
			}

			if inV.IsValid() {
				switch inV.Kind() {
				case reflect.Slice, reflect.Array:
					for i := 0; i < inV.Len(); i++ {
						if mapV := inV.Index(i); mapV.IsValid() {
							if mapV.Kind() == reflect.Interface {
								mapV = mapV.Elem()
							}

							if mapV.IsValid() {
								if mapV.Kind() == reflect.Ptr {
									mapV = mapV.Elem()
								}

								if mapV.IsValid() {
									if mapV.Kind() == reflect.Map {
										if v := DeepGet(mapV.Interface(), key, nil); v != nil {
											rv = append(rv, v)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return rv
}

func Walk(input interface{}, walkFn WalkFunc) error {
	return walkGeneric(input, nil, walkFn)
}

func walkGeneric(parent interface{}, path []string, walkFn WalkFunc) error {
	if parent == nil {
		return nil
	}

	parentV := reflect.ValueOf(parent)

	if parentV.Kind() == reflect.Ptr {
		parentV = parentV.Elem()
		parent = parentV.Interface()
	}

	switch parentV.Kind() {
	case reflect.Map:
		if err := walkFn(parent, path, false); err != nil {
			return returnSkipOrErr(err)
		}

		for _, key := range parentV.MapKeys() {
			valueV := parentV.MapIndex(key)
			subpath := append(path, fmt.Sprintf("%v", key.Interface()))

			if err := walkGeneric(valueV.Interface(), subpath, walkFn); err != nil {
				return returnSkipOrErr(err)
			}
		}

	case reflect.Slice, reflect.Array:
		if err := walkFn(parent, path, false); err != nil {
			return returnSkipOrErr(err)
		}

		for i := 0; i < parentV.Len(); i++ {
			valueV := parentV.Index(i)
			subpath := append(path, fmt.Sprintf("%v", i))

			if err := walkGeneric(valueV.Interface(), subpath, walkFn); err != nil {
				return returnSkipOrErr(err)
			}
		}

	case reflect.Struct:
		if err := walkFn(parent, path, false); err != nil {
			return returnSkipOrErr(err)
		}

		for i := 0; i < parentV.NumField(); i++ {
			fieldV := parentV.Type().Field(i)

			// only operate on exported fields
			if fieldV.PkgPath == `` {
				valueV := parentV.Field(i)
				var subpath []string

				// if this field is embedded, don't add it to the path list because
				// it should be considered a first-class member of the parent struct
				if fieldV.Anonymous {
					subpath = path
				} else if name := fieldNameFromReflect(fieldV); name != `-` {
					subpath = append(path, name)
				} else {
					continue
				}

				if err := walkGeneric(valueV.Interface(), subpath, walkFn); err != nil {
					return returnSkipOrErr(err)
				}
			}
		}

	default:
		return walkFn(parent, path, true)
	}

	return nil
}

func Compact(input map[string]interface{}) (map[string]interface{}, error) {
	output := make(map[string]interface{})

	if err := Walk(input, func(value interface{}, path []string, isLeaf bool) error {
		if isLeaf {
			if !typeutil.IsEmpty(value) {
				DeepSet(output, path, value)
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return output, nil
}

func Merge(first interface{}, second interface{}) (map[string]interface{}, error) {
	if first != nil && !typeutil.IsKind(first, reflect.Map) {
		return nil, fmt.Errorf("first argument must be a map, got %T", first)
	}

	if second != nil && !typeutil.IsKind(second, reflect.Map) {
		return nil, fmt.Errorf("second argument must be a map, got %T", second)
	}

	output := make(map[string]interface{})

	if err := Walk(first, func(value interface{}, path []string, isLeaf bool) error {
		if isLeaf {
			DeepSet(output, path, value)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	if err := Walk(second, func(value interface{}, path []string, isLeaf bool) error {
		if isLeaf {
			if value != nil {
				if currentValue := DeepGet(output, path, nil); currentValue == nil {
					DeepSet(output, path, value)
				} else {
					currentV := reflect.ValueOf(currentValue)

					switch currentV.Type().Kind() {
					case reflect.Slice, reflect.Array:
						newPath := append(path, fmt.Sprintf("%d", currentV.Len()))
						DeepSet(output, newPath, value)

					default:
						if currentValue == value {
							return nil
						}

						DeepSet(output, path, []interface{}{currentValue, value})
					}
				}
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return output, nil
}

func Stringify(input map[string]interface{}) map[string]string {
	output := make(map[string]string)

	if err := Walk(input, func(value interface{}, path []string, isLeaf bool) error {
		if isLeaf {
			if !typeutil.IsEmpty(value) {
				DeepSet(output, path, stringutil.MustString(value))
			}
		}

		return nil
	}); err != nil {
		panic(err.Error())
	}

	return output
}

func Autotype(input interface{}) map[string]interface{} {
	output := make(map[string]interface{})

	if err := Walk(input, func(value interface{}, path []string, isLeaf bool) error {
		// if we encounter a Variant, use its autotyping instead of walking the struct
		if valueVar, ok := value.(typeutil.Variant); ok {
			value = valueVar.Auto()

			if !typeutil.IsEmpty(value) {
				DeepSet(output, path, value)
			}

			return SkipDescendants
		} else if isLeaf {
			if !typeutil.IsEmpty(value) {
				DeepSet(output, path, stringutil.Autotype(value))
			}
		}

		return nil
	}); err != nil {
		panic(err.Error())
	}

	return output
}

func Apply(input interface{}, fn ApplyFunc) map[string]interface{} {
	output := make(map[string]interface{})

	if err := Walk(input, func(value interface{}, path []string, isLeaf bool) error {
		if isLeaf {
			if fn != nil {
				if out, ok := fn(path, value); ok {
					DeepSet(output, path, out)
					return nil
				}
			}

			DeepSet(output, path, value)
		}

		return nil
	}); err != nil {
		panic(err.Error())
	}

	return output
}

func DeepCopy(input interface{}) map[string]interface{} {
	return Apply(input, nil)
}

func returnSkipOrErr(err error) error {
	if err == SkipDescendants {
		return nil
	} else {
		return err
	}
}

func fieldNameFromReflect(field reflect.StructField) string {
	if tag := field.Tag.Get(UnmarshalStructTag); tag != `` {
		if name, _ := stringutil.SplitPair(tag, `,`); name != `` {
			return name
		}
	}

	return field.Name
}
