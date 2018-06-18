package main

import (
	"crypto/sha256"
	"fmt"

	"github.com/gitchander/go-examples/crypto/hex"
)

func main() {
	Sample1()
	Sample2()
}

func Sample1() {

	bs1 := []byte("The quick brown fox jumps over the lazy dog")
	bs2 := []byte("The quick brown fox jumps over the lazy cog")

	h1 := sha256.Sum256(bs1)
	h2 := sha256.Sum256(bs2)

	fmt.Println(hex.HexQuad(h1[:]))
	fmt.Println(hex.HexQuad(h2[:]))
}

func Sample2() {

	bs := []byte("The quick brown fox jumps over the lazy dog")

	h := sha256.New()

	h.Write(bs[:5])
	h.Write(bs[5:8])
	h.Write(bs[8:14])
	h.Write(bs[14:])
	fmt.Println(hex.HexQuad(h.Sum(nil)))

	h.Reset()
	h.Write(bs)
	fmt.Println(hex.HexQuad(h.Sum([]byte{0xAA, 0xBB, 0xCC, 0xDD})))
}
