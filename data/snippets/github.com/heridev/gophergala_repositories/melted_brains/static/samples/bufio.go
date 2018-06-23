package bufio2

import "io"

const (
	defaultBufSize  = 4096
	defaultBufCount = 4
)

type AsyncWriter struct {
	err       error
	buf       []byte
	n         int
	size      int
	countSent int
	countRcvd int

	wr io.Writer

	dataChan   chan []byte
	resultChan chan error
}

func NewAsyncWriterSize(wr io.Writer, size int, count int) *AsyncWriter {
	b, ok := wr.(*AsyncWriter)
	if ok && len(b.buf) >= size {
		return b
	}
	// code omitted ...
}