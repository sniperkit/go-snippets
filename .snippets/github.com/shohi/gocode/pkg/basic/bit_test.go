package basic

import (
	"log"
	"strconv"
	"testing"
)

func TestBitShiftRight(t *testing.T) {
	var b uint8
	b = 255

	log.Println(b >> 1)

	b = 128
	log.Println(b & uint8(1))
}

func TestBitShiftLeft(t *testing.T) {
	var b uint8
	b = 1

	log.Println(b)
	log.Println(b << 1)
	log.Println(b << 1)
}

func TestBitClear(t *testing.T) {
	n := 127
	log.Println("b", strconv.FormatInt(int64(n), 2))
	pos := uint(2)
	mask := ^(1 << pos)
	n &= mask

	log.Println("a", strconv.FormatInt(int64(n), 2))
}

func TestBitOperation(t *testing.T) {
	var a, b uint8
	a = 1
	b = 1

	log.Println(a ^ b)
	log.Printf("binary == %b\n", 0xff^(1<<2))
}

func TestBitShift(t *testing.T) {
	bb := 0xff
	shift := uint(7)
	log.Printf("shift %d ==> %b", shift, bb>>shift)

	shift = uint(6)
	log.Printf("shift %d ==> %b", shift, bb>>shift)
}
