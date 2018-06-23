package factory

import (
	"testing"
)

func TestSonyflake(t *testing.T) {
	id := Sonyflake.NextID()
	if id < 1 {
		t.Errorf("Unexpected Sonyflake ID, should be non-zero, got %d", id)
	}

	id2 := Sonyflake.NextID()
	if id > id2 {
		t.Errorf("IDs should be k-sortable, ascending, got %d > %d", id, id2)
	}
}
