package unittest

import (
	"github.com/Sirupsen/logrus"
	"fmt"
	"io"
	"os"
)

const stdHijackerBufferSize = 4096

// stdout/stderr hijacker.
type stdHijacker struct {
	r *os.File
	w *os.File
	bufs [][]byte
	done chan struct{}
	buf []byte

	// Guaranteed access to stdout.
	Stdout *os.File

	// Guaranteed access to stderr.
	Stderr *os.File
}

// New stdout/stderr hijacker.
func newStdHijacker() (h *stdHijacker, err error) {
	h = &stdHijacker{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		bufs: make([][]byte, 0),
		done: make(chan struct{}),
	}
	h.r, h.w, err = os.Pipe()
	return
}

// Hijack.
func (h *stdHijacker) Hijack() {
	if h.w == nil {
		panic("stdHijacker can only be used once")
	}

	// Assign the stdout and stderr to the pipe.
	os.Stdout = h.w
	os.Stderr = h.w

	// Acquire the log output from logrus.
	logrus.SetOutput(h.w)

	// Start a draining function.
	go func() {
		var n int
		var err error

		// Read out as much data as possible.
		for {
			// Read out the data.
			b := make([]byte, stdHijackerBufferSize)

			if n, err = h.r.Read(b); err != nil {
				if err == io.EOF {
					break
				}

				panic(fmt.Sprintf("failed to read from stdout/stderr hijack pipe's read end: %v", err))
			}

			// Append the buffer.
			h.bufs = append(h.bufs, b[:n])
		}

		// Close the read end of the pipe.
		h.r.Close()

		// Signal completion of draining.
		h.r = nil
		h.done <- struct{}{}
	}()
}

// Release.
func (h *stdHijacker) Release() {
	// Reassign the log output for logrus.
	logrus.SetOutput(h.Stdout)

	// Assign stdout and stderr to the original values.
	os.Stdout = h.Stdout
	os.Stderr = h.Stderr

	// Close the write end of the pipe.
	h.w.Close()
	h.w = nil

	// Wait for the pipe to be drained.
	<- h.done

	// Concatenate the buffers.
	bufsSize := 0

	for _, b := range h.bufs {
		bufsSize += len(b)
	}

	h.buf = make([]byte, bufsSize)
	offset := 0
	for _, b := range h.bufs {
		n := len(b)

		copy(h.buf[offset:offset+n], b)
		offset += n
	}

	h.bufs = nil
}

// Bytes.
//
// Combined output from stdout/stderr. Safe to access after return from
// Release.
func (h *stdHijacker) Bytes() []byte {
	return h.buf
}
