package main

import "fmt"

//
// go run -ldflags "-X main.xyz=abc" command_link.go
//
// go run -ldflags "-X main.xyz=`date -u +%Y.%m.%d-%H:%M:%S`" command_link.go
//

var xyz string

func main() {
	fmt.Println(xyz)
}