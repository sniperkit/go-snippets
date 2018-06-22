package encoding

import (
	"encoding/hex"
	"errors"

	"github.com/distributed-vision/go-resources/encoding/base10"
	"github.com/distributed-vision/go-resources/encoding/base36"
	"github.com/distributed-vision/go-resources/encoding/base58"
	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/encoding/base64url"
	"github.com/distributed-vision/go-resources/encoding/encodertype"
)

var decoders = map[encodertype.EncoderType]func(src string) ([]byte, error){
	encodertype.BASE62:    base62.Decode,
	encodertype.BASE36:    base36.Decode,
	encodertype.BASE58:    base58.Decode,
	encodertype.BASE10:    base10.Decode,
	encodertype.BASE64URL: base64url.Decode,
	encodertype.HEX:       hex.DecodeString,
}

func Decode(toDecode string, encoderType encodertype.EncoderType) ([]byte, error) {
	decoder := decoders[encoderType]

	if decoder != nil {
		return decoder(toDecode)
	}

	if encoderType == encodertype.UTF8 {
		return []byte(toDecode), nil
	}

	return nil, errors.New("Unsupported encoding: " + encoderType.String())
}
