package encoding

import (
	"errors"
	"fmt"
)

// hex bytes
var hextable []byte = []byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f',
}

// Encode implements hexadecimal encoding.
func Encode(dst, src []byte) int {
	for i, v := range src {
		dst[i*2] = hextable[v>>4]
		dst[i*2+1] = hextable[v&0x0F]
	}
	return len(src) * 2
}

var ErrWrongLen error = errors.New("decode source bytes length is wrong")

func Decode(dst, src []byte) (int, error) {
	if len(src)%2 == 1 {
		return 0, ErrWrongLen
	}

	for i := 0; i < len(src)/2; i++ {
		a, ok := fromHexChar(src[i*2])
		if !ok {
			return 0, fmt.Errorf("encoding/hex: invalid byte(%#U)", src[i*2])
		}
		b, ok := fromHexChar(src[i*2+1])
		if !ok {
			return 0, fmt.Errorf("encoding/hex: invalid byte(%#U)", src[i*2+1])
		}
		dst[i] = (a << 4) | b
	}

	return len(src) / 2, nil
}

func fromHexChar(c byte) (byte, bool) {
	switch {
	case c >= '0' && c <= '9':
		return c - '0', true
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10, true
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10, true
	}
	return 0, false
}

func DecodedLen(len int) int { return len / 2 }

func EncodeToString(src []byte) string {
	dst := make([]byte, len(src)*2)
	Encode(dst, src)
	return string(dst)
}

func DecodeString(s string) ([]byte, error) {
	src := []byte(s)
	dst := make([]byte, DecodedLen(len(src)))
	_, err := Decode(dst, src)
	if err != nil {
		return nil, err
	}
	return dst, nil
}
