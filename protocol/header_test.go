package protocol

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"testing"
)

func TestBuildMessageHeader(t *testing.T) {
	hash := crc32.NewIEEE()
	hash.Write([]byte("12345"))

	header := buildMessageHeader(StandardMessage, NoCompression, fmt.Sprintf("%x", hash.Sum(nil)))
	buf := bytes.NewBuffer(nil)
	buf.ReadFrom(&header)

	// t.Logf("%v", buf.Bytes())

	msg, err := Decode(buf.Bytes())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if msg.Header.MsgType != header.MsgType {
		t.Errorf("msgtype headers did not match: %#v, %#v", header, msg.Header)
	}

	if msg.Header.CompressionType != header.CompressionType {
		t.Errorf("compression type headers did not match: %#v, %#v", header, msg.Header)
	}

	if msg.Header.Version != header.Version {
		t.Errorf("version headers did not match: %#v, %#v", header, msg.Header)
	}

	if msg.Header.SourceId != header.SourceId {
		t.Errorf("source id headers did not match: %#v, %#v", header, msg.Header)
	}

	if msg.Header.Timestamp != header.Timestamp {
		t.Errorf("source id headers did not match: %#v, %#v", header, msg.Header)
	}

	// t.Logf("%#v, %#v", header, msg.Header)

	// header.Print()
	// mtlen := len(header.MsgTypeHdr)

	// if mtlen != 1 {
	// 	t.Errorf("expected MsgTypeHdr to be 1 byte, was %d", mtlen)
	// }

	// srcidlen := len(header.SourceIdHdr)

	// if srcidlen != 8 {
	// 	t.Errorf("expected SourceIdHdr to be 8 bytes, was %d", srcidlen)
	// }

	// srctslen := len(header.SourceTimestampHdr)

	// if srctslen != 19 {
	// 	t.Errorf("expected SourceIdHdr to be 19 bytes, was %d", srctslen)
	// }

	// ctlen := len(header.CompressionTypeHdr)

	// if ctlen != 1 {
	// 	t.Errorf("expected CompressionTypeHdr to be 1 bytes, was %d", ctlen)
	// }
}
