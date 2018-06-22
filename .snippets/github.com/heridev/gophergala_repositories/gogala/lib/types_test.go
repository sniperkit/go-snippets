package lib

import (
	"fmt"
	"testing"
)

func TestMessageString(t *testing.T) {

	m := Message{
		Kind: "test",
		Body: "This is a test",
	}
	a := fmt.Sprintf("%s", m)
	b := "test\nThis is a test\n"
	if a != b {
		t.Errorf("got %v want %v", a, b)
	}
}
