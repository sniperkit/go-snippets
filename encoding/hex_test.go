package encoding

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"log"
	"testing"
)

var src []byte
var dstStr string

func init() {
	tmp := md5.Sum([]byte("hello, hex encoding"))
	src = tmp[:]
	dstStr = easyWayEncodeHex(src)

	log.Printf("src(%v) dst string(%s)\n", src, dstStr)
}

func easyWayEncodeHex(src []byte) string {
	var dst string
	for _, v := range src {
		dst += fmt.Sprintf("%02x", v)
	}
	return dst
}

func TestEncode(t *testing.T) {
	dst := make([]byte, len(src)*2)
	count := Encode(dst, src)
	if count != len(src)*2 {
		t.Fatalf("Encode hex src(%v) to dst(%v) count not equal", src, dst, count)
	}
	if string(dst) != dstStr {
		t.Fatalf("Encode hex src(%v) to hex dst(%s) wrong with (%s)", src, string(dst), dstStr)
	}
}

func TestEncodeToString(t *testing.T) {
	dst := EncodeToString(src)
	if dst != dstStr {
		t.Fatalf("Encode to string src(%v) dst string(%s) not same as (%s)", src, dst, dstStr)
	}
}

func TestDecode(t *testing.T) {
	hex := make([]byte, len(src)*2)
	Encode(hex, src)

	dst := make([]byte, len(hex)/2)
	count, err := Decode(dst, hex)
	if err != nil {
		t.Fatalf("Decode src(%s) to original bytes(%v) error: %s", string(hex), src, err)
	}
	if count != len(hex)/2 {
		t.Fatalf("Decode src(%s) decode to wrong length: %d", string(hex), count)
	}
	if bytes.Compare(dst, src) != 0 {
		t.Fatalf("Decode src(%s) got wrong result(%v) right(%v)", string(hex), dst, src)
	}
}

func TestDecodeWrongSrcLen(t *testing.T) {
	src := []byte{89, 34, 23, 1, 0}
	dst := make([]byte, len(src)*2)
	_, err := Decode(dst, src)
	if err == nil {
		t.Fatalf("Decode src(%v) is wrong hex format with length", src)
	}
}

func TestDecodeWrongSrc(t *testing.T) {
	src := []byte{'i', 12, 34, 23, 1, 0}
	dst := make([]byte, len(src)*2)
	_, err := Decode(dst, src)
	if err == nil {
		t.Fatalf("Decode src(%v) is wrong hex format", src)
	}
}

func TestDecodeString(t *testing.T) {
	dst, err := DecodeString(dstStr)
	if err != nil {
		t.Fatalf("Decode string: dst(%v) error: %s", dstStr, err)
	}
	if bytes.Compare(dst, src) != 0 {
		t.Fatalf("Decode string: original(%s) dst(%v) src(%v)", dstStr, dst, src)
	}
}
