package iochan

import (
	"bufio"
	"bytes"
	"io"
	"time"
)

// DelimReader reads from io.Reader till delimiter, send readed bytes to channel
// works while io.Reader
func DelimReader(r io.Reader, delimiter byte) <-chan string {
	ch := make(chan string)

	go func() {
		buf := bufio.NewReader(r)

		for {
			line, err := buf.ReadString(delimiter)
			if err != nil {
				break
			}

			ch <- line
		}

		close(ch)
	}()

	return ch
}

// SizeReader reads from io.Reader upto size, send readed bytes to channel
func SizeReader(r io.Reader, size int) <-chan string {
	ch := make(chan string)

	go func() {
		buf := bufio.NewReader(r)

		for {
			line := make([]byte, size)
			if _, err := buf.Read(line); err != nil {
				break
			}

			ch <- string(line)
		}

		close(ch)
	}()

	return ch
}

// TimedWrite reads from chan and writes to w when timer expires
func TimedWrite(ch <-chan string, timer time.Duration, w io.Writer) {
	go func() {
		var buf bytes.Buffer

		for {
			select {
			case data := <-ch:
				buf.WriteString(data)

			case <-time.After(timer):
				w.Write(buf.Bytes())
				buf.Reset()
			}
		}
	}()
}

// SizedWrite reads from chan and writes to w when data exceeds size bytes
func SizedWrite(ch <-chan string, size int, w io.Writer) {
	go func() {
		var buf bytes.Buffer

		for {
			data := <-ch
			buf.WriteString(data)

			if buf.Len() >= size {
				w.Write(buf.Bytes())
				buf.Reset()
			}
		}
	}()
}
