package main

import (
	"crypto/aes"
	"fmt"
)

func main() {

	key := make([]byte, 32)

	b, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(b.BlockSize())
}
