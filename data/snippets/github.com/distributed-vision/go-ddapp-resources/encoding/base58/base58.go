package base58

import "github.com/distributed-vision/go-resources/encoding/basex"

var base58=basex.BTCEncoder

func Decode(toDecode string) ([]byte, error) {
  return base58.Decode(toDecode)
}

func Encode(toEncode []byte) string  {
  return base58.Encode(toEncode)
}
