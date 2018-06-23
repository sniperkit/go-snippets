package main

import (
	"fmt"
	"math"
	"crypto/md5"
)

func main() {
	fmt.Printf("Math pi = %f\n", math.Pi)
	data := []byte("hello")
	h := md5.New()
	fmt.Printf("%x", h.Sum(data))
}
