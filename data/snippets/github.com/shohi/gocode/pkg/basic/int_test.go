package basic

import (
	"encoding/binary"
	"log"
	"math"
	"testing"
)

func TestUintConvert(t *testing.T) {
	var a int
	a = -1
	log.Println(uint32(a))
}

func TestUint32(t *testing.T) {
	a := uint32(10)
	var b uint32
	b = math.MaxUint32

	log.Println(a, b+uint32(8))
	log.Println(a > (b + uint32(8)))
}

func TestUint64ToByteArray(t *testing.T) {
	var a = uint64(10)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, a)
}

func TestUint32Modulo(t *testing.T) {
	aa := uint32(12)
	bb := uint32(10)
	log.Println(aa / bb)
	log.Println(aa % bb)
}
