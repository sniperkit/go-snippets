package pointer

import (
	"log"
	"testing"
)

func TestDereference(t *testing.T) {

	var strPtr *string
	str := "hello"

	// *strPtr = "hello"
	strPtr = &str
	*strPtr = str

	log.Println(strPtr)
	log.Println(*strPtr)
}

func TestPointer(t *testing.T) {
	type MyT struct {
		val string
	}
	var aa *MyT
	bb := MyT{val: "hello"}
	aa = &bb
	f := func() MyT {
		return *aa
	}

	cc := f()
	cc.val = "world"

	log.Printf("value, aa: %v, bb: %v, cc: %v", aa.val, bb.val, cc.val)
}

func TestPointerZeroValue(t *testing.T) {
	// case 1 - basic type - nil
	var ss *string
	log.Printf("zero value of string pointer: %v", ss)

	// case 2 - complex type - nil
	type MyT struct {
		val string
	}

	var m *MyT
	log.Printf("zero value of complex type's pointer: %v", m)
}
