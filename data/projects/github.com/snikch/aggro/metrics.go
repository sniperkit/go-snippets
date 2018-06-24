package aggro

import (
	"fmt"
	"math"
	"sort"

	"github.com/shopspring/decimal"
)

// Metric represents a type of measurement to be applied to our dataset.
type Metric struct {
	Type  string
	Field string
}

func (m *Metric) measurer() (measurer, error) {
	switch m.Type {
	case "mean":
		return &mean{}, nil
	case "median":
		return &median{}, nil
	case "mode":
		return &mode{}, nil
	case "min":
		return &min{}, nil
	case "max":
		return &max{}, nil
	case "cardinality":
		return &cardinality{}, nil
	case "sum":
		return &sum{}, nil
	case "stdev":
		return &stdev{}, nil
	case "count":
		return &valueCount{}, nil
	default:
		return nil, fmt.Errorf("Unknown metric: %s", m.Type)
	}
}

type measurer interface {
	AddDatum(interface{})
	Result() interface{}
}

// Mean
// Your standard average. Sum all values and divide by the number of values.
type mean struct {
	count int
	sum   decimal.Decimal
}

func (a *mean) AddDatum(datum interface{}) {
	// Cast to *decimal.Decimal.
	amount := datum.(*decimal.Decimal)

	// Increase our count.
	a.count++

	// Add our value to our existing sum.
	a.sum = a.sum.Add(*amount)
}

func (a *mean) Result() interface{} {
	if a.count == 0 {
		return nil
	}

	// Divide our sum by the count.
	result, _ := a.sum.Div(decimal.NewFromFloat(float64(a.count))).Float64()
	return result
}

// Median
// Dataset should be in numerical order. Median is the middle value of our dataset.
// If there is no middle value (due to dataset having even number of values) then
// the median is the mean (average) of the middle two values.
type median struct {
	list []decimal.Decimal
}

func (a *median) AddDatum(datum interface{}) {
	// Cast to *decimal.Decimal.
	amount := datum.(*decimal.Decimal)

	// Append value to median slice.
	a.list = append(a.list, *amount)
}

func (a *median) Result() interface{} {
	if len(a.list) == 0 {
		return nil
	}

	// Our result.
	var median float64

	// Sort our list in numerical order.
	sort.Sort(decimalSortNumerical(a.list))

	// Decimal size of our dataset.
	size := decimal.NewFromFloat(float64(len(a.list)))

	// Find the middle of our dataset.
	middle := size.Div(decimal.New(2, 0))

	// Is our dataset length even? If so, we don't have a correct middle value.
	// In this case, take the middle two values of our dataset and determine the
	// mean of them.
	if size.Mod(decimal.New(2, 0)).Equals(decimal.New(0, 1)) {
		// Find value: middle - 1.
		prev := a.list[middle.Sub(decimal.New(1, 0)).IntPart()]

		// Add two middle values and divide by 2.
		median, _ = a.list[middle.IntPart()].Add(prev).Div(decimal.New(2, 0)).Float64()
	} else {
		// Simply return middle value.
		median, _ = a.list[middle.IntPart()].Float64()
	}
	return median
}

// Mode
// Mode is the value(s) that occur most often within the dataset. If no values
// are repeated (or all values are repeated), then the dataset has no mode.
type mode struct {
	count int
	list  []decimal.Decimal
	value decimal.Decimal
}

func (a *mode) AddDatum(datum interface{}) {
	// Cast to *decimal.Decimal.
	amount := datum.(*decimal.Decimal)

	// Append value to mode slice.
	a.list = append(a.list, *amount)

	// Increase our count.
	a.count++
}

func (a *mode) Result() interface{} {
	if len(a.list) == 0 {
		return nil
	}

	// Results slice.
	modes := []float64{}

	// 'tip' represents our highest frequency count across our entire dataset. A
	// dataset with a tip of '1' means no repeated values were found.
	tip := 0

	// Range our values building a frequency map. This represents each value and
	// the number of times it appears in our dataset.
	freq := make(map[float64]int, len(a.list))
	for _, val := range a.list {
		value, _ := val.Float64()
		freq[value]++
		if freq[value] > tip {
			tip = freq[value]
		}
	}

	// Range our frequency map, checking if our values count matches our tip. If so
	// we have a mode!
	for val, c := range freq {
		if c == tip {
			modes = append(modes, val)
		}
	}

	// If tip is 1 (no repeating values found), or length of resulting modes slice
	// matches our dataset, then return no mode (empty).
	if tip == 1 || len(modes) == len(a.list) {
		modes = []float64{}
	}
	return modes
}

// Min
// Min is the smallest value within the dataset.
type min struct {
	amount *decimal.Decimal
}

func (a *min) AddDatum(datum interface{}) {
	// Cast to *decimal.Decimal.
	amount := datum.(*decimal.Decimal)
	if a.amount == nil {
		a.amount = amount
	}

	// If value is < existing min.amount, assign as min.amount.
	if (a.amount).Cmp(*amount) > -1 {
		a.amount = amount
	}
}

func (a *min) Result() interface{} {
	if a.amount == nil {
		return nil
	}
	result, _ := a.amount.Float64()
	return result
}

// Max
// Max is the largest value within the dataset.
type max struct {
	amount *decimal.Decimal
}

func (a *max) AddDatum(datum interface{}) {
	// Cast to *decimal.Decimal.
	amount := datum.(*decimal.Decimal)
	if a.amount == nil {
		a.amount = amount
	}

	// If value is > existing max.amount, assign as max.amount.
	if (a.amount).Cmp(*amount) < 0 {
		a.amount = amount
	}
}

func (a *max) Result() interface{} {
	if a.amount == nil {
		return nil
	}
	result, _ := a.amount.Float64()
	return result
}

// Cardinality
// Cardinality is a count of unique values in our dataset.
type cardinality struct {
	values map[interface{}]int
}

func (a *cardinality) AddDatum(datum interface{}) {
	if a.values == nil {
		a.values = map[interface{}]int{}
	}

	// Track frequency of our values within the dataset.
	switch t := datum.(type) {
	case *decimal.Decimal:
		floatVal, _ := t.Float64()
		a.values[floatVal]++
	case string:
		a.values[t]++
	}
}

func (a *cardinality) Result() interface{} {
	return len(a.values)
}

// Value Count
// valueCount is the total number of values in the dataset.
type valueCount struct {
	size int
}

func (a *valueCount) AddDatum(datum interface{}) {
	a.size++
}

func (a *valueCount) Result() interface{} {
	return a.size
}

// Sum
// Sum is all dataset values added together.
type sum struct {
	sum decimal.Decimal
}

func (a *sum) AddDatum(datum interface{}) {
	amount := datum.(*decimal.Decimal)

	// Add our value to existing sum.
	a.sum = a.sum.Add(*amount)
}

func (a *sum) Result() interface{} {
	result, _ := a.sum.Float64()
	return result
}

// Standard deviation
// Standard deviation is a representation of how spread out values in the dataset are.
// It's calculated as square root of 'variance'. Variance is the average of the
// squared differences from the mean.
//
// 1) Determine mean.
// 2) Then for each value in the dataset, subtract the mean and square the result.
// 3) Calculate the mean of those squared differences (variance).
// 4) Return square root of variance.
type stdev struct {
	count int
	sum   decimal.Decimal
	list  []decimal.Decimal
}

func (a *stdev) AddDatum(datum interface{}) {
	// Cast to *decimal.Decimal.
	amount := datum.(*decimal.Decimal)

	// Increase our count.
	a.count++

	// Add our value to existing sum.
	a.sum = a.sum.Add(*amount)

	// Append value to stdev slice.
	a.list = append(a.list, *amount)
}

func (a *stdev) Result() interface{} {
	// stdev requires two or more rows to work with.
	if a.count < 2 {
		return nil
	}

	// 1) Determine the mean (avg).
	mean, _ := a.sum.Div(decimal.NewFromFloat(float64(a.count))).Float64()

	// 2) Ranging our dataset, subtract mean from value and square the result.
	total := 0.0
	for _, number := range a.list {
		val, _ := number.Float64()
		total += math.Pow(val-mean, 2)
	}

	// 3) Calculate the variance (mean of our squared results).
	variance := total / float64(len(a.list)-1)

	// 4) Square the result.
	return math.Sqrt(variance)
}

// Sorters
type decimalSortNumerical []decimal.Decimal

func (s decimalSortNumerical) Len() int      { return len(s) }
func (s decimalSortNumerical) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s decimalSortNumerical) Less(i, j int) bool {
	a, _ := s[i].Float64()
	b, _ := s[j].Float64()
	return a < b
}
