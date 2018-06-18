package main

import (
	"crypto/rand"
	"fmt"
)

func randStr(strSize int, randType string) string {

	var dictionary string

	switch randType {
	case "alphanum":
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	case "alpha":
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	case "number":
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

func main() {

	fmt.Println("Alphanum : ", randStr(16, "alphanum"))

	fmt.Println("Alpha : ", randStr(16, "alpha"))

	fmt.Println("Numbers : ", randStr(16, "number"))

}
