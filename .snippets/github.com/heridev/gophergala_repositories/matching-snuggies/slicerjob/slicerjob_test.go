package slicerjob

import "testing"

func TestSlicerJob(t *testing.T) {
	job := New()
	if job.ID == "" {
		t.Fatalf("new job missing ID")
	}
}
