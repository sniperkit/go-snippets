package main

import (
	"crypto/sha1"
	"fmt"
)

func main() {
	s := "sha1 this string"
	// The pattern for generating a hash is sha1.New(), sha1.Write(bytes), then sha1.Sum([]byte{}). Here we start with a new hash.
	h := sha1.New()
	// Write expects bytes. If you have a string s, use []byte(s) to coerce it to bytes.
	// This gets the finalized hash result as a byte slice. The argument to Sum can be used to append to an existing byte slice: it usually isn’t needed.
	h.Write([]byte(s))

	bs := h.Sum(nil)

	fmt.Println(s)
	fmt.Printf("%x\n", bs)
}
