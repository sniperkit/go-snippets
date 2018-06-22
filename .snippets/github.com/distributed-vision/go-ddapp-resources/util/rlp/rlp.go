package rlp

import (
	"errors"

	"github.com/distributed-vision/go-resources/util/hton"
	"github.com/distributed-vision/go-resources/util/ntoh"
)

// EntryType represents the rlp entry type single byte, byte string or list
type EntryType int

const (
	// UNDEFINED is an error type used in errors
	UNDEFINED EntryType = -1
	// BYTE represents a single byte entry
	BYTE EntryType = iota
	// STRING represents a string of bytes
	STRING
	// LIST represents a list containing other entries
	LIST
)

var (
	// ErrCanonLen is returned if the length discovered does not follow the coding rules
	ErrCanonLen = errors.New("rlp: non-canonical length information")
	// ErrBounds is returned for indexing errors on input arrays
	ErrBounds = errors.New("rlp: index out of bounds")
)

func decodeHeader(buffer []byte, index int) (entryType EntryType, hlen uint64, vlen uint64, err error) {

	if len(buffer) <= index {
		return UNDEFINED, 0, 0, ErrBounds
	}

	b := buffer[index]

	switch {
	case b < 0x80:
		// For a single byte whose value is in the [0x00, 0x7F] range, that byte
		// is its own RLP encoding.
		return BYTE, 0, 1, nil
	case b < 0xB8:
		// Otherwise, if a string is 0-55 bytes long,
		// the RLP encoding consists of a single byte with value 0x80 plus the
		// length of the string followed by the string. The range of the first
		// byte is thus [0x80, 0xB7].
		return STRING, 1, uint64(b - 0x80), nil
	case b < 0xC0:
		// If a string is more than 55 bytes long, the
		// RLP encoding consists of a single byte with value 0xB7 plus the length
		// of the length of the string in binary form, followed by the length of
		// the string, followed by the string. For example, a length-1024 string
		// would be encoded as 0xB90400 followed by the string. The range of
		// the first byte is thus [0xB8, 0xBF].
		size := ntoh.UInt(buffer, index, int(b-0xB7))
		if size < 56 {
			err = ErrCanonLen
		}
		return STRING, uint64(b-0xB7) + 1, size, err
	case b < 0xF8:
		// If the total payload of a list
		// (i.e. the combined length of all its items) is 0-55 bytes long, the
		// RLP encoding consists of a single byte with value 0xC0 plus the length
		// of the list followed by the concatenation of the RLP encodings of the
		// items. The range of the first byte is thus [0xC0, 0xF7].
		return LIST, 1, uint64(b - 0xC0), nil
	default:
		// If the total payload of a list is more than 55 bytes long,
		// the RLP encoding consists of a single byte with value 0xF7
		// plus the length of the length of the payload in binary
		// form, followed by the length of the payload, followed by
		// the concatenation of the RLP encodings of the items. The
		// range of the first byte is thus [0xF8, 0xFF].
		size := ntoh.UInt(buffer, index, int(b-0xF7))
		if size < 56 {
			err = ErrCanonLen
		}
		return LIST, uint64(b-0xF7) + 1, size, err
	}
}

func writeHeader(buffer []byte, index int, value int) int {
	if value < 56 {
		buffer[index] = 0x80 + byte(value)
		return 1
	}

	lenlen := hton.UInt(buffer, index+1, uint64(value))
	buffer[index] = 0xB7 + byte(lenlen)
	return lenlen + 1
}

func headerLen(value int) int {
	if value < 56 {
		return 1
	}
	return hton.UIntLen(uint64(value)) + 1
}

func writeEntry(buffer []byte, index int, value []byte) int {
	if len(value) == 1 && value[0] <= 0x7F {
		buffer[index] = value[0]
		return 1
	}

	vlen := len(value)
	hlen := writeHeader(buffer, index, vlen)
	copy(buffer[hlen:], value)
	return hlen + vlen
}

func entryLen(value []byte) int {
	if len(value) == 1 && value[0] <= 0x7F {
		return 1
	}
	vlen := len(value)
	hlen := headerLen(vlen)
	return hlen + vlen
}
