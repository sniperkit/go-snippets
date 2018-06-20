package encdec

import (
	"io"
)

var encodings = make(map[string]Encoding)

// Encoding interface.
type Encoding interface {
	NewEncoder(writer io.Writer) Encoder
	NewDecoder(reader io.Reader) Decoder
}

// Register encoding.
func Register(name string, encoding Encoding) {
	encodings[name] = encoding
}

// Registered encoding.
func Registered(name string) bool {
	_, ok := encodings[name]
	if !ok {
		return false
	}
	return true
}

// Encodings registered.
func Encodings() []string {
	encs := []string{}
	for k := range encodings {
		encs = append(encs, k)
	}
	return encs
}
