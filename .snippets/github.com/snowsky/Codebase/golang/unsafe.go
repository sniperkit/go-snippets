//http://my.oschina.net/achun/blog/122540?p=1

package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	a := []byte("test")
	fmt.Println(String(a))
}

func String(b []byte) (s string) {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pstring.Data = pbytes.Data
	pstring.Len = pbytes.Len
	return
}
