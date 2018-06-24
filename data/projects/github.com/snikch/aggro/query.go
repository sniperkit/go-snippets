package aggro

import "time"

// Query represents the bucketing and aggregate metrics that should be run.
type Query struct {
	Bucket  *Bucket
	Metrics []Metric
}

// Bucket defines how to compare and group data which is then aggregated on.
type Bucket struct {
	Bucket          *Bucket
	Field           *Field
	DatetimeOptions *DatetimeBucketOptions
	Sort            *SortOptions
	RangeOptions    *RangeBucketOptions
}

// SortOptions represent how this Bucket should be sorted.
type SortOptions struct {
	Type   string
	Metric string
	Desc   bool
}

// DatetimeBucketOptions provides additional configuration for datetime bucketing.
type DatetimeBucketOptions struct {
	// Start will, if provided, ensure buckets start at this date.
	Start *time.Time
	// End will, if provided, ensure buckets continue to this date.
	End *time.Time
	// What interval period are the results to be bucketed at.
	Period DatetimePeriod
	// Datetimes should be bucketed based on the date in this location.
	Location *time.Location
}

// RangeBucketOptions provides additional configuration for custom range bucketing.
type RangeBucketOptions struct {
	Period []interface{}
}
