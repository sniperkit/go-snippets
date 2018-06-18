package life

import (
	"encoding/binary"
)

type hashBytes [8]byte

func hashPointToBytes(p Point) (h hashBytes) {
	var byteOrder = binary.LittleEndian
	byteOrder.PutUint32(h[0:], uint32(p.X))
	byteOrder.PutUint32(h[4:], uint32(p.Y))
	return
}

func hashPointToUint64(p Point) uint64 {
	var (
		a = uint32(p.X)
		b = uint32(p.Y)
	)
	return (uint64(a) << 32) | uint64(b)
}

func hashPointToUint64_v2(p Point) uint64 {
	return (uint64(p.X) << 32) | (uint64(p.Y) & 0xFFFFFFFF)
}
