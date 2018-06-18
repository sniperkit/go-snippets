package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ConvertType int

const (
	Invalid ConvertType = iota
	Nil
	String
	Boolean
	Float
	Integer
	Time
	Bytes
)

func (self ConvertType) String() string {
	switch self {
	case String:
		return `str`
	case Boolean:
		return `bool`
	case Float:
		return `float`
	case Integer:
		return `int`
	case Time:
		return `time`
	case Bytes:
		return `bytes`
	default:
		return ``
	}
}

var rxLeadingZeroes = regexp.MustCompile(`^0+\d+$`)
var NilStrings = []string{`null`, `NULL`, `<nil>`, `nil`, `Nil`, `None`, `undefined`, ``}
var BooleanTrueValues = []string{`true`, `yes`, `on`}
var BooleanFalseValues = []string{`false`, `no`, `off`}

var TimeFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC850,
	time.RFC822,
	time.RFC822Z,
	time.RFC1123,
	time.RFC1123Z,
	time.Kitchen,
	`2006-01-02 15:04:05.000000000`,
	`2006-01-02 15:04:05.000000`,
	`2006-01-02 15:04:05.000`,
	`2006-01-02 15:04:05 -0700 MST`,
	`2006-01-02 15:04:05Z07:00`,
	`2006-01-02 15:04:05`,
	`2006-01-02 15:04`,
	`2006-01-02`,
	`2006-01-02T15:04:05.000000000`,
	`2006-01-02T15:04:05.000000`,
	`2006-01-02T15:04:05.000`,
	`2006-01-02T15:04:05 -0700 MST`,
	`2006-01-02T15:04:05Z07:00`,
	`2006-01-02T15:04:05`,
}

func ToString(in interface{}) (string, error) {
	if in == nil {
		return ``, nil
	}

	if inStr, ok := in.(fmt.Stringer); ok {
		return inStr.String(), nil
	}

	if inT := reflect.TypeOf(in); inT != nil {
		switch inT.Kind() {
		case reflect.Float32:
			return strconv.FormatFloat(float64(in.(float32)), 'f', -1, 32), nil
		case reflect.Float64:
			return strconv.FormatFloat(in.(float64), 'f', -1, 64), nil
		case reflect.Bool:
			return strconv.FormatBool(in.(bool)), nil
		case reflect.String:
			if inStr, ok := in.(string); ok {
				return inStr, nil
			}
		}

		if !IsKind(in, CompoundTypes...) {
			return fmt.Sprintf("%v", in), nil
		}
	}

	return ``, fmt.Errorf("Unable to convert type '%T' to string", in)
}

func ConvertTo(toType ConvertType, inI interface{}) (interface{}, error) {
	var inS string
	var inSerr error

	if inV, ok := inI.(reflect.Value); ok {
		if inV.CanInterface() {
			inI = inV.Interface()
		} else {
			return nil, fmt.Errorf("reflect.Value given, but cannot retrieve interface value")
		}
	}

	inS, inSerr = ToString(inI)

	switch toType {
	case Float:
		if inS == `` {
			return float64(0), nil
		}

		return strconv.ParseFloat(inS, 64)
	case Integer:
		if inTime, ok := inI.(time.Time); ok {
			return int64(inTime.UnixNano()), nil
		} else if inS == `` {
			return int64(0), nil
		} else if layout := DetectTimeFormat(inS); layout != `` && layout != `epoch` {
			if tm, err := time.Parse(layout, inS); err == nil {
				return tm.UnixNano(), nil
			} else {
				return nil, err
			}
		}

		return strconv.ParseInt(inS, 10, 64)
	case Boolean:
		if inI == nil {
			return false, nil
		}

		if IsBooleanTrue(inS) {
			return true, nil
		} else if IsBooleanFalse(inS) {
			return false, nil
		} else {
			return nil, fmt.Errorf("Cannot convert '%s' into a boolean value", inS)
		}
	case Time:
		if inTime, ok := inI.(time.Time); ok {
			return inTime, nil
		}

		inS = strings.Trim(strings.TrimSpace(inS), `"'`)

		if DetectTimeFormat(inS) == `epoch` {
			if v, err := strconv.ParseInt(inS, 10, 64); err == nil {
				return time.Unix(v, 0), nil
			}
		}

		for _, format := range TimeFormats {
			if tm, err := time.Parse(format, strings.TrimSpace(inS)); err == nil {
				return tm, nil
			}
		}

		switch inS {
		case `now`:
			return time.Now(), nil
		default:
			// handle time zero values
			tmS := strings.Map(func(r rune) rune {
				switch r {
				case '-', ':', ' ', 'T', 'Z':
					return '0'
				}

				return r
			}, inS)

			if v, err := strconv.ParseInt(tmS, 10, 64); err == nil && v == 0 {
				return time.Time{}, nil
			}

			return nil, fmt.Errorf("Cannot convert '%s' into a date/time value", inS)
		}

	case Bytes:
		if inI == nil {
			return []byte{}, nil
		} else if inB, ok := inI.([]byte); ok {
			return inB, nil
		} else if inB, ok := inI.([]uint8); ok {
			outB := make([]byte, len(inB))

			for i, v := range inB {
				outB[i] = byte(v)
			}

			return outB, nil
		} else if IsKind(inI, reflect.Slice, reflect.Array) {
			outB := make([]byte, 0)

			if err := SliceEach(inI, func(i int, value interface{}) error {
				if bb, err := ConvertToInteger(value); err == nil {
					outB = append(outB, byte(bb))
					return nil
				} else {
					return err
				}
			}); err == nil {
				return outB, nil
			} else {
				return nil, err
			}
		} else if inSerr == nil {
			return []byte(inS), inSerr
		} else {
			return nil, inSerr
		}

	case String:
		// special case: assume incoming byte slices are actually strings
		if inB, ok := inI.([]byte); ok {
			// convert byte slices to strings directly
			return string(inB), nil
		} else if inB, ok := inI.([]uint8); ok {
			// convert byte slices to strings directly
			return string(inB), nil
		} else if inV := reflect.ValueOf(inI); inV.Kind() == reflect.Ptr {
			// dereference pointers to strings and stringify the result
			inS, inSerr = ToString(inV.Elem())
		}

		return inS, inSerr

	default:
		return inI, nil
	}
}

func ConvertToInteger(in interface{}) (int64, error) {
	if v, err := ConvertTo(Integer, in); err == nil {
		return v.(int64), nil
	} else {
		return int64(0), err
	}
}

func ConvertToFloat(in interface{}) (float64, error) {
	if v, err := ConvertTo(Float, in); err == nil {
		return v.(float64), nil
	} else {
		return float64(0.0), err
	}
}

func ConvertToString(in interface{}) (string, error) {
	if v, err := ConvertTo(String, in); err == nil {
		return v.(string), nil
	} else {
		return ``, err
	}
}

func ConvertToBool(in interface{}) (bool, error) {
	if v, err := ConvertTo(Boolean, in); err == nil {
		return v.(bool), nil
	} else {
		return false, err
	}
}

func ConvertToTime(in interface{}) (time.Time, error) {
	switch in.(type) {
	case time.Time:
		return in.(time.Time), nil
	default:
		if v, err := ConvertTo(Time, in); err == nil {
			return v.(time.Time), nil
		} else {
			return time.Time{}, err
		}
	}
}

func ConvertToBytes(in interface{}) ([]byte, error) {
	if v, err := ConvertTo(Bytes, in); err == nil {
		return v.([]byte), nil
	} else {
		return nil, err
	}
}

func DetectTimeFormat(in string) string {
	if IsInteger(in) {
		return `epoch`
	}

	for _, layout := range TimeFormats {
		if _, err := time.Parse(layout, in); err == nil {
			return layout
		}
	}

	return ``
}

func Autotype(in interface{}) interface{} {
	if in == nil {
		return nil
	}

	if IsTime(in) {
		if v, err := ConvertTo(Time, in); err == nil {
			return v
		}
	}

	// effectively, this detects strings that are numeric, but have leading zeroes.
	// we should treat those as meaningful and return that string outright
	//
	// (e.g.: handle the "US Zip Code" problem)
	if vStr, ok := in.(string); ok {
		if rxLeadingZeroes.MatchString(vStr) {
			return vStr
		}

		// certain known string values should convert to nil directly
		for _, nilStr := range NilStrings {
			if vStr == nilStr {
				return nil
			}
		}
	}

	for _, ctype := range []ConvertType{
		Boolean,
		Integer,
		Float,
		String,
	} {
		if value, err := ConvertTo(ctype, in); err == nil {
			return value
		}
	}

	return in
}

func DetectConvertType(in interface{}) ConvertType {
	if in == nil {
		return Nil
	}

	if IsTime(in) {
		return Time
	}

	// effectively, this detects strings that are numeric, but have leading zeroes.
	// we should treat those as meaningful and return that string outright
	//
	// (e.g.: handle the "US Zip Code" problem)
	if vStr, ok := in.(string); ok {
		if rxLeadingZeroes.MatchString(vStr) {
			return String
		}

		// certain known string values should convert to nil directly
		for _, nilStr := range NilStrings {
			if vStr == nilStr {
				return Nil
			}
		}
	}

	for _, ctype := range []ConvertType{
		Boolean,
		Integer,
		Float,
		String,
	} {
		if _, err := ConvertTo(ctype, in); err == nil {
			return ctype
		}
	}

	return Invalid
}
