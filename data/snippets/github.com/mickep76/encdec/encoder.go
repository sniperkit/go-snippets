package encdec

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

// Encoder interface.
type Encoder interface {
	SetIndent(indent string)
	Encode(value interface{}) error
}

// EncoderOption variadic function.
type EncoderOption func(Encoder)

// NewEncoder variadic constructor.
func NewEncoder(encoding string, writer io.Writer, options ...EncoderOption) Encoder {
	c, ok := encodings[encoding]
	if !ok {
		return nil
	}

	enc := c.NewEncoder(writer)
	for _, option := range options {
		option(enc)
	}

	return enc
}

// WithIndent output setter for JSON.
func WithIndent(indent string) EncoderOption {
	return func(e Encoder) {
		e.SetIndent(indent)
	}
}

// ToBytes method.
func ToBytes(encoding string, value interface{}, options ...EncoderOption) ([]byte, error) {
	var buf bytes.Buffer

	if err := NewEncoder(encoding, &buf, options...).Encode(value); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ToFile method.
func ToFile(encoding string, file string, value interface{}, options ...EncoderOption) error {
	fp, err := os.Create(file)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(fp)

	if err := NewEncoder(encoding, w, options...).Encode(value); err != nil {
		return err
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return fp.Close()
}
