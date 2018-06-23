package encdec

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

// Decoder interface.
type Decoder interface {
	SetMapString()
	Decode(value interface{}) error
}

// DecoderOption function.
type DecoderOption func(Decoder)

// WithMapString convert map[interface{}]interface{} to map[string]interface{} for YAML.
func WithMapString() DecoderOption {
	return func(d Decoder) {
		d.SetMapString()
	}
}

// NewDecoder variadic constructor.
func NewDecoder(encoding string, reader io.Reader, options ...DecoderOption) Decoder {
	c, ok := encodings[encoding]
	if !ok {
		return nil
	}

	dec := c.NewDecoder(reader)
	for _, option := range options {
		option(dec)
	}

	return dec
}

// FromBytes method.
func FromBytes(encoding string, encoded []byte, value interface{}, options ...DecoderOption) error {
	r := bufio.NewReader(bytes.NewReader(encoded))
	return NewDecoder(encoding, r, options...).Decode(value)
}

// FromFile method.
func FromFile(encoding string, file string, value interface{}, options ...DecoderOption) error {
	fp, err := os.Open(file)
	if err != nil {
		return err
	}

	r := bufio.NewReader(fp)

	if err := NewDecoder(encoding, r, options...).Decode(value); err != nil {
		return err
	}

	return fp.Close()
}
