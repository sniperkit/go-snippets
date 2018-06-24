package protoio

import (
	"bufio"
	"encoding/binary"
	"io"
	"math"

	"github.com/gogo/protobuf/proto"
)

// Reader reads raw protobuf protocol from streams.
type Reader struct {
	*bufio.Reader

	tmp []byte
}

// ReadField reads the next field.
func (r *Reader) ReadField() (uint32, uint8, error) {
	u, err := r.ReadVarint()
	if err != nil {
		return 0, 0, err
	}
	return uint32(u) >> 3, uint8(u) & 0x07, nil
}

// ReadVarint reads a number.
func (r *Reader) ReadVarint() (uint64, error) {
	return binary.ReadUvarint(r)
}

// ReadString reads a string.
func (r *Reader) ReadString() (string, error) {
	b, err := r.readBytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ReadDouble reads a float64.
func (r *Reader) ReadDouble() (float64, error) {
	b := r.getBuffer(8)
	_, err := io.ReadFull(r, b)
	if err != nil {
		return 0, err
	}

	return math.Float64frombits(binary.LittleEndian.Uint64(b)), nil
}

// ReadMessage reads a message.
func (r *Reader) ReadMessage(m proto.Message) error {
	b, err := r.readBytes()
	if err != nil {
		return err
	}
	return proto.Unmarshal(b, m)
}

func (r *Reader) readBytes() ([]byte, error) {
	u, err := r.ReadVarint()
	if err != nil {
		return nil, err
	}

	b := r.getBuffer(int(u))
	_, err = io.ReadFull(r, b)
	return b, err
}

func (r *Reader) getBuffer(n int) []byte {
	if cap(r.tmp) < n {
		r.tmp = make([]byte, 0, n)
	}
	return r.tmp[:n]
}

// ---------------------------------------------------------------------

// Reader reads raw protobuf protocol from streams.
type Writer struct {
	*bufio.Writer

	tmp []byte
}

// WriteField writes field info.
func (w *Writer) WriteField(tag uint32, wire uint8) error {
	n := (uint64(tag) << 3) | uint64(wire&0x07)
	return w.WriteVarint(n)
}

// WriteVarint writes a number.
func (w *Writer) WriteVarint(u uint64) error {
	b := w.getBuffer(binary.MaxVarintLen64)
	n := binary.PutUvarint(b, u)
	_, err := w.Write(b[:n])
	return err
}

// WriteDouble writes a float64.
func (w *Writer) WriteDouble(v float64) error {
	b := w.getBuffer(8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(v))
	_, err := w.Write(b)
	return err
}

// WriteVarintField writes a number field.
func (w *Writer) WriteVarintField(tag uint32, u uint64) error {
	if err := w.WriteField(tag, proto.WireVarint); err != nil {
		return err
	}
	return w.WriteVarint(u)
}

// WriteStringField writes a string.
func (w *Writer) WriteStringField(tag uint32, s string) error {
	if err := w.WriteField(tag, proto.WireBytes); err != nil {
		return err
	}

	if err := w.WriteVarint(uint64(len(s))); err != nil {
		return err
	}

	_, err := w.Writer.WriteString(s)
	return err
}

// WriteMessageField writes a message.
func (w *Writer) WriteMessageField(tag uint32, m proto.Message) error {
	data, err := proto.Marshal(m)
	if err != nil {
		return err
	}

	if err := w.WriteField(tag, proto.WireBytes); err != nil {
		return err
	}

	if err := w.WriteVarint(uint64(len(data))); err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

func (w *Writer) getBuffer(n int) []byte {
	if cap(w.tmp) < n {
		w.tmp = make([]byte, n)
	}
	return w.tmp[:n]
}
