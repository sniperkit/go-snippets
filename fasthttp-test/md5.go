package main

import (
	"crypto/md5"
	"fmt"
)

func main() {
	data := []byte("hello,world")
	md5sum := md5.Sum(data)
	for _, c := range md5sum {
		fmt.Printf("%x", c)
	}
	fmt.Println("")
	//log.Printf("data: %s md5sum: %s\n", string(data), string(md5sum))
}
