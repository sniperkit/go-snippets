package basic

import (
	"bytes"
	"encoding/gob"
	"log"
	"strconv"
	"testing"
)

// GetBytes - convert arbitrary interface to byte array
// ref, https://stackoverflow.com/questions/23003793/convert-arbitrary-golang-interface-to-byte-array
func getBytes(val interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(val)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

//
func TestByteBinaryPrint(t *testing.T) {
	var b = byte(0x01)
	log.Println(b)
	log.Printf("%b", b)
	log.Printf("%v", b)
	log.Println(strconv.FormatInt(int64(b), 2))

	var bb = []byte{0x01, 0x02}
	log.Printf("%v", bb)
}

func TestByteLiteral(t *testing.T) {
	b := byte(0x21)
	log.Println(b)

	bb := []byte{0x00, 0x01}
	log.Println(bb)
}
