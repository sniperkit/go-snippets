package base62

import (
	"fmt"

	"github.com/distributed-vision/go-resources/encoding/basex"
)

var base62 = basex.NewEncoder("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func Decode(toDecode string) ([]byte, error) {
	return base62.Decode(toDecode)
}

func MustDecode(toDecode string) []byte {
	result, err := base62.Decode(toDecode)

	if err != nil {
		panic(fmt.Sprintf("base62.decode failed: %s", err))
	}

	return result
}

func Encode(toEncode []byte) string {
	return base62.Encode(toEncode)
}
