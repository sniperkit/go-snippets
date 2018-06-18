package humanize

import "testing"

// TestOneline verifies various values get converted correctly
func TestOneline(t *testing.T) {
	var test1 = `Foo bar baz
	quux`
	var test2 = `foo?  Foo.
	bar?
	baz!
	quux.`

	var expected, actual string

	expected = "Foo bar baz quux"
	actual = Oneline(test1)
	if expected != actual {
		t.Errorf("Expected %q, but got %q", expected, actual)
	}

	// Note the horrible spacing after sentence enders which were followed by a
	// newline.  SO GRODY.
	expected = "foo?  Foo. bar? baz! quux."
	actual = Oneline(test2)
	if expected != actual {
		t.Errorf("Expected %q, but got %q", expected, actual)
	}
}
