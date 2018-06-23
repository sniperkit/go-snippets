package textsim

import "io"

type circularBuffer struct {
	buffer       [][]byte
	startPointer int
	numElements  int
}

func newCircularBuffer(size int) *circularBuffer {
	return &circularBuffer{
		buffer:       make([][]byte, size),
		startPointer: 0,
		numElements:  0,
	}
}

// addElem returns a displaced element iff the circular buffer is
// currently full.
func (c *circularBuffer) addElem(elem []byte) (displaced []byte) {
	if c.isFull() {
		displaced = c.buffer[c.startPointer]
	} else {
		c.numElements += 1
	}
	c.buffer[c.startPointer] = elem
	c.startPointer = (c.startPointer + 1) % len(c.buffer)
	return
}

func (c *circularBuffer) isFull() bool {
	return c.numElements == len(c.buffer)
}

func (c *circularBuffer) write(writer io.Writer) {
	if !c.isFull() {
		panic("Cannot write a non-full buffer")
	}
	for i := 0; i < len(c.buffer); i++ {
		writer.Write(c.buffer[(c.startPointer+i)%len(c.buffer)])
	}
}
