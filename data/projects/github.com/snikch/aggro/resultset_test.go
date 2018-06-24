package aggro

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/gomega"
)

func TestFlattenResultSetSimple(t *testing.T) {
	RegisterTestingT(t)
	expected := ResultTable{
		Rows: [][]map[string]interface{}{
			{
				{
					"salary:max": nil,
				},
				{
					"salary:max": 150000,
				},
				{
					"salary:max": 120000,
				},
			},
			{
				{
					"salary:max": nil,
				},
				{
					"salary:max": 120000,
				},
				{
					"salary:max": 220000,
				},
			},
		},
		RowTitles: [][]string{
			{"Auckland"},
			{"Wellington"},
		},
		ColumnTitles: [][]string{
			{"2015-12-01T00:00:00+13:00"},
			{"2016-01-01T00:00:00+13:00"},
			{"2016-02-01T00:00:00+13:00"},
		},
	}
	input := &Resultset{
		Buckets: []*ResultBucket{
			{
				Value: "Auckland",
				Buckets: []*ResultBucket{
					{
						Value: "2015-12-01T00:00:00+13:00",
						// Buckets: []*ResultBucket{},
						Metrics: map[string]interface{}{
							"salary:max": nil,
						},
					},
					{
						Value: "2016-01-01T00:00:00+13:00",
						// Buckets: []*ResultBucket{},
						Metrics: map[string]interface{}{
							"salary:max": 150000,
						},
					},
					{
						Value:   "2016-02-01T00:00:00+13:00",
						Buckets: []*ResultBucket{},
						Metrics: map[string]interface{}{
							"salary:max": 120000,
						},
					},
				},
			},
			{
				Value: "Wellington",
				Buckets: []*ResultBucket{
					{
						Value:   "2015-12-01T00:00:00+13:00",
						Buckets: []*ResultBucket{},
						Metrics: map[string]interface{}{
							"salary:max": nil,
						},
					},
					{
						Value:   "2016-01-01T00:00:00+13:00",
						Buckets: []*ResultBucket{},
						Metrics: map[string]interface{}{
							"salary:max": 120000,
						},
					},
					{
						Value:   "2016-02-01T00:00:00+13:00",
						Buckets: []*ResultBucket{},
						Metrics: map[string]interface{}{
							"salary:max": 220000,
						},
					},
				},
			},
		},
	}
	results, err := Tabulate(input, 1)
	if err != nil {
		t.Fatalf("Unexpected error converting results: %s", err.Error())
	}

	rm, _ := json.MarshalIndent(*results, "", "  ")
	em, _ := json.MarshalIndent(expected, "", "  ")
	Expect(rm).To(MatchJSON(em))
}

func TestFlattenResultSetWithHoles(t *testing.T) {
	RegisterTestingT(t)
	expected := ResultTable{
		Rows: [][]map[string]interface{}{
			{
				{
					"salary:max": nil,
				},
				{
					"salary:max": 150000,
				},
				nil,
			},
			{
				{
					"salary:max": nil,
				},
				nil,
				{
					"salary:max": 120000,
				},
			},
		},
		RowTitles: [][]string{
			{"Auckland"},
			{"Wellington"},
		},
		ColumnTitles: [][]string{
			{"2015-12-01T00:00:00+13:00"},
			{"2016-01-01T00:00:00+13:00"},
			{"2016-02-01T00:00:00+13:00"},
		},
	}
	input := &Resultset{
		Buckets: []*ResultBucket{
			{
				Value: "Auckland",
				Buckets: []*ResultBucket{
					{
						Value:   "2015-12-01T00:00:00+13:00",
						Buckets: []*ResultBucket{},
						Metrics: map[string]interface{}{
							"salary:max": nil,
						},
					},
					{
						Value:   "2016-01-01T00:00:00+13:00",
						Buckets: []*ResultBucket{},
						Metrics: map[string]interface{}{
							"salary:max": 150000,
						},
					},
				},
			},
			{
				Value: "Wellington",
				Buckets: []*ResultBucket{
					{
						Value:   "2015-12-01T00:00:00+13:00",
						Buckets: []*ResultBucket{},
						Metrics: map[string]interface{}{
							"salary:max": nil,
						},
					},
					{
						Value:   "2016-02-01T00:00:00+13:00",
						Buckets: []*ResultBucket{},
						Metrics: map[string]interface{}{
							"salary:max": 120000,
						},
					},
				},
			},
		},
	}
	results, err := Tabulate(input, 1)
	if err != nil {
		t.Fatalf("Unexpected error converting results: %s", err.Error())
	}

	rm, _ := json.MarshalIndent(*results, "", "  ")
	em, _ := json.MarshalIndent(expected, "", "  ")
	Expect(rm).To(MatchJSON(em))
}

func TestFlattenResultDeep(t *testing.T) {
	RegisterTestingT(t)
	expected := ResultTable{
		Rows: [][]map[string]interface{}{
			{
				{
					"salary:max": 1111,
				},
				{
					"salary:max": 1121,
				},
				nil,
			},
			{
				{
					"salary:max": 1211,
				},
				nil,
				{
					"salary:max": 1212,
				},
			},
			{
				nil,
				{
					"salary:max": 2321,
				},
				nil,
			},
		},
		RowTitles: [][]string{
			{"A1", "B1"},
			{"A1", "B2"},
			{"A2", "B3"},
		},
		ColumnTitles: [][]string{
			{"C1", "D1"},
			{"C2", "D1"},
			{"C1", "D2"},
		},
	}
	input := &Resultset{
		Buckets: []*ResultBucket{
			{
				Value: "A1",
				Buckets: []*ResultBucket{
					{
						Value: "B1",
						Buckets: []*ResultBucket{
							{
								Value: "C1",
								Buckets: []*ResultBucket{
									{
										Value:   "D1",
										Buckets: []*ResultBucket{},
										Metrics: map[string]interface{}{
											"salary:max": 1111,
										},
									},
								},
							},
							{
								Value: "C2",
								Buckets: []*ResultBucket{
									{
										Value:   "D1",
										Buckets: []*ResultBucket{},
										Metrics: map[string]interface{}{
											"salary:max": 1121,
										},
									},
								},
							},
						},
					},
					{
						Value: "B2",
						Buckets: []*ResultBucket{
							{
								Value: "C1",
								Buckets: []*ResultBucket{
									{
										Value:   "D1",
										Buckets: []*ResultBucket{},
										Metrics: map[string]interface{}{
											"salary:max": 1211,
										},
									},
									{
										Value:   "D2",
										Buckets: []*ResultBucket{},
										Metrics: map[string]interface{}{
											"salary:max": 1212,
										},
									},
								},
							},
						},
					},
				},
			},
			{
				Value: "A2",
				Buckets: []*ResultBucket{
					{
						Value: "B3",
						Buckets: []*ResultBucket{
							{
								Value: "C2",
								Buckets: []*ResultBucket{
									{
										Value:   "D1",
										Buckets: []*ResultBucket{},
										Metrics: map[string]interface{}{
											"salary:max": 2321,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	results, err := Tabulate(input, 2)
	if err != nil {
		t.Fatalf("Unexpected error converting results: %s", err.Error())
	}

	rm, _ := json.MarshalIndent(*results, "", "  ")
	em, _ := json.MarshalIndent(expected, "", "  ")
	Expect(rm).To(MatchJSON(em))
}

func TestTabulateDepthTooHigh(t *testing.T) {
	input := &Resultset{
		Buckets: []*ResultBucket{
			{
				Value: "A1",
				Buckets: []*ResultBucket{
					{
						Value: "B1",
						Buckets: []*ResultBucket{
							{
								Value: "C1",
								Buckets: []*ResultBucket{
									{
										Value:   "D1",
										Buckets: []*ResultBucket{},
										Metrics: map[string]interface{}{
											"salary:max": 1111,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	_, err := Tabulate(input, 0)
	if err != ErrTargetDepthTooLow {
		t.Fatalf("Expected error converting results, got none")
	}
	_, err = Tabulate(input, 4)
	if err != ErrTargetDepthNotReached {
		t.Fatalf("Expected error converting results, got none")
	}
}
