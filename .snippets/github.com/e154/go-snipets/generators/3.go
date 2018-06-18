package main

import (
	"encoding/base64"
	"crypto/rand"
	"fmt"
)

func main() {
	size := 32 // change the length of the generated random string here

	rb := make([]byte,size)
	_, err := rand.Read(rb)


	if err != nil {
		fmt.Println(err)
	}

	rs := base64.URLEncoding.EncodeToString(rb)

	fmt.Println(rs)
}
