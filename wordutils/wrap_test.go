package wordutils

import (
	"testing"
)

func TestWrapSimple(t *testing.T) {
	text := "12345 678 12345 678 12345 678"
	newText := Wrap(text, 10)
	expected := "12345 678\n12345 678\n12345 678"
	if newText != expected {
		t.Fatalf("Expected %s, got %s", expected, newText)
	}

	newText = Wrap(text, 19)
	expected = "12345 678 12345\n678 12345 678"
	if newText != expected {
		t.Fatalf("Expected %s, got %s", expected, newText)
	}

	newText = Wrap(text, 20)
	expected = "12345 678 12345 678\n12345 678"
	if newText != expected {
		t.Fatalf("Expected %s, got %s", expected, newText)
	}
}
