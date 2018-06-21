package protocol

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"regexp"
	"sync"
	"time"

	"github.com/zerklabs/auburn/log"
)

// const HeaderSize = 36

type MessageHeader struct {
	lock sync.Mutex

	Version         int
	MsgType         int
	CompressionType int
	SourceId        string
	Timestamp       int64
	BodyLen         int
	CompBodyLen     int

	off int // buf offset
}

// rebuildHeader will return the starting index of the body and the MessageHeader
func rebuildHeader(b []byte) (int, MessageHeader, error) {
	// we need to find the NULSEP byte sequence so we know where the header ends
	var bodyidx int
	var headeridx int

	re := regexp.MustCompile("NULSEP")
	sepidx := re.FindIndex(b)
	if sepidx == nil || len(sepidx) == 0 {
		return 0, MessageHeader{}, fmt.Errorf("failed to parse message header")
	}

	headeridx = sepidx[0]
	bodyidx = sepidx[1]

	if headeridx == 0 {
		return 0, MessageHeader{}, fmt.Errorf("failed to parse the message header")
	}

	buf := bytes.NewBuffer(b[:headeridx])
	dec := gob.NewDecoder(buf)
	var header MessageHeader
	if err := dec.Decode(&header); err != nil {
		if err != io.EOF {
			return 0, MessageHeader{}, err
		}
	}

	return bodyidx, header, nil
}

func buildMessageHeader(msgtype int, comptype int, id string) MessageHeader {
	now := time.Now().UnixNano()

	return MessageHeader{
		Version:         1,
		MsgType:         msgtype,
		SourceId:        id,
		Timestamp:       now,
		CompressionType: comptype,
		BodyLen:         0,
		CompBodyLen:     0,
		off:             0, // offset
	}
}

func (h *MessageHeader) Print() {
	log.Infof("%#v", h)
}

func (h *MessageHeader) Length() int {
	b, err := h.encode()
	if err != nil {
		panic(err)
	}

	return len(b)
}

func (h *MessageHeader) encode() ([]byte, error) {
	bytebuf := bytes.NewBuffer(nil)

	enc := gob.NewEncoder(bytebuf)
	if err := enc.Encode(h); err != nil {
		return make([]byte, 0), err
	}

	bytebuf.Write([]byte("NULSEP"))

	return bytebuf.Bytes(), nil
}

func (h *MessageHeader) Read(p []byte) (n int, err error) {
	b, err := h.encode()
	if err != nil {
		return 0, err
	}

	if h.off >= len(b) {
		if len(p) == 0 {
			return
		}
		return 0, io.EOF
	}

	n = copy(p, b[h.off:])
	h.off += n

	return n, nil
}
