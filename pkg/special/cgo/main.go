package main

// #include <stdio.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"
)

// Random ...
func Random() int {
	return int(C.random())
}

// Seed ...
func Seed(i int) {
	C.srandom(C.uint(i))
}

// Print ...
func Print(s string) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))

	C.fputs(cs, (*C.FILE)(C.stdout))
}

func main() {
	fmt.Println(Random())
	Print("hello world")
}
