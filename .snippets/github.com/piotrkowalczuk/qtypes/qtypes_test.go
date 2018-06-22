package qtypes

import (
	"fmt"
	"reflect"
	"testing"

	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func Example() {
	query := struct {
		Name  *String
		Age   *Int64
		Money *Float64
	}{
		Name:  ParseString("sub:John"),
		Age:   GreaterInt64(18),
		Money: EqualFloat64(0.0),
	}

	fmt.Println(query.Name.Value())
	fmt.Println(query.Name.Type)
	fmt.Println(query.Age.Value())
	fmt.Println(query.Age.Type)
	fmt.Println(query.Money.Value())
	fmt.Println(query.Money.Type)

	// Output:
	// John
	// SUBSTRING
	// 18
	// GREATER
	// 0
	// EQUAL
}

func ExampleEqualString() {
	ex := EqualString("text")

	fmt.Println(ex.Valid)
	fmt.Println(ex.Negation)
	fmt.Println(ex.Type)
	fmt.Println(ex.Value())

	// Output:
	// true
	// false
	// EQUAL
	// text
}

func TestParseString(t *testing.T) {
	cases := map[string]struct {
		given    string
		expected String
	}{
		"null": {
			given: "null:",
			expected: String{
				Values: []string{""},
				Type:   QueryType_NULL,
				Valid:  true,
			},
		},
		"not-null": {
			given: "nnull:",
			expected: String{
				Values:   []string{""},
				Type:     QueryType_NULL,
				Valid:    true,
				Negation: true,
			},
		},
		"equal": {
			given: "eq:123",
			expected: String{
				Values: []string{"123"},
				Type:   QueryType_EQUAL,
				Valid:  true,
			},
		},
		"has-prefix": {
			given: "hp:New",
			expected: String{
				Values: []string{"New"},
				Type:   QueryType_HAS_PREFIX,
				Valid:  true,
			},
		},
		"has-prefix-insensitive": {
			given: "hpi:New",
			expected: String{
				Values:      []string{"New"},
				Type:        QueryType_HAS_PREFIX,
				Valid:       true,
				Insensitive: true,
			},
		},
		"has-suffix": {
			given: "hs:New",
			expected: String{
				Values: []string{"New"},
				Type:   QueryType_HAS_SUFFIX,
				Valid:  true,
			},
		},
		"has-suffix-insensitive": {
			given: "hsi:New",
			expected: String{
				Values:      []string{"New"},
				Type:        QueryType_HAS_SUFFIX,
				Valid:       true,
				Insensitive: true,
			},
		},
		"substring": {
			given: "sub:anything",
			expected: String{
				Values: []string{"anything"},
				Type:   QueryType_SUBSTRING,
				Valid:  true,
			},
		},
		"substring-insensitive": {
			given: "subi:anything",
			expected: String{
				Values:      []string{"anything"},
				Type:        QueryType_SUBSTRING,
				Valid:       true,
				Insensitive: true,
			},
		},
		"pattern": {
			given: "rgx:.*",
			expected: String{
				Values: []string{".*"},
				Type:   QueryType_PATTERN,
				Valid:  true,
			},
		},
		"max-length": {
			given: "maxl:4",
			expected: String{
				Values: []string{"4"},
				Type:   QueryType_MAX_LENGTH,
				Valid:  true,
			},
		},
		"min-length": {
			given: "minl:555",
			expected: String{
				Values: []string{"555"},
				Type:   QueryType_MIN_LENGTH,
				Valid:  true,
			},
		},
		"empty": {
			given:    "",
			expected: String{},
		},
		"without-condition": {
			given: "text",
			expected: String{
				Values: []string{"text"},
				Type:   QueryType_EQUAL,
				Valid:  true,
			},
		},
		"with-condition-but-without-value": {
			given: "neq:",
			expected: String{
				Values:   []string{""},
				Type:     QueryType_EQUAL,
				Valid:    true,
				Negation: true,
			},
		},
		"has-element": {
			given: "he:555",
			expected: String{
				Values: []string{"555"},
				Type:   QueryType_HAS_ELEMENT,
				Valid:  true,
			},
		},
		"has-any-elements": {
			given: "hae:555,222",
			expected: String{
				Values: []string{"555", "222"},
				Type:   QueryType_HAS_ANY_ELEMENT,
				Valid:  true,
			},
		},
		"has-all-elements": {
			given: "hle:111,222",
			expected: String{
				Values: []string{"111", "222"},
				Type:   QueryType_HAS_ALL_ELEMENTS,
				Valid:  true,
			},
		},
		"not-greater-than": {
			given: "ngt:111",
			expected: String{
				Values:   []string{"111"},
				Type:     QueryType_GREATER,
				Negation: true,
				Valid:    true,
			},
		},
		"greater-than": {
			given: "gt:111",
			expected: String{
				Values: []string{"111"},
				Type:   QueryType_GREATER,
				Valid:  true,
			},
		},
		"less-than": {
			given: "lt:111",
			expected: String{
				Values: []string{"111"},
				Type:   QueryType_LESS,
				Valid:  true,
			},
		},
	}

CasesLoop:
	for hint, c := range cases {
		got := ParseString(c.given)

		if got == nil {
			t.Error("unexpected nil")
			continue CasesLoop
		}
		if !reflect.DeepEqual(c.expected, *got) {
			t.Errorf("%s: wrong output,\nexpected:\n	%v\nbut got:\n	%v\n", hint, &c.expected, got)
		}
	}
}

func TestEqualString(t *testing.T) {
	values := []string{"a"}
	testString(t, EqualString(values[0]), false, true, QueryType_EQUAL, values...)
}

func TestSubString(t *testing.T) {
	values := []string{"a"}
	testString(t, SubString(values[0]), false, true, QueryType_SUBSTRING, values...)
}

func TestHasPrefixString(t *testing.T) {
	values := []string{"a"}
	testString(t, HasPrefixString(values[0]), false, true, QueryType_HAS_PREFIX, values...)
}

func TestHasSuffixString(t *testing.T) {
	values := []string{"a"}
	testString(t, HasSuffixString(values[0]), false, true, QueryType_HAS_SUFFIX, values...)
}

func TestNullString(t *testing.T) {
	testString(t, NullString(), false, true, QueryType_NULL)
}

func testString(t *testing.T, s *String, n, v bool, tp QueryType, values ...string) {
	if s.Negation != n {
		t.Errorf("wrong negation, exiected %t but got %t", n, s.Negation)
	}

	if len(values) > 0 {
		if s.Value() != values[0] {
			t.Errorf("wrong first value, expected %s but got %s", values[0], s.Value())
		}
		if len(s.Values) == len(values) {
			for j, v := range values {
				if s.Values[j] != v {
					t.Errorf("%d: wrong value, expected %s but got %s", j, v, s.Values[j])
				}
			}
		} else {
			t.Errorf("wrong number of values, expected %d but got %d", len(values), len(s.Values))
		}
	}
	if s.Valid != v {
		t.Errorf("expected valid to be %t", v)
	}
	if s.Type != tp {
		t.Errorf("wrong type, expected %s but got %s", tp, s.Type)
	}
}

func TestBetweenTimestamp(t *testing.T) {
	cases := map[string]struct {
		from     *timestamp.Timestamp
		to       *timestamp.Timestamp
		expected Timestamp
	}{
		"valid": {
			from: &timestamp.Timestamp{Seconds: 0, Nanos: 0},
			to:   &timestamp.Timestamp{Seconds: 0, Nanos: 1},
			expected: Timestamp{
				Valid:    true,
				Negation: false,
				Type:     QueryType_BETWEEN,
				Values: []*timestamp.Timestamp{
					{Seconds: 0, Nanos: 0},
					{Seconds: 0, Nanos: 1},
				},
			},
		},
		"after-first": {
			from: &timestamp.Timestamp{Seconds: 1, Nanos: 0},
			to:   &timestamp.Timestamp{Seconds: 0, Nanos: 0},
			expected: Timestamp{
				Valid: false,
				Type:  QueryType_BETWEEN,
				Values: []*timestamp.Timestamp{
					{Seconds: 1, Nanos: 0},
					{Seconds: 0, Nanos: 0},
				},
			},
		},
		"after-first-seconds": {
			from: &timestamp.Timestamp{Seconds: 1, Nanos: 1},
			to:   &timestamp.Timestamp{Seconds: 1, Nanos: 0},
			expected: Timestamp{
				Valid: false,
				Type:  QueryType_BETWEEN,
				Values: []*timestamp.Timestamp{
					{Seconds: 1, Nanos: 1},
					{Seconds: 1, Nanos: 0},
				},
			},
		},
		"nil-arguments": {
			from:     nil,
			to:       nil,
			expected: Timestamp{},
		},
		"nil-argument-first": {
			from:     nil,
			to:       &timestamp.Timestamp{Seconds: 0, Nanos: 1},
			expected: Timestamp{},
		},
		"nil-argument-second": {
			from:     &timestamp.Timestamp{Seconds: 0, Nanos: 1},
			to:       nil,
			expected: Timestamp{},
		},
	}

	for hint, c := range cases {
		bt := BetweenTimestamp(c.from, c.to)
		if !reflect.DeepEqual(c.expected, *bt) {
			t.Errorf("%s: unexpected output, expected:\n%v\ngot:\n%v\n", hint, c.expected, *bt)
		}
	}
}

func TestInt64_Value(t *testing.T) {
	cases := map[string]struct {
		given    Int64
		expected int64
	}{
		"single": {
			given: Int64{
				Values: []int64{1},
				Valid:  true,
				Type:   QueryType_EQUAL,
			},
			expected: 1,
		},
		"none": {
			given: Int64{
				Valid: true,
				Type:  QueryType_EQUAL,
			},
			expected: 0,
		},
		"multiple": {
			given: Int64{
				Values: []int64{3, 2, 1},
				Valid:  true,
				Type:   QueryType_EQUAL,
			},
			expected: 3,
		},
	}

	for hint, c := range cases {
		if c.given.Value() != c.expected {
			t.Errorf("%s: unexpected value, expected %d but got %d", hint, c.expected, c.given.Value())
		}
	}
}

func TestNaNInt64(t *testing.T) {
	testInt64(t, NullInt64(), false, true, QueryType_NULL)
}

func TestEqualInt64(t *testing.T) {
	value := int64(1111)
	testInt64(t, EqualInt64(value), false, true, QueryType_EQUAL, value)
}

func TestNotEqualInt64(t *testing.T) {
	value := int64(1111)
	testInt64(t, NotEqualInt64(value), true, true, QueryType_EQUAL, value)
}

func TestGreaterInt64(t *testing.T) {
	value := int64(1111)
	testInt64(t, GreaterInt64(value), false, true, QueryType_GREATER, value)
}

func TestGreaterEqualInt64(t *testing.T) {
	value := int64(1111)
	testInt64(t, GreaterEqualInt64(value), false, true, QueryType_GREATER_EQUAL, value)
}

func TestBetweenInt64(t *testing.T) {
	values := []int64{1111, 2222}
	testInt64(t, BetweenInt64(values[0], values[1]), false, true, QueryType_BETWEEN, values...)
}

func TestLessInt64(t *testing.T) {
	value := int64(1111)
	testInt64(t, LessInt64(value), false, true, QueryType_LESS, value)
}

func TestLessEqualInt64(t *testing.T) {
	value := int64(1111)
	testInt64(t, LessEqualInt64(value), false, true, QueryType_LESS_EQUAL, value)
}

func TestInInt64(t *testing.T) {
	values := []int64{1111, 2222, 3333, 4444}
	testInt64(t, InInt64(values...), false, true, QueryType_IN, values...)
}

func TestParseInt64(t *testing.T) {
	cases := map[string]struct {
		given    string
		expected Int64
	}{
		"empty": {
			given:    "",
			expected: Int64{},
		},
		"null": {
			given: "null:",
			expected: Int64{
				Values: []int64{},
				Type:   QueryType_NULL,
				Valid:  true,
			},
		},
		"number": {
			given: "15",
			expected: Int64{
				Values: []int64{15},
				Type:   QueryType_EQUAL,
				Valid:  true,
			},
		},
		"not-null": {
			given: "nnull:",
			expected: Int64{
				Values:   []int64{},
				Type:     QueryType_NULL,
				Valid:    true,
				Negation: true,
			},
		},
		"equal": {
			given: "eq:123",
			expected: Int64{
				Values: []int64{123},
				Type:   QueryType_EQUAL,
				Valid:  true,
			},
		},
		"not-equal": {
			given: "neq:123",
			expected: Int64{
				Values:   []int64{123},
				Type:     QueryType_EQUAL,
				Valid:    true,
				Negation: true,
			},
		},
		"greater": {
			given: "gt:555",
			expected: Int64{
				Values: []int64{555},
				Type:   QueryType_GREATER,
				Valid:  true,
			},
		},
		"greater-equal": {
			given: "gte:666",
			expected: Int64{
				Values: []int64{666},
				Type:   QueryType_GREATER_EQUAL,
				Valid:  true,
			},
		},
		"lesser": {
			given: "lt:777",
			expected: Int64{
				Values: []int64{777},
				Type:   QueryType_LESS,
				Valid:  true,
			},
		},
		"lesser-equal": {
			given: "lte:888",
			expected: Int64{
				Values: []int64{888},
				Type:   QueryType_LESS_EQUAL,
				Valid:  true,
			},
		},
		"between": {
			given: "bw:111,222",
			expected: Int64{
				Values: []int64{111, 222},
				Type:   QueryType_BETWEEN,
				Valid:  true,
			},
		},
		"not-between": {
			given: "nbw:111,222",
			expected: Int64{
				Values:   []int64{111, 222},
				Type:     QueryType_BETWEEN,
				Valid:    true,
				Negation: true,
			},
		},
	}

CasesLoop:
	for hint, c := range cases {
		got, err := ParseInt64(c.given)
		if err != nil {
			t.Errorf("%s: unexpected error: %s", hint, err.Error())
			continue CasesLoop
		}
		if got == nil {
			t.Errorf("%s: unexpected nil", hint)
			continue CasesLoop
		}
		if !reflect.DeepEqual(c.expected, *got) {
			t.Errorf("%s: wrong output,\nexpected:\n	%v\nbut got:\n	%v\n", hint, &c.expected, got)
		}
	}
}

func TestParseInt64_text(t *testing.T) {
	got, err := ParseInt64("ne:long-text")
	if err == nil {
		t.Fatal("expected error")
	}
	if got != nil {
		t.Fatal("expected nil")
	}
}

func TestParseTimestamp(t *testing.T) {
	parseTimestamp := func(t *testing.T, s string) *timestamp.Timestamp {
		pt, err := time.Parse(time.RFC3339Nano, s)
		if err != nil {
			t.Fatalf("string cant be parsed into time: %s", err.Error())
		}
		tt, err := ptypes.TimestampProto(pt)
		if err != nil {
			t.Fatalf("tmie cant be converted into timestamp: %s", err.Error())
		}
		return tt
	}
	cases := map[string]struct {
		given    string
		expected Timestamp
	}{
		"empty": {
			given:    "",
			expected: Timestamp{},
		},
		"null": {
			given: "null:",
			expected: Timestamp{
				Values: []*timestamp.Timestamp{},
				Type:   QueryType_NULL,
				Valid:  true,
			},
		},
		"not-null": {
			given: "nnull:",
			expected: Timestamp{
				Values:   []*timestamp.Timestamp{},
				Type:     QueryType_NULL,
				Valid:    true,
				Negation: true,
			},
		},
		"equal": {
			given: "eq:2009-11-10T23:00:00Z",
			expected: Timestamp{
				Values: []*timestamp.Timestamp{
					parseTimestamp(t, "2009-11-10T23:00:00Z"),
				},
				Type:  QueryType_EQUAL,
				Valid: true,
			},
		},
		"greater-equal": {
			given: "gte:2009-11-10T23:00:00Z",
			expected: Timestamp{
				Values: []*timestamp.Timestamp{
					parseTimestamp(t, "2009-11-10T23:00:00Z"),
				},
				Type:  QueryType_GREATER_EQUAL,
				Valid: true,
			},
		},
		"greater": {
			given: "gt:2009-11-10T23:00:00Z",
			expected: Timestamp{
				Values: []*timestamp.Timestamp{
					parseTimestamp(t, "2009-11-10T23:00:00Z"),
				},
				Type:  QueryType_GREATER,
				Valid: true,
			},
		},
		"less": {
			given: "lt:2009-11-10T23:00:00Z",
			expected: Timestamp{
				Values: []*timestamp.Timestamp{
					parseTimestamp(t, "2009-11-10T23:00:00Z"),
				},
				Type:  QueryType_LESS,
				Valid: true,
			},
		},
		"less-equal": {
			given: "lte:2009-11-10T23:00:00Z",
			expected: Timestamp{
				Values: []*timestamp.Timestamp{
					parseTimestamp(t, "2009-11-10T23:00:00Z"),
				},
				Type:  QueryType_LESS_EQUAL,
				Valid: true,
			},
		},
		"between": {
			given: "bw:2009-11-10T23:00:00Z,2009-12-10T23:00:00Z",
			expected: Timestamp{
				Values: []*timestamp.Timestamp{
					parseTimestamp(t, "2009-11-10T23:00:00Z"),
					parseTimestamp(t, "2009-12-10T23:00:00Z"),
				},
				Type:  QueryType_BETWEEN,
				Valid: true,
			},
		},
		"in": {
			given: "in:2009-10-10T23:00:00Z,2009-11-10T23:00:00Z,2009-12-10T23:00:00Z",
			expected: Timestamp{
				Values: []*timestamp.Timestamp{
					parseTimestamp(t, "2009-10-10T23:00:00Z"),
					parseTimestamp(t, "2009-11-10T23:00:00Z"),
					parseTimestamp(t, "2009-12-10T23:00:00Z"),
				},
				Type:  QueryType_IN,
				Valid: true,
			},
		},
	}

CasesLoop:
	for hint, c := range cases {
		got, err := ParseTimestamp(c.given)
		if err != nil {
			t.Errorf("%s: unexpected error: %s", hint, err.Error())
			continue CasesLoop
		}
		if got == nil {
			t.Errorf("%s: unexpected nil", hint)
			continue CasesLoop
		}
		if !reflect.DeepEqual(c.expected, *got) {
			t.Errorf("%s: wrong output,\nexpected:\n	%v\nbut got:\n	%v\n", hint, &c.expected, got)
		}
	}
}

func testInt64(t *testing.T, i *Int64, n, v bool, tp QueryType, values ...int64) {
	if i.Negation != n {
		t.Errorf("wrong negation, exiected %t but got %t", n, i.Negation)
	}

	if len(values) > 0 {
		if i.Value() != values[0] {
			t.Errorf("wrong first value, expected %d but got %d", values[0], i.Value())
		}
		if len(i.Values) == len(values) {
			for j, v := range values {
				if i.Values[j] != v {
					t.Errorf("%d: wrong value, expected %d but got %d", j, v, i.Values[j])
				}
			}
		} else {
			t.Errorf("wrong number of values, expected %d but got %d", len(values), len(i.Values))
		}
	}
	if i.Valid != v {
		t.Errorf("expected valid to be %t", v)
	}
	if i.Type != tp {
		t.Errorf("wrong type, expected %s but got %s", tp, i.Type)
	}
}

func TestBetweenFloat64(t *testing.T) {
	values := []float64{1111, 2222}
	testFloat64(t, BetweenFloat64(values[0], values[1]), false, true, QueryType_BETWEEN, values...)
}

func TestParseFloat64(t *testing.T) {
	cases := map[string]struct {
		given    string
		expected Float64
	}{
		"empty": {
			given:    "",
			expected: Float64{},
		},
		"null": {
			given: "null:",
			expected: Float64{
				Values: []float64{},
				Type:   QueryType_NULL,
				Valid:  true,
			},
		},
		"number": {
			given: "15.15",
			expected: Float64{
				Values: []float64{15.15},
				Type:   QueryType_EQUAL,
				Valid:  true,
			},
		},
		"not-null": {
			given: "nnull:",
			expected: Float64{
				Values:   []float64{},
				Type:     QueryType_NULL,
				Valid:    true,
				Negation: true,
			},
		},
		"equal": {
			given: "eq:123.15",
			expected: Float64{
				Values: []float64{123.15},
				Type:   QueryType_EQUAL,
				Valid:  true,
			},
		},
		"not-equal": {
			given: "neq:123.55555",
			expected: Float64{
				Values:   []float64{123.55555},
				Type:     QueryType_EQUAL,
				Valid:    true,
				Negation: true,
			},
		},
		"greater": {
			given: "gt:555",
			expected: Float64{
				Values: []float64{555.00},
				Type:   QueryType_GREATER,
				Valid:  true,
			},
		},
		"greater-equal": {
			given: "gte:666.666",
			expected: Float64{
				Values: []float64{666.666},
				Type:   QueryType_GREATER_EQUAL,
				Valid:  true,
			},
		},
		"lesser": {
			given: "lt:777.666",
			expected: Float64{
				Values: []float64{777.666},
				Type:   QueryType_LESS,
				Valid:  true,
			},
		},
		"lesser-equal": {
			given: "lte:888.666",
			expected: Float64{
				Values: []float64{888.666},
				Type:   QueryType_LESS_EQUAL,
				Valid:  true,
			},
		},
		"between": {
			given: "bw:111.666,222.666",
			expected: Float64{
				Values: []float64{111.666, 222.666},
				Type:   QueryType_BETWEEN,
				Valid:  true,
			},
		},
		"not-between": {
			given: "nbw:111.666,222",
			expected: Float64{
				Values:   []float64{111.666, 222},
				Type:     QueryType_BETWEEN,
				Valid:    true,
				Negation: true,
			},
		},
		"not-less": {
			given: "nlt:111.666",
			expected: Float64{
				Values:   []float64{111.666},
				Type:     QueryType_LESS,
				Valid:    true,
				Negation: true,
			},
		},
		"not-greater-than-or-equal": {
			given: "ngte:111.666",
			expected: Float64{
				Values:   []float64{111.666},
				Type:     QueryType_GREATER_EQUAL,
				Valid:    true,
				Negation: true,
			},
		},
		"not-less-than-or-equal": {
			given: "nlte:111.666",
			expected: Float64{
				Values:   []float64{111.666},
				Type:     QueryType_LESS_EQUAL,
				Valid:    true,
				Negation: true,
			},
		},
		"not-in": {
			given: "nin:111.666,222.444",
			expected: Float64{
				Values:   []float64{111.666, 222.444},
				Type:     QueryType_IN,
				Valid:    true,
				Negation: true,
			},
		},
		"contains": {
			given: "cts:111.666,222.444",
			expected: Float64{
				Values: []float64{111.666, 222.444},
				Type:   QueryType_CONTAINS,
				Valid:  true,
			},
		},
		"is-contained-by": {
			given: "icb:111.666,222.444",
			expected: Float64{
				Values: []float64{111.666, 222.444},
				Type:   QueryType_IS_CONTAINED_BY,
				Valid:  true,
			},
		},
		"overlap": {
			given: "ovl:111.666,222.444",
			expected: Float64{
				Values: []float64{111.666, 222.444},
				Type:   QueryType_OVERLAP,
				Valid:  true,
			},
		},
	}

CasesLoop:
	for hint, c := range cases {
		got, err := ParseFloat64(c.given)
		if err != nil {
			t.Errorf("%s: unexpected error: %s", hint, err.Error())
			continue CasesLoop
		}
		if got == nil {
			t.Errorf("%s: unexpected nil", hint)
			continue CasesLoop
		}
		if !reflect.DeepEqual(c.expected, *got) {
			t.Errorf("%s: wrong output,\nexpected:\n	%v\nbut got:\n	%v\n", hint, &c.expected, got)
		}
	}
}

func testFloat64(t *testing.T, f *Float64, n, v bool, tp QueryType, values ...float64) {
	if f.Negation != n {
		t.Errorf("wrong negation, exiected %t but got %t", n, f.Negation)
	}

	if len(values) > 0 {
		if f.Value() != values[0] {
			t.Errorf("wrong first value, expected %g but got %g", values[0], f.Value())
		}
		if len(f.Values) == len(values) {
			for j, v := range values {
				if f.Values[j] != v {
					t.Errorf("%d: wrong value, expected %g but got %g", j, v, f.Values[j])
				}
			}
		} else {
			t.Errorf("wrong number of values, expected %d but got %d", len(values), len(f.Values))
		}
	}
	if f.Valid != v {
		t.Errorf("expected valid to be %t", v)
	}
	if f.Type != tp {
		t.Errorf("wrong type, expected %s but got %s", tp, f.Type)
	}
}
