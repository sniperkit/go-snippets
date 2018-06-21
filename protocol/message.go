package protocol

import (
	"bytes"
	"compress/flate"
	"compress/zlib"
	"fmt"
	"io"
	"sync"

	"github.com/mreiferson/go-snappystream"
	"github.com/zerklabs/auburn/log"
)

type Message struct {
	lock   sync.Mutex
	Header MessageHeader
	buf    []byte
	off    int // read offset
}

func (m *Message) Print() {
	log.Infof("%#v", m)
}

// Length will return the entire message length with the compressed body, including the header
func (m *Message) Length() int {
	return m.Header.Length() + len(m.buf)
}

// DecodedLength will return the entire message length with the decompressed body, including the header
func (m *Message) DecodedLength() int {
	b, err := m.decodebody()
	if err != nil {
		panic(err)
	}

	return m.Header.Length() + len(b)
}

// Version returns the message format version
func (m *Message) Version() int {
	return m.Header.Version
}

// MessageType returns the message type
func (m *Message) MessageType() int {
	return m.Header.MsgType
}

// CompressionType returns the type of compression used on the
// body of the message
func (m *Message) CompressionType() int {
	return m.Header.CompressionType
}

// Timestamp returns the unix timestamp when the message was created
// on the sender side. (nanosecond format)
func (m *Message) Timestamp() int64 {
	return m.Header.Timestamp
}

// Sender returns the string ID of the host who sent the message
func (m *Message) Sender() string {
	return m.Header.SourceId
}

// Body returns the decompressed body as a byte slice
func (m *Message) Body() ([]byte, error) {
	return m.decodebody()
}

// decodebody returns the decompressed body as a byte slice. It will
// check the compression type by the header value
func (m *Message) decodebody() ([]byte, error) {
	buf := bytes.NewBuffer(m.buf)
	rawbuf := bytes.NewBuffer(nil)

	switch m.Header.CompressionType {
	case NoCompression:
		rawbuf.ReadFrom(buf)

		break
	case DeflateCompression:
		r := flate.NewReader(buf)
		if _, nerr := rawbuf.ReadFrom(r); nerr != nil {
			return rawbuf.Bytes(), fmt.Errorf("error reading from flate stream: %v", nerr)
		}

		if nerr := r.Close(); nerr != nil {
			return rawbuf.Bytes(), fmt.Errorf("error closing flate stream: %v", nerr)
		}
	case ZlibCompression:
		gz, nerr := zlib.NewReader(buf)
		if nerr != nil {
			return rawbuf.Bytes(), nerr
		}

		if _, nerr := rawbuf.ReadFrom(gz); nerr != nil {
			return rawbuf.Bytes(), fmt.Errorf("error reading from zlib stream: %v", nerr)
		}

		if nerr := gz.Close(); nerr != nil {
			return rawbuf.Bytes(), fmt.Errorf("error closing zlib stream: %v", nerr)
		}
	case SnappyCompression:
		w := snappystream.NewReader(buf, false)
		if _, nerr := rawbuf.ReadFrom(w); nerr != nil {
			return rawbuf.Bytes(), fmt.Errorf("error reading from snappystream: %v", nerr)
		}
	}

	return rawbuf.Bytes(), nil
}

func (m *Message) Read(p []byte) (n int, err error) {
	rawbuf := bytes.NewBuffer(nil)
	if _, err := rawbuf.ReadFrom(&m.Header); err != nil {
		return 0, err
	}

	bodybytes, err := m.decodebody()
	if err != nil {
		return 0, err
	}

	rawbuf.Write(bodybytes)

	if m.off >= rawbuf.Len() {
		if len(p) == 0 {
			return
		}
		return 0, io.EOF
	}

	n = copy(p, rawbuf.Bytes()[m.off:])
	m.off += n

	return
}

// Write passes the encoding of the input to the encode() function
func (m *Message) Write(p []byte) (int64, error) {
	rawbuf := bytes.NewBuffer(p)
	return m.encode(rawbuf)
}

func (m *Message) encode(rawbuf *bytes.Buffer) (int64, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	var written int64
	buf := bytes.NewBuffer(nil)

	// set the body length to the uncompressed input length
	m.Header.BodyLen = rawbuf.Len()

	switch m.Header.CompressionType {
	case NoCompression:
		buf.ReadFrom(rawbuf)
		written += int64(buf.Len())
	case DeflateCompression:
		w, err := flate.NewWriter(buf, flate.BestCompression)
		if err != nil {
			return 0, err
		}

		nn, err := rawbuf.WriteTo(w)
		written += nn
		if err != nil {
			return written, fmt.Errorf("error writing to flate stream: %v", err)
		}

		if err := w.Flush(); err != nil {
			return written, fmt.Errorf("error flushing flate stream: %v", err)
		}

		if err := w.Close(); err != nil {
			return written, fmt.Errorf("error closing flate stream: %v", err)
		}
	case ZlibCompression:
		gz, err := zlib.NewWriterLevel(buf, flate.BestCompression)
		if err != nil {
			return 0, err
		}

		nn, err := rawbuf.WriteTo(gz)
		written += nn
		if err != nil {
			return written, fmt.Errorf("error writing to zlib stream: %v", err)
		}

		if err := gz.Flush(); err != nil {
			return written, fmt.Errorf("error flushing zlib stream: %v", err)
		}

		if err := gz.Close(); err != nil {
			return written, fmt.Errorf("error closing zlib stream: %v", err)
		}
	case SnappyCompression:
		w := snappystream.NewWriter(buf)

		nn, err := rawbuf.WriteTo(w)
		written += nn
		if err != nil {
			return written, fmt.Errorf("error writing to snappystream: %v", err)
		}
	}

	m.buf = buf.Bytes()
	m.Header.CompBodyLen = len(m.buf)

	return written, nil
}
