package main

//extern void SayHello(_GoString_ s1, _GoString_ s2);
import "C"
import (
	"fmt"
)

//export SayHello
func SayHello(s string, t string) {
	fmt.Println(s + t)
}

func main() {
	C.SayHello("GopherChina", "2018")
}
