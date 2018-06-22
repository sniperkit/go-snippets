package base36

import "github.com/distributed-vision/go-resources/encoding/basex"

var base36=basex.NewEncoder("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func Decode(toDecode string) ([]byte, error) {
  return base36.Decode(toDecode)
}

func Encode(toEncode []byte) string  {
  return base36.Encode(toEncode)
}
