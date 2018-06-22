package lib

import (
	"testing"
	"time"
)

func TestMakeArgs(t *testing.T) {
	a := MakeArgs(1, "b", []int{1, 2})
	if len(a) != 3 {
		t.Errorf("got %v want 3", len(a))
	}
}

func TestAppendString(t *testing.T) {
	a := AppendString("Hello", " Go", "lang")
	b := "Hello Golang"

	if len(a) != len(b) {
		t.Errorf("got %v want %v", len(a), len(b))
	}

	if a != b {
		t.Errorf("got %v want %v", a, b)
	}
}

func TestPrintTimeStamp(t *testing.T) {
	a := PrintTimeStamp()
	b := time.Now().Format(time.Kitchen)

	if a != b {
		t.Errorf("got %v want %v", a, b)
	}

}
