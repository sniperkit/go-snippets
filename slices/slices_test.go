package slices_test

import (
	"testing"

	"github.com/sahilm/go-tour/slices"
)

func TestPic(t *testing.T) {
	got := len(slices.Pic(4, 4))
	want := 4
	if got != want {
		t.Errorf("got len: %v, want len: %v", got, want)
	}
}
