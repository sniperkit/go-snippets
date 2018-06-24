package aggro

import (
	"testing"
	"time"
)

func TestDatetimePeriodQuarter(t *testing.T) {
	for _, example := range []struct {
		t        time.Time
		expected string
	}{
		{time.Date(2016, 1, 1, 12, 12, 20, 0, time.UTC), "2016-01-01T00:00:00Z"},
		{time.Date(2016, 2, 2, 12, 12, 20, 0, time.UTC), "2016-01-01T00:00:00Z"},
		{time.Date(2016, 3, 20, 12, 12, 20, 0, time.UTC), "2016-01-01T00:00:00Z"},
		{time.Date(2016, 4, 22, 12, 12, 20, 0, time.UTC), "2016-04-01T00:00:00Z"},
		{time.Date(2016, 5, 22, 12, 12, 20, 0, time.UTC), "2016-04-01T00:00:00Z"},
		{time.Date(2016, 6, 22, 12, 12, 20, 0, time.UTC), "2016-04-01T00:00:00Z"},
		{time.Date(2016, 7, 22, 12, 12, 20, 0, time.UTC), "2016-07-01T00:00:00Z"},
		{time.Date(2016, 8, 22, 12, 12, 20, 0, time.UTC), "2016-07-01T00:00:00Z"},
		{time.Date(2016, 9, 22, 12, 12, 20, 0, time.UTC), "2016-07-01T00:00:00Z"},
		{time.Date(2016, 10, 22, 12, 12, 20, 0, time.UTC), "2016-10-01T00:00:00Z"},
		{time.Date(2016, 11, 22, 12, 12, 20, 0, time.UTC), "2016-10-01T00:00:00Z"},
		{time.Date(2016, 12, 22, 12, 12, 20, 0, time.UTC), "2016-10-01T00:00:00Z"},
	} {
		result, err := (&DatetimeCell{value: &example.t}).ValueForPeriod(Quarter, time.UTC)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		if result != example.expected {
			t.Fatalf("Unexpected result:\n\n\t%s did not equal expected %s", result, example.expected)
		}
	}
}

func TestDatetimePeriodWeek(t *testing.T) {
	for _, example := range []struct {
		t        time.Time
		expected string
	}{
		{time.Date(2016, 1, 1, 12, 12, 20, 0, time.UTC), "2015-12-27T00:00:00Z"},
		{time.Date(2016, 1, 2, 12, 12, 20, 0, time.UTC), "2015-12-27T00:00:00Z"},
		{time.Date(2016, 1, 3, 12, 12, 20, 0, time.UTC), "2016-01-03T00:00:00Z"},
		{time.Date(2016, 1, 4, 12, 12, 20, 0, time.UTC), "2016-01-03T00:00:00Z"},
		{time.Date(2016, 1, 5, 12, 12, 20, 0, time.UTC), "2016-01-03T00:00:00Z"},
		{time.Date(2016, 1, 6, 12, 12, 20, 0, time.UTC), "2016-01-03T00:00:00Z"},
		{time.Date(2016, 1, 7, 12, 12, 20, 0, time.UTC), "2016-01-03T00:00:00Z"},
		{time.Date(2016, 1, 8, 12, 12, 20, 0, time.UTC), "2016-01-03T00:00:00Z"},
		{time.Date(2016, 1, 9, 12, 12, 20, 0, time.UTC), "2016-01-03T00:00:00Z"},
		{time.Date(2016, 1, 10, 12, 12, 20, 0, time.UTC), "2016-01-10T00:00:00Z"},
	} {
		result, err := (&DatetimeCell{value: &example.t}).ValueForPeriod(Week, time.UTC)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		if result != example.expected {
			t.Fatalf("Unexpected result:\n\n\t%s did not equal expected %s", result, example.expected)
		}
	}
}
