package basew

import "github.com/distributed-vision/go-resources/encoding/basex"

var basew = basex.NewEncoder("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_")

func Decode(toDecode string) ([]byte, error) {
	return basew.Decode(toDecode)
}

func Encode(toEncode []byte) string {
	return basew.Encode(toEncode)
}
