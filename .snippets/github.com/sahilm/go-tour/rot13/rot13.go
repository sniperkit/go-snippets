package rot13

import (
	"io"
	"unicode"
)

// Reader is a Rot13 reader
type Reader struct {
	R io.Reader
}

// Read reads p converting it into a rot13 encoded byte array.
// Returns the number of bytes read and nil error.
func (reader Reader) Read(p []byte) (n int, err error) {
	n, err = reader.R.Read(p)
	if err != nil {
		return
	}
	for i := 0; i < n; i++ {
		p[i] = toRot13(p[i])
	}
	return
}

func toRot13(b byte) byte {
	transform := func(b, base byte) byte {
		return base + (b-base+13)%26
	}

	switch {
	case unicode.IsLower(rune(b)):
		return transform(b, 'a')
	case unicode.IsUpper(rune(b)):
		return transform(b, 'A')
	default:
		return b
	}
}
