package aggro

import (
	"encoding/json"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

var (
	table = &Table{
		Fields: []Field{
			{"location", "string"},
			{"department", "string"},
			// {"name", "string"},
			{"salary", "number"},
			{"start_date", "datetime"},
		},
	}

	rows = []map[string]interface{}{
		{"location": "Auckland", "department": "Engineering", "salary": 120000, "start_date": "2016-01-31T22:00:00Z"},
		{"location": "Auckland", "department": "Engineering", "salary": 80000, "start_date": "2016-03-23T22:00:00Z"},
		{"location": "Auckland", "department": "Marketing", "salary": 90000, "start_date": "2016-01-31T22:00:00Z"},
		{"location": "Auckland", "department": "Marketing", "salary": 150000, "start_date": "2016-01-23T22:00:00Z"},
		{"location": "Wellington", "department": "Engineering", "salary": 120000, "start_date": "2016-01-23T22:00:00Z"},
		{"location": "Wellington", "department": "Engineering", "salary": 160000, "start_date": "2016-03-23T22:00:00Z"},
		{"location": "Wellington", "department": "Engineering", "salary": 120000, "start_date": "2016-02-02T22:00:00Z"},
	}
)

func TestBucketByString(t *testing.T) {
	RegisterTestingT(t)
	dataset := &Dataset{
		Table: table,
	}

	err := dataset.AddRows(rows...)
	if err != nil {
		t.Fatalf("Unexpected error creating dataset: %s", err.Error())
	}

	query := &Query{
		Metrics: []Metric{
			{Type: "max", Field: "salary"},
			{Type: "min", Field: "salary"},
			{Type: "mean", Field: "salary"},
			{Type: "median", Field: "salary"},
			{Type: "mode", Field: "salary"},
			{Type: "stdev", Field: "salary"},
			{Type: "cardinality", Field: "salary"},
			{Type: "count", Field: "salary"},
			{Type: "sum", Field: "salary"},
		},
		Bucket: &Bucket{
			Field: &Field{
				Name: "location",
				Type: "string",
			},
			Sort: &SortOptions{
				Type: "alphabetical",
			},
			Bucket: &Bucket{
				Field: &Field{
					Name: "department",
					Type: "string",
				},
				Sort: &SortOptions{
					Type: "alphabetical",
				},
			},
		},
	}

	results, err := dataset.Run(query)
	if err != nil {
		t.Fatalf("Unexpected error running query: %s", err.Error())
	}
	if results == nil {
		t.Fatalf("Unexpectedly got an empty resultset running query")
	}

	expected := Resultset{
		Buckets: []*ResultBucket{
			{
				Value: "Auckland",
				Buckets: []*ResultBucket{
					{
						Value: "Engineering",
						Metrics: map[string]interface{}{
							"salary:cardinality": 2,
							"salary:count":       2,
							"salary:max":         120000,
							"salary:mean":        100000,
							"salary:median":      100000,
							"salary:min":         80000,
							"salary:mode":        []float64{},
							"salary:stdev":       28284.2712474619,
							"salary:sum":         200000,
						},
					},
					{
						Value: "Marketing",
						Metrics: map[string]interface{}{
							"salary:cardinality": 2,
							"salary:count":       2,
							"salary:max":         150000,
							"salary:mean":        120000,
							"salary:median":      120000,
							"salary:min":         90000,
							"salary:mode":        []float64{},
							"salary:stdev":       42426.40687119285,
							"salary:sum":         240000,
						},
					},
				},
			},
			{
				Value: "Wellington",
				Buckets: []*ResultBucket{
					{
						Value: "Engineering",
						Metrics: map[string]interface{}{
							"salary:cardinality": 2,
							"salary:count":       3,
							"salary:max":         160000,
							"salary:mean":        133333.33333333334,
							"salary:median":      120000,
							"salary:min":         120000,
							"salary:mode": []float64{
								120000,
							},
							"salary:stdev": 23094.01076758503,
							"salary:sum":   400000,
						},
					},
				},
			},
		},
	}
	rm, _ := json.Marshal(*results)
	em, _ := json.Marshal(expected)
	Expect(rm).To(MatchJSON(em))
}

func TestBucketByStringDesc(t *testing.T) {
	RegisterTestingT(t)
	dataset := &Dataset{
		Table: table,
	}

	err := dataset.AddRows(rows...)
	if err != nil {
		t.Fatalf("Unexpected error creating dataset: %s", err.Error())
	}

	query := &Query{
		Metrics: []Metric{
			{Type: "max", Field: "salary"},
			{Type: "min", Field: "salary"},
			{Type: "mean", Field: "salary"},
			{Type: "median", Field: "salary"},
			{Type: "mode", Field: "salary"},
			{Type: "stdev", Field: "salary"},
			{Type: "cardinality", Field: "salary"},
			{Type: "count", Field: "salary"},
			{Type: "sum", Field: "salary"},
		},
		Bucket: &Bucket{
			Field: &Field{
				Name: "location",
				Type: "string",
			},
			Sort: &SortOptions{
				Type: "alphabetical",
				Desc: true,
			},
			Bucket: &Bucket{
				Field: &Field{
					Name: "department",
					Type: "string",
				},
				Sort: &SortOptions{
					Type: "alphabetical",
					Desc: true,
				},
			},
		},
	}

	results, err := dataset.Run(query)
	if err != nil {
		t.Fatalf("Unexpected error running query: %s", err.Error())
	}
	if results == nil {
		t.Fatalf("Unexpectedly got an empty resultset running query")
	}

	expected := Resultset{
		Buckets: []*ResultBucket{
			{
				Value: "Wellington",
				Buckets: []*ResultBucket{
					{
						Value: "Engineering",
						Metrics: map[string]interface{}{
							"salary:cardinality": 2,
							"salary:count":       3,
							"salary:max":         160000,
							"salary:mean":        133333.33333333334,
							"salary:median":      120000,
							"salary:min":         120000,
							"salary:mode": []float64{
								120000,
							},
							"salary:stdev": 23094.01076758503,
							"salary:sum":   400000,
						},
					},
				},
			},
			{
				Value: "Auckland",
				Buckets: []*ResultBucket{
					{
						Value: "Marketing",
						Metrics: map[string]interface{}{
							"salary:cardinality": 2,
							"salary:count":       2,
							"salary:max":         150000,
							"salary:mean":        120000,
							"salary:median":      120000,
							"salary:min":         90000,
							"salary:mode":        []float64{},
							"salary:stdev":       42426.40687119285,
							"salary:sum":         240000,
						},
					},
					{
						Value: "Engineering",
						Metrics: map[string]interface{}{
							"salary:cardinality": 2,
							"salary:count":       2,
							"salary:max":         120000,
							"salary:mean":        100000,
							"salary:median":      100000,
							"salary:min":         80000,
							"salary:mode":        []float64{},
							"salary:stdev":       28284.2712474619,
							"salary:sum":         200000,
						},
					},
				},
			},
		},
	}
	rm, _ := json.Marshal(*results)
	em, _ := json.Marshal(expected)
	Expect(rm).To(MatchJSON(em))
}

func TestBucketByDate(t *testing.T) {
	RegisterTestingT(t)
	dataset := &Dataset{
		Table: table,
	}

	err := dataset.AddRows(rows...)
	if err != nil {
		t.Fatalf("Unexpected error creating dataset: %s", err.Error())
	}

	start := time.Date(2015, 12, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2016, 4, 1, 0, 0, 0, 0, time.UTC)
	query := &Query{
		Metrics: []Metric{
			{Type: "mean", Field: "salary"},
			{Type: "max", Field: "salary"},
			{Type: "min", Field: "salary"},
		},
		Bucket: &Bucket{
			Field: &Field{
				Name: "location",
				Type: "string",
			},
			Sort: &SortOptions{
				Type: "alphabetical",
			},
			Bucket: &Bucket{
				Field: &Field{
					Name: "start_date",
					Type: "datetime",
				},
				DatetimeOptions: &DatetimeBucketOptions{
					Period:   Month,
					Start:    &start,
					End:      &end,
					Location: time.UTC,
				},
				Sort: &SortOptions{
					Type: "alphabetical",
				},
			},
		},
	}

	results, err := dataset.Run(query)
	if err != nil {
		t.Fatalf("Unexpected error running query: %s", err.Error())
	}
	if results == nil {
		t.Fatalf("Unexpectedly got an empty resultset running query")
	}

	expected := Resultset{
		Buckets: []*ResultBucket{
			{
				Value: "Auckland",
				Buckets: []*ResultBucket{
					{
						Value: "2015-12-01T00:00:00Z",
						Metrics: map[string]interface{}{
							"salary:mean": nil,
							"salary:max":  nil,
							"salary:min":  nil,
						},
					},
					{
						Value: "2016-01-01T00:00:00Z",
						Metrics: map[string]interface{}{
							"salary:max":  150000,
							"salary:mean": 120000,
							"salary:min":  90000,
						},
					},
					{
						Value: "2016-02-01T00:00:00Z",
						Metrics: map[string]interface{}{
							"salary:max":  nil,
							"salary:mean": nil,
							"salary:min":  nil,
						},
					},
					{
						Value: "2016-03-01T00:00:00Z",
						Metrics: map[string]interface{}{
							"salary:max":  80000,
							"salary:mean": 80000,
							"salary:min":  80000,
						},
					},
					{
						Value: "2016-04-01T00:00:00Z",
						Metrics: map[string]interface{}{
							"salary:max":  nil,
							"salary:mean": nil,
							"salary:min":  nil,
						},
					},
				},
			},
			{
				Value: "Wellington",
				Buckets: []*ResultBucket{
					{
						Value: "2015-12-01T00:00:00Z",
						Metrics: map[string]interface{}{
							"salary:max":  nil,
							"salary:mean": nil,
							"salary:min":  nil,
						},
					},
					{
						Value: "2016-01-01T00:00:00Z",
						Metrics: map[string]interface{}{
							"salary:max":  120000,
							"salary:mean": 120000,
							"salary:min":  120000,
						},
					},
					{
						Value: "2016-02-01T00:00:00Z",
						Metrics: map[string]interface{}{
							"salary:max":  120000,
							"salary:mean": 120000,
							"salary:min":  120000,
						},
					},
					{
						Value: "2016-03-01T00:00:00Z",
						Metrics: map[string]interface{}{
							"salary:max":  160000,
							"salary:mean": 160000,
							"salary:min":  160000,
						},
					},
					{
						Value: "2016-04-01T00:00:00Z",
						Metrics: map[string]interface{}{
							"salary:max":  nil,
							"salary:mean": nil,
							"salary:min":  nil,
						},
					},
				},
			},
		},
	}
	rm, _ := json.MarshalIndent(*results, "", "  ")
	em, _ := json.MarshalIndent(expected, "", "  ")
	Expect(rm).To(MatchJSON(em))
}

func TestBucketByDateTZ(t *testing.T) {
	RegisterTestingT(t)
	dataset := &Dataset{
		Table: table,
	}

	err := dataset.AddRows(rows...)
	if err != nil {
		t.Fatalf("Unexpected error creating dataset: %s", err.Error())
	}

	start := time.Date(2015, 12, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2016, 4, 1, 0, 0, 0, 0, time.UTC)
	loc, _ := time.LoadLocation("Pacific/Auckland")
	query := &Query{
		Metrics: []Metric{
			{Type: "mean", Field: "salary"},
			{Type: "max", Field: "salary"},
			{Type: "min", Field: "salary"},
		},
		Bucket: &Bucket{
			Field: &Field{
				Name: "location",
				Type: "string",
			},
			Sort: &SortOptions{
				Type: "alphabetical",
			},
			Bucket: &Bucket{
				Field: &Field{
					Name: "start_date",
					Type: "datetime",
				},
				DatetimeOptions: &DatetimeBucketOptions{
					Period:   Month,
					Start:    &start,
					End:      &end,
					Location: loc,
				},
				Sort: &SortOptions{
					Type: "alphabetical",
				},
			},
		},
	}

	results, err := dataset.Run(query)
	if err != nil {
		t.Fatalf("Unexpected error running query: %s", err.Error())
	}
	if results == nil {
		t.Fatalf("Unexpectedly got an empty resultset running query")
	}

	expected := Resultset{
		Buckets: []*ResultBucket{
			{
				Value: "Auckland",
				Buckets: []*ResultBucket{
					{
						Value: "2015-12-01T00:00:00+13:00",
						Metrics: map[string]interface{}{
							"salary:max":  nil,
							"salary:mean": nil,
							"salary:min":  nil,
						},
					},
					{
						Value: "2016-01-01T00:00:00+13:00",
						Metrics: map[string]interface{}{
							"salary:max":  150000,
							"salary:mean": 150000,
							"salary:min":  150000,
						},
					},
					{
						Value: "2016-02-01T00:00:00+13:00",
						Metrics: map[string]interface{}{
							"salary:max":  120000,
							"salary:mean": 105000,
							"salary:min":  90000,
						},
					},
					{
						Value: "2016-03-01T00:00:00+13:00",
						Metrics: map[string]interface{}{
							"salary:max":  80000,
							"salary:mean": 80000,
							"salary:min":  80000,
						},
					},
					{
						Value: "2016-04-01T00:00:00+13:00",
						Metrics: map[string]interface{}{
							"salary:max":  nil,
							"salary:mean": nil,
							"salary:min":  nil,
						},
					},
				},
			},
			{
				Value: "Wellington",
				Buckets: []*ResultBucket{
					{
						Value: "2015-12-01T00:00:00+13:00",
						Metrics: map[string]interface{}{
							"salary:max":  nil,
							"salary:mean": nil,
							"salary:min":  nil,
						},
					},
					{
						Value: "2016-01-01T00:00:00+13:00",
						Metrics: map[string]interface{}{
							"salary:max":  120000,
							"salary:mean": 120000,
							"salary:min":  120000,
						},
					},
					{
						Value: "2016-02-01T00:00:00+13:00",
						Metrics: map[string]interface{}{
							"salary:max":  120000,
							"salary:mean": 120000,
							"salary:min":  120000,
						},
					},
					{
						Value: "2016-03-01T00:00:00+13:00",
						Metrics: map[string]interface{}{
							"salary:mean": 160000,
							"salary:max":  160000,
							"salary:min":  160000,
						},
					},
					{
						Value: "2016-04-01T00:00:00+13:00",
						Metrics: map[string]interface{}{
							"salary:max":  nil,
							"salary:mean": nil,
							"salary:min":  nil,
						},
					},
				},
			},
		},
	}
	rm, _ := json.MarshalIndent(*results, "", "  ")
	em, _ := json.MarshalIndent(expected, "", "  ")
	Expect(rm).To(MatchJSON(em))
}

func TestBucketByRange(t *testing.T) {
	RegisterTestingT(t)
	dataset := &Dataset{
		Table: table,
	}

	err := dataset.AddRows(rows...)
	if err != nil {
		t.Fatalf("Unexpected error creating dataset: %s", err.Error())
	}

	query := &Query{
		Metrics: []Metric{
			{Type: "count", Field: "salary"},
		},
		Bucket: &Bucket{
			Field: &Field{
				Name: "location",
				Type: "string",
			},
			Sort: &SortOptions{
				Type: "alphabetical",
			},
			Bucket: &Bucket{
				Field: &Field{
					Name: "salary",
					Type: "number",
				},
				RangeOptions: &RangeBucketOptions{
					Period: []interface{}{
						20000,
						50000,
						100000,
						150000,
						200000,
						300000,
					},
				},
				Sort: &SortOptions{
					Type: "numerical",
					Desc: false,
				},
			},
		},
	}

	results, err := dataset.Run(query)
	if err != nil {
		t.Fatalf("Unexpected error running query: %s", err.Error())
	}
	if results == nil {
		t.Fatalf("Unexpectedly got an empty resultset running query")
	}

	expected := Resultset{
		Buckets: []*ResultBucket{
			{
				Value: "Auckland",
				Buckets: []*ResultBucket{
					{
						Value:   "20000",
						Metrics: nil,
					},
					{
						Value:   "50000",
						Metrics: nil,
					},
					{
						Value: "100000",
						Metrics: map[string]interface{}{
							"salary:count": 2,
						},
					},
					{
						Value: "150000",
						Metrics: map[string]interface{}{
							"salary:count": 2,
						},
					},
					{
						Value:   "200000",
						Metrics: nil,
					},
					{
						Value:   "300000",
						Metrics: nil,
					},
				},
			},
			{
				Value: "Wellington",
				Buckets: []*ResultBucket{
					{
						Value:   "20000",
						Metrics: nil,
					},
					{
						Value:   "50000",
						Metrics: nil,
					},
					{
						Value:   "100000",
						Metrics: nil,
					},
					{
						Value: "150000",
						Metrics: map[string]interface{}{
							"salary:count": 2,
						},
					},
					{
						Value: "200000",
						Metrics: map[string]interface{}{
							"salary:count": 1,
						},
					},
					{
						Value:   "300000",
						Metrics: nil,
					},
				},
			},
		},
	}
	rm, _ := json.Marshal(*results)
	em, _ := json.Marshal(expected)
	Expect(rm).To(MatchJSON(em))
}
