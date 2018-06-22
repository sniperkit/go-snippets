package human

import (
	"encoding"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"unicode"

	"github.com/hashicorp/go-multierror"
	"github.com/speijnik/go-errortree"
)

// Encoder writes human readable text to an output stream.
type Encoder struct {
	stream      *FlushableBuffer
	tagName     string
	indent      uint
	listSymbols []string
}

// Encode writes the human encoding of v to the stream.
func (e *Encoder) Encode(v interface{}) error {
	value := reflect.ValueOf(v)
	if err := e.encodeValue(v, value, -1, false); err != nil {
		e.stream.Reset()
		return err
	}

	_, err := e.stream.Flush()
	e.stream.Reset()
	return err
}

func (e *Encoder) encodeStruct(v reflect.Value, indentLevel int, inList bool) (err error) {
	t := v.Type()

	if v.Kind() == reflect.Ptr && v.IsValid() && !v.IsNil() {
		v = v.Elem()
		t = v.Type()
	}

	if !v.IsValid() {
		// no-op for nil/invalid values
		return
	}

	for i := 0; i < t.NumField(); i++ {
		fieldDefinition := t.Field(i)
		fieldValue := v.Field(i)

		if fieldDefinition.Anonymous {
			// Anonymous field handling

			if fieldValue.Kind() != reflect.Struct && fieldValue.Kind() != reflect.Ptr {
				// skip anonymous field that is neither a struct, nor a pointer
				continue
			} else if fieldValue.Kind() == reflect.Ptr && (!fieldValue.IsValid() || fieldValue.IsNil()) {
				// skip anonymous nil pointers
				continue
			} else if fieldValue.Kind() == reflect.Ptr || fieldValue.Kind() == reflect.Struct {
				// We are handling a valid pointer or a struct

				if fieldValue.Kind() == reflect.Ptr && fieldValue.Elem().Kind() != reflect.Struct {
					// Not a struct pointer, skip anonymous field
					continue
				} else if fieldValue.Kind() == reflect.Ptr {
					// Struct pointer, dereference pointer
					fieldValue = fieldValue.Elem()
				}

				// Getting this far means we are handling a struct
				if fieldErr := e.encodeStruct(fieldValue, indentLevel, inList); fieldErr != nil {
					err = errortree.Add(err, fieldDefinition.Name, fieldErr)
				}

				// We continue in any way, to skip the handling below
				continue
			}
		}

		fieldName := fieldDefinition.Name

		if !unicode.IsUpper([]rune(fieldName)[0]) {
			// Ignore private fields
			continue
		}

		fieldName, omitEmpty, tagErr := parseTagFromStructField(fieldDefinition, e.tagName)
		if tagErr != nil {
			// Parsing the tag failed, ignore the field and carry on
			err = multierror.Append(err, tagErr)
			continue
		}

		fieldInterface := fieldValue.Interface()
		if fieldInterface == nil || fieldName == "-" || (omitEmpty && IsNilOrEmpty(fieldInterface, fieldValue)) {
			// Skip field if:
			// - field is a nil-value
			// - field name specifies that the field shall be omitted
			// - omitEmpty is set and the field is nil or empty
			continue
		}

		// Getting this far means we are handling a non-empty field
		// if the struct is in a list adapt the first element's indent to the list symbol
		if inList {
			fmt.Fprint(e.stream, " "+fieldName+":")
			inList = false
		} else {
			fmt.Fprint(e.stream, strings.Repeat(" ", int(e.indent)*indentLevel)+fieldName+":")
		}
		if fieldEncodeErr := e.encodeValue(fieldInterface, fieldValue, indentLevel, false); fieldEncodeErr != nil {
			err = errortree.Add(err, fieldName, fieldEncodeErr)
		}
	}

	return
}

func (e *Encoder) encodeSlice(v reflect.Value, indentLevel int) error {

	listSymbol := e.listSymbols[int(indentLevel-1)%len(e.listSymbols)]

	for i := 0; i < v.Len(); i++ {
		valueV := v.Index(i)
		valueI := valueV.Interface()
		fmt.Fprint(e.stream, strings.Repeat(" ", int(e.indent)*indentLevel)+listSymbol)
		if err := e.encodeValue(valueI, valueV, indentLevel, true); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeMap(v reflect.Value, indentLevel int) error {

	listSymbol := e.listSymbols[int(indentLevel-1)%len(e.listSymbols)]

	keys := v.MapKeys()

	mapKeysStringMap := make(map[string]reflect.Value, len(keys))
	mapKeyStringList := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		keyV := keys[i]
		keyI := keyV.Interface()
		keyString := fmt.Sprint(keyI)
		mapKeysStringMap[keyString] = keyV
		mapKeyStringList[i] = keyString
	}

	sort.Strings(mapKeyStringList)

	for _, keyString := range mapKeyStringList {
		keyV := mapKeysStringMap[keyString]
		valueV := v.MapIndex(keyV)
		fmt.Fprint(e.stream, strings.Repeat(" ", int(e.indent)*indentLevel)+listSymbol+" "+keyString+":")
		if err := e.encodeValue(valueV.Interface(), valueV, indentLevel, true); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeValue(i interface{}, v reflect.Value, indentLevel int, inList bool) (err error) {
	// At this point it is safe to get rid of a possible pointer...
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	} else if v.Kind() == reflect.Ptr {
		// No-op for nil-pointers
		return
	}

	// Check if the passed interface implements encoding.TextMarshaler, in which case we use the marshaler
	// for generating the value
	if marshaler, ok := i.(encoding.TextMarshaler); ok {
		text, marshalErr := marshaler.MarshalText()
		if marshalErr != nil {
			err = marshalErr
		}
		// As MarshalText is expected to return a textual representation, print it to our stream
		fmt.Fprintf(e.stream, " %s\n", text)
		return
	} else if stringer, ok := i.(fmt.Stringer); ok {
		fmt.Fprintf(e.stream, " %s\n", stringer.String())
		return
	}

	// Per-type handling
	switch v.Kind() {
	case reflect.Struct:
		// Handle struct
		if !inList {
			fmt.Fprintln(e.stream, "")
		}
		err = e.encodeStruct(v, indentLevel+1, inList)
	case reflect.Slice, reflect.Array:
		// Handle slice
		fmt.Fprintln(e.stream, "")
		err = e.encodeSlice(v, indentLevel+1)
	case reflect.Map:
		// Handle map
		fmt.Fprintln(e.stream, "")
		err = e.encodeMap(v, indentLevel+1)

	default:
		// All other types are mapped as-is
		// missuse Fprint's sepereration spaces to introduce a space in front of the value
		fmt.Fprintln(e.stream, "", i)
	}

	return
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer, opts ...Option) (encoder *Encoder, err error) {
	encoder = &Encoder{
		stream: NewFlushableBuffer(w),
	}

	// apply options
	for _, opt := range append(defaultOptions, opts...) {
		if optErr := opt(encoder); optErr != nil {
			err = multierror.Append(err, optErr)
		}
	}
	// check if any option returned an error and set encoder nil
	if err != nil {
		encoder = nil
	}
	return
}
