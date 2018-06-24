package aggro

import (
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// MetricDelimeter is a string used to separate metric field's from names.
var MetricDelimeter = ":"

type queryProcessor struct {
	dataset     *Dataset
	query       *Query
	tipBuckets  map[*ResultBucket]bool
	measurables []*[]Cell
	err         error
	results     *Resultset
	composition []interface{}
	hasDatetime bool
	hasRange    bool
}

func (p *queryProcessor) Run() (*Resultset, error) {
	p.prepare()
	p.aggregate()
	p.measure()
	return p.results, p.err
}

func (p *queryProcessor) prepare() {
	if p.err != nil {
		return
	}

	// Initialise the root & tip buckets, and full bucket lookup.
	p.tipBuckets = map[*ResultBucket]bool{}
}

// aggregate is responsible for sorting the dataset's rows into buckets.
func (p *queryProcessor) aggregate() {
	if p.err != nil {
		return
	}
	if p.query.Bucket == nil {
		p.err = errors.New("Query has no root bucket")
		return
	}
	// Loop over each row, adding all nest query buckets to the value buckets.
	buckets := map[string]*ResultBucket{}
	for i, row := range p.dataset.Rows {
		buckets = p.recurse(0, i, row, p.query.Bucket, buckets)
	}

	buckets = p.fillDatetimeGaps(buckets)

	buckets = p.fillRangeGaps(buckets)

	p.results = &Resultset{
		Buckets:     p.sort(buckets),
		Composition: p.composition,
	}
}

func (p *queryProcessor) recurse(depth, index int, row map[string]Cell, aggregate *Bucket, results map[string]*ResultBucket) map[string]*ResultBucket {
	// If there's no aggregate, we're done.
	if aggregate == nil {
		return results
	}

	// Ensure we have the details required to bucket on.
	if aggregate.Field.Type == fieldTypeDatetime && aggregate.DatetimeOptions == nil {
		p.err = errors.New("Bucketing by datetime without DatetimeOptions set")
		return results
	}

	// Grab the cell that we're aggregating on.
	cell := row[aggregate.Field.Name]

	// Handle nil cell.
	if cell == nil {
		return results
	}

	// And grab the underlying aggregatable string value.
	value := ""

	switch tCell := cell.(type) {
	case *StringCell:
		// String Cell's are easy, it's just the value.
		value = tCell.value
		p.composition = append(p.composition, tCell.data)
	case *DatetimeCell:
		p.hasDatetime = true
		// Datetime Cell's are a bit more complicated, and need the period value.
		value, p.err = tCell.ValueForPeriod(aggregate.DatetimeOptions.Period, aggregate.DatetimeOptions.Location)
		if p.err != nil {
			return results
		}
		p.composition = append(p.composition, tCell.data)
	case *NumberCell:
		if aggregate.RangeOptions != nil {
			p.hasRange = true
			value, p.err = tCell.ValueForPeriod(aggregate.RangeOptions.Period)
			if p.err != nil {
				return results
			}
			p.composition = append(p.composition, tCell.data)
		} else {
			p.err = fmt.Errorf("Non aggregatable cell found without RangeOptions at depth %d, index %d", depth, index)
		}
	default:
		p.err = fmt.Errorf("Non aggregatable cell found at depth %d, index %d", depth, index)
		return results
	}

	// Ensure we have a result bucket for this value, making one if we don't.
	bucket := ensureValueBucket(results, value)

	// If there's no next bucket, we're at the deepest point. Add data to measure.
	if aggregate.Bucket == nil {
		bucket.sourceRows = append(bucket.sourceRows, row)
		p.tipBuckets[bucket] = true
	}

	// Bump depth and recurse to next level, passing in the children as the results.
	depth++
	bucket.bucketLookup = p.recurse(depth, index, row, aggregate.Bucket, bucket.bucketLookup)

	// Update the current results bucket with the new values, then return.
	results[value] = bucket
	return results
}

func ensureValueBucket(results map[string]*ResultBucket, value string) *ResultBucket {
	bucket := results[value]
	if bucket == nil {
		bucket = &ResultBucket{
			Value:        value,
			bucketLookup: map[string]*ResultBucket{},
		}
	}
	return bucket
}

func (p *queryProcessor) fillDatetimeGaps(results map[string]*ResultBucket) map[string]*ResultBucket {
	if !p.hasDatetime {
		return results
	}
	return p.fillBucketDatetimeGaps(p.query.Bucket, results)
}

func (p *queryProcessor) fillBucketDatetimeGaps(bucket *Bucket, results map[string]*ResultBucket) map[string]*ResultBucket {
	if bucket == nil || len(results) < 0 {
		return results
	}
	if (bucket.Field.Type == fieldTypeDatetime) || p.hasRange {
		// Get the max and min values.
		var min, max *string
		// Set the min to the start if there is one.
		if bucket.DatetimeOptions.Start != nil {
			var start string
			start, p.err = datetimeValueForPeriod(
				bucket.DatetimeOptions.Start,
				bucket.DatetimeOptions.Period,
				bucket.DatetimeOptions.Location,
			)
			if p.err != nil {
				return results
			}
			min = &start
		}
		// Set the max to the end if there is one.
		if bucket.DatetimeOptions.End != nil {
			var end string
			end, p.err = datetimeValueForPeriod(
				bucket.DatetimeOptions.End,
				bucket.DatetimeOptions.Period,
				bucket.DatetimeOptions.Location,
			)
			if p.err != nil {
				return results
			}
			max = &end
		}
		// Now extend the start and end depending on the values in the results.
		for key := range results {
			value := key
			// Min being nil means max is too. Set them to the first result.
			if min == nil {
				min = &value
				max = &value
			} else {
				if *min > value {
					min = &value
				}
				if *max < value {
					max = &value
				}
			}
		}
		// No need to do anything if we only have a single bucket length.
		if *min == *max {
			return results
		}

		loopValue := *min
		var loopDate time.Time
		loopDate, p.err = time.Parse(time.RFC3339, loopValue)
		if p.err != nil {
			return results
		}
		// Now loop until we hit the max point, ensuring each period exists.
		for loopValue <= *max {
			// Make sure this period exists.
			results[loopValue] = ensureValueBucket(results, loopValue)
			if bucket.Bucket == nil {
				p.tipBuckets[results[loopValue]] = true
			}

			// Now bump the date up one period, and loop.
			date, err := datetimeAddPeriod(&loopDate, bucket.DatetimeOptions.Period)
			if err != nil {
				p.err = err
				return results
			}

			loopDate = *date
			loopValue, p.err = datetimeValueForPeriod(
				&loopDate,
				bucket.DatetimeOptions.Period,
				bucket.DatetimeOptions.Location,
			)
			if p.err != nil {
				return results
			}
		}
	}

	// Now recurse into any children result sets.
	for _, result := range results {
		result.bucketLookup = p.fillBucketDatetimeGaps(bucket.Bucket, result.bucketLookup)
	}

	return results
}

func (p *queryProcessor) sort(results map[string]*ResultBucket) []*ResultBucket {
	return sortMap(p.query.Bucket, results)
}

func (p *queryProcessor) fillRangeGaps(results map[string]*ResultBucket) map[string]*ResultBucket {
	if !p.hasRange {
		return results
	}
	return p.fillBucketRangeGaps(p.query.Bucket, results)
}

func (p *queryProcessor) fillBucketRangeGaps(bucket *Bucket, results map[string]*ResultBucket) map[string]*ResultBucket {
	if bucket == nil || len(results) < 0 {
		return results
	}

	if (bucket.Field.Type == fieldTypeNumber) && (bucket.RangeOptions.Period != nil) {
		for _, p := range bucket.RangeOptions.Period {

			var v float64
			switch i := p.(type) {
			case float64:
				v = i
			case float32:
				v = float64(i)
			case int32:
				v = float64(i)
			case int64:
				v = float64(i)
			case int:
				v = float64(i)
			}

			// Make sure this period exists.
			index := decimal.NewFromFloat(v)
			results[index.String()] = ensureValueBucket(results, index.String())
		}
	}

	// Now recurse into any children result sets.
	for _, result := range results {
		result.bucketLookup = p.fillBucketRangeGaps(bucket.Bucket, result.bucketLookup)
	}

	return results
}

func (p *queryProcessor) measure() {
	if p.err != nil {
		return
	}

	// We only add metrics for the tip buckets, i.e. the deepest nesting.
	for bucket := range p.tipBuckets {
		// Create measurers for each of the metrics, then feed data into them.
		bucket.Metrics = map[string]interface{}{}
		var m measurer

		for i := range p.query.Metrics {
			metric := &p.query.Metrics[i]
			// Create a measurer.
			m, p.err = metric.measurer()
			if p.err != nil {
				return
			}
			// Now add all of the data to the measurer.
			for j := range bucket.sourceRows {
				row := bucket.sourceRows[j]

				// Check the field is of a metricable type.
				if !row[metric.Field].IsMetricable(m) {
					p.err = fmt.Errorf("Non metricable cell found (`%s:%s`)", metric.Field, metric.Type)
					return
				}

				m.AddDatum(row[metric.Field].MeasurableCell().Value())
			}

			// And then push the result into the metric resultset.
			bucket.Metrics[metric.Field+MetricDelimeter+metric.Type] = m.Result()
		}
	}
}
