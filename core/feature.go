package core

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/cespare/xxhash"
)

// IsCategorical returns true if categorical
func (k Feature_Kind) IsCategorical() bool { return k == Feature_CATEGORICAL }

// IsNumerical returns true if numerical
func (k Feature_Kind) IsNumerical() bool { return k == Feature_NUMERICAL }

// --------------------------------------------------------------------

// NewNumericalFeature initialises a new numerical feature.
func NewNumericalFeature(name string) *Feature {
	return &Feature{
		Name: name,
		Kind: Feature_NUMERICAL,
	}
}

// NewCategoricalFeature initialises a new categorical feature with a vocabulary.
func NewCategoricalFeature(name string, vocabulary []string) *Feature {
	return &Feature{
		Name:       name,
		Kind:       Feature_CATEGORICAL,
		Vocabulary: vocabulary,
	}
}

// NewCategoricalFeatureExpandable initialises a new categorical feature with an exandable vocabulary.
func NewCategoricalFeatureExpandable(name string, vocabulary []string) *Feature {
	return &Feature{
		Name:       name,
		Kind:       Feature_CATEGORICAL,
		Strategy:   Feature_EXPANDABLE,
		Vocabulary: vocabulary,
	}
}

// NewCategoricalFeatureIdentity initialises a new categorical feature with identity.
func NewCategoricalFeatureIdentity(name string) *Feature {
	return &Feature{
		Name:     name,
		Kind:     Feature_CATEGORICAL,
		Strategy: Feature_IDENTITY,
	}
}

// NewCategoricalFeatureHashBuckets initialises a new categorical feature with a number of hash buckets.
// Values are converted to integers via:
//   HASH(value) % numBuckets
func NewCategoricalFeatureHashBuckets(name string, numBuckets uint32) *Feature {
	return &Feature{
		Name:        name,
		Kind:        Feature_CATEGORICAL,
		HashBuckets: numBuckets,
	}
}

// NewCategoricalFeatureVocabularyHashBuckets initialises a new categorical feature with a fixed vocabulary
// and hash slots for unknown values.
// By default values are looked up from vocabulary, but unknown values are converted via:
//   HASH(value) % numBuckets + len(vocabulary)
func NewCategoricalFeatureVocabularyHashBuckets(name string, vocabulary []string, numBuckets uint32) *Feature {
	return &Feature{
		Name:        name,
		Kind:        Feature_CATEGORICAL,
		Vocabulary:  vocabulary,
		HashBuckets: numBuckets,
	}
}

// NumCategories returns the total number of categories associated
// with this feature. Will return -1 if unknown.
func (f *Feature) NumCategories() int {
	if f.Kind != Feature_CATEGORICAL {
		return 0
	}
	if f.Strategy == Feature_IDENTITY {
		return -1
	}
	return int(f.HashBuckets) + len(f.Vocabulary)
}

// Category will extract the categorical value from the example.
// Will return NoCategory if value is unknown/unobtainable.
func (f *Feature) Category(x Example) Category {
	if f.Kind != Feature_CATEGORICAL {
		return NoCategory
	}

	return f.CategoryOf(x.GetExampleValue(f.Name))
}

// CategoryOf attempts to retrieve to category for the given value.
// It may return NoCategory.
func (f *Feature) CategoryOf(v interface{}) Category {
	if v == nil || f.Kind != Feature_CATEGORICAL {
		return NoCategory
	}

	if f.Strategy == Feature_IDENTITY {
		return categorize(v)
	}

	s := stringify(v)
	for i, vv := range f.Vocabulary {
		if s == vv {
			return Category(i)
		}
	}

	if f.HashBuckets != 0 {
		return Category(xxhash.Sum64String(s)%uint64(f.HashBuckets)) + Category(len(f.Vocabulary))
	}

	if f.Strategy == Feature_EXPANDABLE {
		f.Vocabulary = append(f.Vocabulary, s)
		return Category(len(f.Vocabulary) - 1)
	}

	return NoCategory
}

// ValueOf returns the string value of a category.
// Will return "?" if value is unknown/unobtainable.
func (f *Feature) ValueOf(cat Category) string {
	if f.Kind != Feature_CATEGORICAL || !IsCat(cat) {
		return "?"
	}

	pos := int(cat)
	if f.Strategy == Feature_IDENTITY {
		return strconv.Itoa(pos)
	}

	if pos < len(f.Vocabulary) {
		return f.Vocabulary[pos]
	}

	if n := pos - len(f.Vocabulary); n < int(f.HashBuckets) {
		return "#" + strconv.Itoa(n)
	}

	return "?"
}

// Number will extract the numeric value from the example.
// Will return NaN if value is unknown/unobtainable.
func (f *Feature) Number(x Example) float64 {
	if f.Kind != Feature_NUMERICAL {
		return math.NaN()
	}
	return numerify(x.GetExampleValue(f.Name))
}

func stringify(v interface{}) string {
	switch w := v.(type) {
	case string:
		return w
	case []byte:
		return string(w)
	default:
		return fmt.Sprintf("%v", w)
	}
}

func categorize(v interface{}) Category {

	// fast path
	switch n := v.(type) {
	case []byte:
		return categorize(string(n))
	case string:
		if nn, err := strconv.ParseInt(n, 10, 64); err == nil {
			return Category(nn)
		}
		return NoCategory
	case Category:
		return n
	case int:
		return Category(n)
	case int8:
		return Category(n)
	case int16:
		return Category(n)
	case int32:
		return Category(n)
	case int64:
		return Category(n)
	case uint:
		return Category(n)
	case uint8:
		return Category(n)
	case uint16:
		return Category(n)
	case uint32:
		return Category(n)
	case uint64:
		return Category(n)
	case float32:
		return Category(n)
	case float64:
		return Category(n)
	case bool:
		if n {
			return 1
		}
		return 0
	}

	// slow path
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int:
		return categorize(rv.Int())
	case reflect.String:
		return categorize(rv.String())
	case reflect.Float32, reflect.Float64:
		return categorize(rv.Float())
	case reflect.Bool:
		return categorize(rv.Bool())
	case reflect.Ptr:
		return categorize(rv.Elem().Interface())
	}

	return NoCategory
}

func numerify(v interface{}) float64 {

	// fast path
	switch n := v.(type) {
	case []byte:
		return numerify(string(n))
	case string:
		if nn, err := strconv.ParseFloat(n, 64); err == nil {
			return nn
		}
		return math.NaN()
	case int:
		return float64(n)
	case int8:
		return float64(n)
	case int16:
		return float64(n)
	case int32:
		return float64(n)
	case int64:
		return float64(n)
	case uint:
		return float64(n)
	case uint8:
		return float64(n)
	case uint16:
		return float64(n)
	case uint32:
		return float64(n)
	case uint64:
		return float64(n)
	case float32:
		return float64(n)
	case float64:
		return n
	}

	// slow path
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int:
		return numerify(rv.Int())
	case reflect.String:
		return numerify(rv.String())
	case reflect.Float32, reflect.Float64:
		return numerify(rv.Float())
	case reflect.Ptr:
		return numerify(rv.Elem().Interface())
	}

	return math.NaN()
}
