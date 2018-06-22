package human

import (
	"bytes"
	"io"
)

// FlushableBuffer wraps around an io.Writer and provides
// the ability to lazily write to the buffer when the Flush method
// is called.
//
// This allows for delay of actual writing to the io.Buffer, providing
// an interface similar to the io.Writer.
type FlushableBuffer struct {
	*bytes.Buffer
	stream io.Writer
}

// Flush writes the data written to the FlushableBuffer to the underlying
// stream.
func (b *FlushableBuffer) Flush() (n int, err error) {
	return b.stream.Write(b.Bytes())
}

// NewFlushableBuffer creates a new FlushableBuffer for a given io.Writer
func NewFlushableBuffer(stream io.Writer) *FlushableBuffer {
	return &FlushableBuffer{
		Buffer: bytes.NewBufferString(""),
		stream: stream,
	}
}
