# aggro 😡

[![GoDoc](https://godoc.org/github.com/snikch/aggro?status.svg)](https://godoc.org/github.com/snikch/aggro)
[![Go Report Card](https://goreportcard.com/badge/github.com/snikch/aggro)](https://goreportcard.com/report/github.com/snikch/aggro)

In memory dataset bucketing and metrics, inspired by Elasticsearch aggregations

Installation
-----------------

`go get -u github.com/snikch/aggro`

Example
-------------

Given a dataset...

```go
rows := []map[string]interface{}{
		{"location": "Auckland", "department": "Engineering", "salary": 120000, "start_date": "2016-01-31T22:00:00Z"},
		{"location": "Auckland", "department": "Engineering", "salary": 80000, "start_date": "2016-03-23T22:00:00Z"},
		{"location": "Auckland", "department": "Marketing", "salary": 90000, "start_date": "2016-01-31T22:00:00Z"},
		{"location": "Auckland", "department": "Marketing", "salary": 150000, "start_date": "2016-01-23T22:00:00Z"},
		{"location": "Wellington", "department": "Engineering", "salary": 120000, "start_date": "2016-01-23T22:00:00Z"},
		{"location": "Wellington", "department": "Engineering", "salary": 160000, "start_date": "2016-03-23T22:00:00Z"},
		{"location": "Wellington", "department": "Engineering", "salary": 120000, "start_date": "2016-02-02T22:00:00Z"},
	}
```

Initialize aggro, build aggregations and run...

```go
// Build a dataset that contains a *Table representing your data.
dataset := &Dataset{
	Table: &Table{
		Fields: []Field{
			{"location", "string"},
			{"department", "string"},
			{"salary", "number"},
			{"start_date", "datetime"},
		},
	},
}

// Add our rows to our dataset.
err := dataset.AddRows(rows...)
if err != nil {
	return err
}

// Build our query specifying preferred metrics and bucket composition.
query := &Query{
    Metrics: []Metric{
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
                Name: "department",
                Type: "string",
            },
            Sort: &SortOptions{
                Type: "alphabetical",
            },
        },
    },
}

// Run it.
results, err := dataset.Run(query)
if err != nil {
	return err
}
```

Find a list of available [measurers here](https://github.com/snikch/aggro/blob/master/metrics.go#L15)
