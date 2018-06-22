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

var encoders = map[encodertype.EncoderType]func(src []byte) string{
	encodertype.BASE62:    base62.Encode,
	encodertype.BASE36:    base36.Encode,
	encodertype.BASE58:    base58.Encode,
	encodertype.BASE10:    base10.Encode,
	encodertype.BASE64URL: base64url.Encode,
	encodertype.HEX:       hex.EncodeToString}

func Encode(toEncode []byte, encoderType encodertype.EncoderType) (string, error) {
	encoder := encoders[encoderType]

	if encoder != nil {
		return encoder(toEncode), nil
	}

	if encoderType == encodertype.UTF8 {
		return string(toEncode), nil
	}

	return "", errors.New("Unsupported encoding: " + encoderType.String())
}
