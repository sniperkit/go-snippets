package slice

import (
	"fmt"
	"log"
	"testing"
)

func TestSlice(t *testing.T) {
	var aa []int
	for i := 0; i < 5; i++ {
		aa = append(aa, 0)
	}

	log.Println(aa)
}

func TestSliceCapAndLen(t *testing.T) {
	aa := make([]byte, 5)
	log.Println(len(aa), cap(aa))

	cc := make([]byte, 10, 20)
	log.Println(len(cc), cap(cc))

	var bb []byte
	log.Println(len(bb), cap(bb), bb, bb == nil)
	bb = append(bb, 0x10)
	log.Println(len(bb), cap(bb), bb, bb == nil)
}

func TestSliceAppend(t *testing.T) {
	var aa []int
	aa = append(aa, 10, 20)
	log.Println(aa)

	aa = make([]int, 5)
	aa = append(aa, 30, 40)
	log.Println(aa)
}

func TestSliceInit(t *testing.T) {
	var aa []int
	log.Printf("default value for slice: %v", aa)

	var bb []int = nil
	log.Printf("value for nil-initialized slice: %v", bb)
}

func TestSliceInitialization(t *testing.T) {
	// initialize slice without length parameter will cause error
	// aa := make([]int)
	aa := make([]int, 10)
	log.Println(aa)
}

func TestNilSliceTraverse(t *testing.T) {
	var a []*int

	for k, v := range a {
		log.Println(k, v)
	}

	bb := []string{"hello", "world"}
	for k := range bb {
		log.Println("key: ", k)
	}

	for k, v := range bb {
		log.Println("key: ", k, ", value: ", v)
	}
}

func TestSlicePrint(t *testing.T) {
	fn := func(strs ...string) {
		fmt.Printf("slice: %v", strs)
	}
	var aa []string
	fn(aa...)

	aa = append(aa, []string{"a", "b", "c"}...)
	fn(aa...)
}

func TestSliceSub(t *testing.T) {
	c := "hello world"
	log.Printf("string ==> %s", c[:0])
}

func TestEmptyAndNilSlice(t *testing.T) {
	// only way to declare `nil` slice
	var s []int
	log.Printf("s is nil ==> %v", s == nil)
	log.Printf("string ==> %v", s)

	s = []int{}
	log.Printf("s is nil ==> %v, %v", s == nil, len(s))

	s = make([]int, 0)
	log.Printf("s is nil ==> %v, %v", s == nil, len(s))
}
