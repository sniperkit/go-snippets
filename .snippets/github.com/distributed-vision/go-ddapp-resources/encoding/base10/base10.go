package base10

import "github.com/distributed-vision/go-resources/encoding/basex"

var base10 = basex.NewEncoder("0123456789")

func Decode(toDecode string) ([]byte, error) {
	return base10.Decode(toDecode)
}

func Encode(toEncode []byte) string {
	return base10.Encode(toEncode)
}
