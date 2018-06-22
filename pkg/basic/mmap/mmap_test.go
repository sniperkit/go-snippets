package mmap

import (
	"log"
	"testing"
)

func TestGetFromNilMap(t *testing.T) {
	var m map[string]string

	val := m["hello"]

	log.Printf("m is empty: %v, val is \"\": %v", m == nil, val == "")
}
