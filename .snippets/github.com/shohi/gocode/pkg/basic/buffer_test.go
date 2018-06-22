package basic

import (
	"bytes"
	"log"
	"testing"
)

func TestBufferFromBytes(t *testing.T) {
	var data []byte
	buf := bytes.NewBuffer(data)

	log.Println(buf)
}
