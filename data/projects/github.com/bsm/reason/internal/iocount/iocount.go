package iocount

import (
	"io"
)

// Reader wraps an io.Reader
type Reader struct {
	N int64
	R io.Reader
}

// Read implements io.Reader
func (r *Reader) Read(buf []byte) (int, error) {
	n, err := r.R.Read(buf)
	r.N += int64(n)
	return n, err
}

// Writer wraps an io.Writer
type Writer struct {
	N int64
	W io.Writer
}

// Writer implements io.Writer
func (w *Writer) Write(buf []byte) (int, error) {
	n, err := w.W.Write(buf)
	w.N += int64(n)
	return n, err
}
