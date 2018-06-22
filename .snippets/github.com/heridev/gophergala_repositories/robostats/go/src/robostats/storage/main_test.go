package storage

import (
	"testing"
)

func TestConnect(t *testing.T) {
	d := DB()
	if d == nil {
		t.Fatal("Database expected.")
	}
}

func TestGetCollection(t *testing.T) {
	c := C("test_collection")
	if c == nil {
		t.Fatal("Collection expected.")
	}
}
