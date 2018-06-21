package protocol

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"testing"
)

func testhostname() string {
	hash := crc32.NewIEEE()
	hash.Write([]byte("12345"))

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func TestNewMessage(t *testing.T) {
	msg := NewMessage(StandardMessage, NoCompression, testhostname())

	_, err := msg.Write([]byte("this is a message"))
	if err != nil {
		t.Error(err)
	}

	// t.Logf("wrote %d bytes", n)
	// msg.Print()
}

func TestNewMessage_Snappy(t *testing.T) {
	msg := NewMessage(StandardMessage, SnappyCompression, testhostname())

	_, err := msg.Write([]byte("this is a message"))
	if err != nil {
		t.Error(err)
	}

	// t.Logf("wrote %d bytes", n)
}

func TestNewMessage_Deflate(t *testing.T) {
	msg := NewMessage(StandardMessage, DeflateCompression, testhostname())

	_, err := msg.Write([]byte("this is a message"))
	if err != nil {
		t.Error(err)
	}

	// t.Logf("wrote %d bytes", n)
}

func TestDecode(t *testing.T) {
	msg := NewMessage(StandardMessage, NoCompression, testhostname())

	content := "Supercalifragilisticexpialidocious"

	msg.Write([]byte(content))

	msgb := bytes.NewBuffer(nil)
	msgb.ReadFrom(msg)

	t.Logf("msg length (compressed): %d", msg.Length())
	t.Logf("msg length (uncompressed): %d", msg.DecodedLength())

	decmsg, err := Decode(msgb.Bytes())
	if err != nil {
		t.Error(err)
	}

	t.Logf("decoded msg length (compressed): %d", decmsg.Length())
	t.Logf("decoded msg length (uncompressed): %d", decmsg.DecodedLength())

	if decmsg.Header.MsgType != msg.Header.MsgType {
		t.Errorf("decoding msg type failed")
	}

	if decmsg.Header.CompressionType != msg.Header.CompressionType {
		t.Errorf("decoding compression type failed")
	}

	if decmsg.Header.SourceId != testhostname() {
		t.Errorf("decoding source id header failed, found %s, expected %s", decmsg.Header.SourceId, testhostname())
	}

	if decmsg.Header.Timestamp != msg.Header.Timestamp {
		t.Errorf("decoding source Timestamp header failed, found %d, expected %d", decmsg.Header.Timestamp, msg.Header.Timestamp)
	}

	if decmsg.Length() != msg.Length() {
		t.Errorf("decoded message length did not match original message length: %d, %d", decmsg.Length(), msg.Length())
	}

	if !bytes.Equal(decmsg.buf, msg.buf) {
		t.Errorf("decoding message content failed, found:\n%v\nexpected:\n%v\n", decmsg.buf, msg.buf)
	}

	body, err := decmsg.Body()
	if err != nil {
		t.Error(err)
	}

	bb := string(body)

	if content != bb {
		t.Errorf("incorrect body decoded, found: %s, expected: %s", string(body), content)
	}
}

func TestDecode_Deflate(t *testing.T) {
	msg := NewMessage(StandardMessage, DeflateCompression, testhostname())

	content := "Supercalifragilisticexpialidocious"

	msg.Write([]byte(content))

	msgb := bytes.NewBuffer(nil)
	msgb.ReadFrom(msg)

	t.Logf("msg length (compressed): %d", msg.Length())
	t.Logf("msg length (uncompressed): %d", msg.DecodedLength())

	decmsg, err := Decode(msgb.Bytes())
	if err != nil {
		t.Error(err)
	}

	t.Logf("decoded msg length (compressed): %d", decmsg.Length())
	t.Logf("decoded msg length (uncompressed): %d", decmsg.DecodedLength())

	if decmsg.Header.MsgType != msg.Header.MsgType {
		t.Errorf("decoding msg type failed")
	}

	if decmsg.Header.CompressionType != msg.Header.CompressionType {
		t.Errorf("decoding compression type failed")
	}

	if decmsg.Header.SourceId != testhostname() {
		t.Errorf("decoding source id header failed, found %s, expected %s", decmsg.Header.SourceId, testhostname())
	}

	if decmsg.Header.Timestamp != msg.Header.Timestamp {
		t.Errorf("decoding source Timestamp header failed, found %d, expected %d", decmsg.Header.Timestamp, msg.Header.Timestamp)
	}

	if !bytes.Equal(decmsg.buf, msg.buf) {
		t.Errorf("decoding message content failed, found:\n%v\nexpected:\n%v\n", decmsg.buf, msg.buf)
	}

	body, err := decmsg.Body()
	if err != nil {
		t.Error(err)
	}

	if content != string(body) {
		t.Errorf("incorrect body decoded")
	}
}

func TestDecode_Zlib(t *testing.T) {
	msg := NewMessage(StandardMessage, ZlibCompression, testhostname())

	content := "Supercalifragilisticexpialidocious"

	msg.Write([]byte(content))

	msgb := bytes.NewBuffer(nil)
	msgb.ReadFrom(msg)

	t.Logf("msg length (compressed): %d", msg.Length())
	t.Logf("msg length (uncompressed): %d", msg.DecodedLength())

	decmsg, err := Decode(msgb.Bytes())
	if err != nil {
		t.Error(err)
	}

	t.Logf("decoded msg length (compressed): %d", decmsg.Length())
	t.Logf("decoded msg length (uncompressed): %d", decmsg.DecodedLength())

	if decmsg.Header.MsgType != msg.Header.MsgType {
		t.Errorf("decoding msg type failed")
	}

	if decmsg.Header.CompressionType != msg.Header.CompressionType {
		t.Errorf("decoding compression type failed")
	}

	if decmsg.Header.SourceId != testhostname() {
		t.Errorf("decoding source id header failed, found %s, expected %s", decmsg.Header.SourceId, testhostname())
	}

	if decmsg.Header.Timestamp != msg.Header.Timestamp {
		t.Errorf("decoding source Timestamp header failed, found %d, expected %d", decmsg.Header.Timestamp, msg.Header.Timestamp)
	}

	if decmsg.Length() != msg.Length() {
		t.Errorf("decoded message length did not match original message length: %d, %d", decmsg.Length(), msg.Length())
	}

	if !bytes.Equal(decmsg.buf, msg.buf) {
		t.Errorf("decoding message content failed, found:\n%v\nexpected:\n%v\n", decmsg.buf, msg.buf)
	}

	body, err := decmsg.Body()
	if err != nil {
		t.Error(err)
	}

	if content != string(body) {
		t.Errorf("incorrect body decoded")
	}
}

func TestDecode_Snappy(t *testing.T) {
	msg := NewMessage(StandardMessage, SnappyCompression, testhostname())

	content := "Supercalifragilisticexpialidocious"

	msg.Write([]byte(content))

	msgb := bytes.NewBuffer(nil)
	msgb.ReadFrom(msg)

	t.Logf("original msg length (compressed): %d", msg.Length())
	t.Logf("original msg length (uncompressed): %d", msg.DecodedLength())

	decmsg, err := Decode(msgb.Bytes())
	if err != nil {
		t.Error(err)
	}

	t.Logf("decoded msg length (compressed): %d", decmsg.Length())
	t.Logf("decoded msg length (uncompressed): %d", decmsg.DecodedLength())

	if decmsg.Header.MsgType != msg.Header.MsgType {
		t.Errorf("decoding msg type failed")
	}

	if decmsg.Header.CompressionType != msg.Header.CompressionType {
		t.Errorf("decoding compression type failed")
	}

	if decmsg.Header.SourceId != testhostname() {
		t.Errorf("decoding source id header failed, found %s, expected %s", decmsg.Header.SourceId, testhostname())
	}

	if decmsg.Header.Timestamp != msg.Header.Timestamp {
		t.Errorf("decoding source Timestamp header failed, found %d, expected %d", decmsg.Header.Timestamp, msg.Header.Timestamp)
	}

	if decmsg.Length() != msg.Length() {
		t.Errorf("decoded message length did not match original message length: %d, %d", decmsg.Length(), msg.Length())
	}

	if decmsg.Length() != msg.Length() {
		t.Errorf("decoded message length did not match original message length: %d, %d", decmsg.Length(), msg.Length())
	}

	if !bytes.Equal(decmsg.buf, msg.buf) {
		t.Errorf("decoding message content failed, found:\n%v\nexpected:\n%v\n", decmsg.buf, msg.buf)
	}

	body, err := decmsg.Body()
	if err != nil {
		t.Error(err)
	}

	if content != string(body) {
		t.Errorf("incorrect body decoded")
	}
}
