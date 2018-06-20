package factory

import (
	"testing"
)

func TestDatabase(t *testing.T) {
	db := DB{}
	if db.DB != nil {
		t.Errorf("DB instance expected nil")
	}
	if db.Profiler != nil {
		t.Errorf("DB profiler expected nil")
	}

	db.Profiler = &DatabaseProfilerStdout{}

	if db.Quiet().Profiler != nil {
		t.Errorf("DB quiet profiler expected nil")
	}
}
