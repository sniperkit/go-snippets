package guid

import (
	"fmt"
	"testing"
	"time"
)

const (
	// for manual debugging
	showIds = false
)

func Test(t *testing.T) {
	id1 := New()
	if len(id1) != size {
		t.Fatalf("len(id1) != %d (=%d)", size, len(id1))
	}
	id2 := New()
	if len(id2) != size {
		t.Fatalf("len(id2) != %d (=%d)", size, len(id2))
	}
	if id1 == id2 {
		t.Fatalf("generated same ids (id1: '%s', id2: '%s')", id1, id2)
	}
	if showIds {
		fmt.Printf("%s\n", id1)
		fmt.Printf("%s\n", id2)
		time.Sleep(2 * time.Millisecond)
		fmt.Printf("%s\n", New())
		time.Sleep(2 * time.Millisecond)
		fmt.Printf("%s\n", New())
		time.Sleep(2 * time.Millisecond)
		fmt.Printf("%s\n", New())
	}
}
