package io

import (
	"bytes"
	"io"
	"log"
	"testing"
)

func TestReaderToByte(t *testing.T) {
	bs := []byte("reader")
	p := make([]byte, len(bs))

	// normal case
	r := bytes.NewReader(bs)
	_, err := io.ReadFull(r, p)
	log.Printf("content: %s, err: %v", string(p), err)

	// failure case
	p = make([]byte, 1024)
	r = bytes.NewReader(bs)
	n, err := r.Read(p)
	log.Printf("n: %v, cap: %v, err: %v", n, len(bs), err)
}
