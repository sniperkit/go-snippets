package encodertype

import (
	"errors"
	"strings"
)

type EncoderType int

const (
	BASE62 EncoderType = iota
	BASE36
	BASE58
	BASE10
	BASE64URL
	HEX
	UTF8
)

func (this EncoderType) String() string {
	switch this {
	case BASE62:
		return "base62"
	case BASE36:
		return "base36"
	case BASE58:
		return "base58"
	case BASE10:
		return "base10"
	case BASE64URL:
		return "base64url"
	case HEX:
		return "hex"
	case UTF8:
		return "utf8"
	default:
		return "invalid"
	}
}

func Parse(value string) (EncoderType, error) {
	switch strings.ToUpper(value) {
	case "BASE62":
		return BASE62, nil
	case "BASE36":
		return BASE36, nil
	case "BASE58":
		return BASE58, nil
	case "BASE10":
		return BASE10, nil
	case "BASE64URL":
		return BASE64URL, nil
	case "HEX":
		return HEX, nil
	case "UTF8":
		return UTF8, nil
	default:
		return -1, errors.New("Unknown encoding type: " + value)
	}
}
