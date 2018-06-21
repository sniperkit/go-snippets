package main

import (
	"encoding/base64"
	"log"
	"os"
)

func main() {
	encoder()
	encodeToString()
	encoding()
}

func encoder() {
	input := []byte("你好，我是ETO")
	encoder := base64.NewEncoder(base64.StdEncoding, os.Stdout)
	encoder.Write(input)
	// Must close the encoder when finished to flush any partial blocks.
	// If you comment out the following line, the last partial block "r"
	// won't be encoded.
	encoder.Close()
	// Output:
	// Zm9vAGJhcg==
}

func encodeToString() {
	msg := "Hello, 世界"
	encoded := base64.StdEncoding.EncodeToString([]byte(msg))
	log.Println(encoded)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Println("decode error:", err)
		return
	}
	log.Println(string(decoded))
}

func encoding() {
	//This code will effect the generated result
	const encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	const encodeURL = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	log.Println(len(encodeStd))
	log.Println(len([]byte(encodeStd)))
	var funnyEncoding = base64.NewEncoding(encodeStd).WithPadding(rune('@'))

	encoded := funnyEncoding.EncodeToString([]byte("Hello world!"))
	log.Println(encoded)
	decoded, err := funnyEncoding.DecodeString(encoded)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(decoded))
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
