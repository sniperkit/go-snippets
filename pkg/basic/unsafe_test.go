package basic

import (
	"fmt"
	"log"
	"math"
	"strings"
	"testing"
	"unsafe"
)

func TestUnsafeSizeOf(t *testing.T) {
	var bb int
	aa := make([]int, math.MaxUint32)
	log.Println(int(unsafe.Sizeof(&bb)) * len(aa) / (1024 * 1024))
}

func TestUnsafeSizeOfForBaseStruct(t *testing.T) {
	a := int(123)
	b := int64(123)
	c := strings.Repeat("foo", 100)
	d := struct {
		FieldA float32
		FieldB string
	}{0, "bar"}
	e := struct{}{}

	fmt.Printf("a: %T, %d\n", a, unsafe.Sizeof(a))
	fmt.Printf("b: %T, %d\n", b, unsafe.Sizeof(b))
	fmt.Printf("c: %T, %d\n", c, unsafe.Sizeof(c))
	fmt.Printf("d: %T, %d\n", d, unsafe.Sizeof(d))
	fmt.Printf("e: %T, %d\n", d, unsafe.Sizeof(e))
}

func TestUnsafeSizeOfDifferentType(t *testing.T) {
	var s struct{}
	log.Println(unsafe.Sizeof(s))
	var i interface{}
	log.Println(unsafe.Sizeof(i))
	var b bool
	log.Println(unsafe.Sizeof(b))
}
