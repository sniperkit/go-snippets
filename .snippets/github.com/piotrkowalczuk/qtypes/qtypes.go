// Package qtypes provides set of types that helps to build complex protobuf messages that can express conditional statements.
package qtypes

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	pbts "github.com/golang/protobuf/ptypes/timestamp"
)

const (
	arraySeparator = ","
	// Null ...
	Null = "null"
	// NotNull ...
	NotNull = "nnull"
	// Equal ...
	Equal = "eq"
	// NotEqual ...
	NotEqual = "neq"
	// GreaterThan ...
	GreaterThan = "gt"
	// NotGreaterThan ...
	NotGreaterThan = "ngt"
	// GreaterThanOrEqual ...
	GreaterThanOrEqual = "gte"
	// NotGreaterThanOrEqual ...
	NotGreaterThanOrEqual = "ngte"
	// LessThan ...
	LessThan = "lt"
	// NotLessThan ...
	NotLessThan = "nlt"
	// LessThanOrEqual ...
	LessThanOrEqual = "lte"
	// NotLessThanOrEqual ...
	NotLessThanOrEqual = "nlte"
	// Between ...
	Between = "bw"
	// NotBetween ...
	NotBetween = "nbw"
	// HasPrefix ...
	HasPrefix = "hp"
	// HasPrefixInsensitive ...
	HasPrefixInsensitive = "hpi"
	// HasSuffix ...
	HasSuffix = "hs"
	// HasSuffixInsensitive ...
	HasSuffixInsensitive = "hsi"
	// Substring ...
	Substring = "sub"
	// SubstringInsensitive` ...
	SubstringInsensitive = "subi"
	// HasElement ...
	HasElement = "he"
	// HasAnyElement ...
	HasAnyElement = "hae"
	// HasAllElements ...
	HasAllElements = "hle"
	// In ...
	In = "in"
	// NotIn ...
	NotIn = "nin"
	// Pattern ...
	Pattern = "rgx"
	// MinLength ...
	MinLength = "minl"
	// MaxLength ...
	MaxLength = "maxl"
	// Contains ...
	Contains = "cts"
	// IsContainedBy ...
	IsContainedBy = "icb"
	// Overlap ...
	Overlap = "ovl"
)

var (
	prefixes = map[string]string{
		Null:                  Null + ":",
		NotNull:               NotNull + ":",
		Equal:                 Equal + ":",
		NotEqual:              NotEqual + ":",
		GreaterThan:           GreaterThan + ":",
		NotGreaterThan:        NotGreaterThan + ":",
		GreaterThanOrEqual:    GreaterThanOrEqual + ":",
		NotGreaterThanOrEqual: NotGreaterThanOrEqual + ":",
		LessThan:              LessThan + ":",
		NotLessThan:           NotLessThan + ":",
		LessThanOrEqual:       LessThanOrEqual + ":",
		NotLessThanOrEqual:    NotLessThanOrEqual + ":",
		Between:               Between + ":",
		NotBetween:            NotBetween + ":",
		HasPrefix:             HasPrefix + ":",
		HasPrefixInsensitive:  HasPrefixInsensitive + ":",
		HasSuffix:             HasSuffix + ":",
		HasSuffixInsensitive:  HasSuffixInsensitive + ":",
		In:                    In + ":",
		NotIn:                 NotIn + ":",
		Substring:             Substring + ":",
		SubstringInsensitive:  SubstringInsensitive + ":",
		Pattern:               Pattern + ":",
		MinLength:             MinLength + ":",
		MaxLength:             MaxLength + ":",
		Contains:              Contains + ":",
		IsContainedBy:         IsContainedBy + ":",
		Overlap:               Overlap + ":",
		HasElement:            HasElement + ":",
		HasAnyElement:         HasAnyElement + ":",
		HasAllElements:        HasAllElements + ":",
	}
)

// Value returns first value or empty string if none.
func (qs *String) Value() string {
	if len(qs.Values) == 0 {
		return ""
	}

	return qs.Values[0]
}

// ParseString allocates new String object based on given string.
// If string is prefixed with known operator e.g. 'hp:New'
// returned object will get same type.
func ParseString(s string) *String {
	if s == "" {
		return &String{}
	}

	for c, p := range prefixes {
		if strings.HasPrefix(s, p) {
			t, n, i := queryType(c)
			return &String{
				Values:      strings.Split(strings.TrimPrefix(s, p), arraySeparator),
				Type:        t,
				Negation:    n,
				Insensitive: i,
				Valid:       true,
			}
		}
	}
	return &String{
		Values: strings.Split(s, arraySeparator),
		Type:   QueryType_EQUAL,
		Valid:  true,
	}
}

// EqualString ...
func EqualString(s string) *String {
	return &String{
		Values: []string{s},
		Valid:  true,
		Type:   QueryType_EQUAL,
	}
}

// HasPrefixString ...
func HasPrefixString(s string) *String {
	return &String{
		Values: []string{s},
		Valid:  true,
		Type:   QueryType_HAS_PREFIX,
	}
}

// HasSuffixString ...
func HasSuffixString(s string) *String {
	return &String{
		Values: []string{s},
		Valid:  true,
		Type:   QueryType_HAS_SUFFIX,
	}
}

// SubString ...
func SubString(s string) *String {
	return &String{
		Values: []string{s},
		Valid:  true,
		Type:   QueryType_SUBSTRING,
	}
}

// NullString ...
func NullString() *String {
	return &String{
		Valid: true,
		Type:  QueryType_NULL,
	}
}

// NullInt64 allocates valid Int64 object of type not a number with given value.
func NullInt64() *Int64 {
	return &Int64{
		Valid: true,
		Type:  QueryType_NULL,
	}
}

// EqualInt64 allocates valid Int64 object of type equal with given value.
func EqualInt64(i int64) *Int64 {
	return &Int64{
		Values: []int64{i},
		Valid:  true,
		Type:   QueryType_EQUAL,
	}
}

// NotEqualInt64 allocates valid Int64 negated object of type equal with given value.
func NotEqualInt64(i int64) *Int64 {
	return &Int64{
		Values:   []int64{i},
		Valid:    true,
		Negation: true,
		Type:     QueryType_EQUAL,
	}
}

// InInt64 allocates valid Int64 object of type in with given values.
func InInt64(v ...int64) *Int64 {
	return &Int64{
		Values: v,
		Valid:  true,
		Type:   QueryType_IN,
	}
}

// BetweenInt64 allocates valid Int64 object of type between with given values.
func BetweenInt64(a, b int64) *Int64 {
	return &Int64{
		Values: []int64{a, b},
		Valid:  true,
		Type:   QueryType_BETWEEN,
	}
}

// GreaterInt64 allocates valid Int64 object of type greater with given value.
func GreaterInt64(i int64) *Int64 {
	return &Int64{
		Values: []int64{i},
		Valid:  true,
		Type:   QueryType_GREATER,
	}
}

// GreaterEqualInt64 allocates valid Int64 object of type greater equal with given value.
func GreaterEqualInt64(i int64) *Int64 {
	return &Int64{
		Values: []int64{i},
		Valid:  true,
		Type:   QueryType_GREATER_EQUAL,
	}
}

// LessInt64 allocates valid Int64 object of type less with given value.
func LessInt64(i int64) *Int64 {
	return &Int64{
		Values: []int64{i},
		Valid:  true,
		Type:   QueryType_LESS,
	}
}

// LessEqualInt64 allocates valid Int64 object of type less equal with given value.
func LessEqualInt64(i int64) *Int64 {
	return &Int64{
		Values: []int64{i},
		Valid:  true,
		Type:   QueryType_LESS_EQUAL,
	}
}

// Value ...
func (i *Int64) Value() int64 {
	if len(i.Values) == 0 {
		return 0
	}

	return i.Values[0]
}

// ParseInt64 ...
func ParseInt64(s string) (*Int64, error) {
	if s == "" {
		return &Int64{}, nil
	}
	incoming, t, n, _ := handleNumericPrefix(s)
	outgoing := make([]int64, 0, len(incoming))
	for i, v := range incoming {
		if v == "" {
			break
		}
		vv, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("qtypes: query int64 parsing error for key %d and value %v: %s", i, vv, err.Error())
		}
		outgoing = append(outgoing, vv)
	}
	return &Int64{
		Values:   outgoing,
		Type:     t,
		Negation: n,
		Valid:    true,
	}, nil
}

// EqualFloat64 allocates valid Float64 object of type equal with given value.
func EqualFloat64(i float64) *Float64 {
	return &Float64{
		Values: []float64{i},
		Valid:  true,
		Type:   QueryType_EQUAL,
	}
}

// BetweenFloat64 allocates valid Float64 object if both numbers are not 0 and from is not greater than to.
func BetweenFloat64(from, to float64) *Float64 {
	if from == 0 && to == 0 {
		return &Float64{}
	}
	if from > to {
		return &Float64{}
	}
	return &Float64{
		Values: []float64{from, to},
		Type:   QueryType_BETWEEN,
		Valid:  true,
	}
}

// Value returns first available value or 0 if none available.
func (f *Float64) Value() float64 {
	if len(f.Values) == 0 {
		return 0.0
	}

	return f.Values[0]
}

// ParseFloat64 ...
func ParseFloat64(s string) (*Float64, error) {
	if s == "" {
		return &Float64{}, nil
	}
	incoming, t, n, _ := handleNumericPrefix(s)

	outgoing := make([]float64, 0, len(incoming))
	for i, v := range incoming {
		if v == "" {
			break
		}
		vv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("qtypes: query float64 parsing error for valur %d: %s", i, err.Error())
		}
		outgoing = append(outgoing, vv)
	}
	return &Float64{
		Values:   outgoing,
		Type:     t,
		Negation: n,
		Valid:    true,
	}, nil
}

// BetweenTimestamp allocates valid Timestamp object if both timestamps are not nil
// and first is before the second.
func BetweenTimestamp(from, to *pbts.Timestamp) *Timestamp {
	if from == nil || to == nil {
		return &Timestamp{}
	}

	v := true
	if to.Seconds < from.Seconds {
		v = false
	}
	if to.Seconds == from.Seconds && to.Nanos < from.Nanos {
		v = false
	}
	return &Timestamp{
		Values: []*pbts.Timestamp{from, to},
		Type:   QueryType_BETWEEN,
		Valid:  v,
	}
}

// Value returns first value or nil if none.
func (t *Timestamp) Value() *pbts.Timestamp {
	if len(t.Values) == 0 {
		return nil
	}

	return t.Values[0]
}

// ParseTimestamp ...
func ParseTimestamp(s string) (*Timestamp, error) {
	if s == "" {
		return &Timestamp{}, nil
	}

	incoming, t, n, _ := handleNumericPrefix(s)

	outgoing := make([]*pbts.Timestamp, 0, len(incoming))
	for i, v := range incoming {
		if v == "" {
			break
		}
		t, err := time.Parse(time.RFC3339Nano, v)
		if err != nil {
			return nil, fmt.Errorf("qtypes: query timestamp parsing error for value %d: %s", i, err.Error())
		}
		tt, err := ptypes.TimestampProto(t)
		if err != nil {
			return nil, fmt.Errorf("qtypes: time to proto timestamp conversion error for value %d: %s", i, err.Error())
		}
		outgoing = append(outgoing, tt)
	}
	return &Timestamp{
		Values:   outgoing,
		Type:     t,
		Negation: n,
		Valid:    true,
	}, nil
}

func handleNumericPrefix(s string) (incoming []string, t QueryType, n, i bool) {
	if parts := strings.Split(s, ":"); len(parts) == 1 {
		return []string{s}, QueryType_EQUAL, false, false
	}
	for c, p := range prefixes {
		if strings.HasPrefix(s, p) {
			t, n, i = queryType(c)
			incoming = strings.Split(strings.TrimPrefix(s, p), arraySeparator)
		}
	}
	if len(incoming) == 0 {
		incoming = strings.Split(s, arraySeparator)
	}

	return
}

func queryType(p string) (t QueryType, n bool, i bool) {
	switch p {
	case Null:
		t = QueryType_NULL
	case NotNull:
		t = QueryType_NULL
		n = true
	case Equal:
		t = QueryType_EQUAL
	case NotEqual:
		t = QueryType_EQUAL
		n = true
	case GreaterThan:
		t = QueryType_GREATER
	case NotGreaterThan:
		t = QueryType_GREATER
		n = true
	case GreaterThanOrEqual:
		t = QueryType_GREATER_EQUAL
	case NotGreaterThanOrEqual:
		t = QueryType_GREATER_EQUAL
		n = true
	case LessThan:
		t = QueryType_LESS
	case NotLessThan:
		t = QueryType_LESS
		n = true
	case LessThanOrEqual:
		t = QueryType_LESS_EQUAL
	case NotLessThanOrEqual:
		t = QueryType_LESS_EQUAL
		n = true
	case Between:
		t = QueryType_BETWEEN
	case NotBetween:
		t = QueryType_BETWEEN
		n = true
	case HasElement:
		t = QueryType_HAS_ELEMENT
	case HasAllElements:
		t = QueryType_HAS_ALL_ELEMENTS
	case HasAnyElement:
		t = QueryType_HAS_ANY_ELEMENT
	case HasPrefix:
		t = QueryType_HAS_PREFIX
	case HasPrefixInsensitive:
		t = QueryType_HAS_PREFIX
		i = true
	case HasSuffix:
		t = QueryType_HAS_SUFFIX
	case HasSuffixInsensitive:
		t = QueryType_HAS_SUFFIX
		i = true
	case Substring:
		t = QueryType_SUBSTRING
	case SubstringInsensitive:
		t = QueryType_SUBSTRING
		i = true
	case Pattern:
		t = QueryType_PATTERN
	case MinLength:
		t = QueryType_MIN_LENGTH
	case MaxLength:
		t = QueryType_MAX_LENGTH
	case In:
		t = QueryType_IN
	case NotIn:
		t = QueryType_IN
		n = true
	case Contains:
		t = QueryType_CONTAINS
	case IsContainedBy:
		t = QueryType_IS_CONTAINED_BY
	case Overlap:
		t = QueryType_OVERLAP
	}
	return
}
