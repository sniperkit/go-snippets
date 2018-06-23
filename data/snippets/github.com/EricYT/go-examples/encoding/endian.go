package encoding

import "unsafe"

// int endian() {
//     union {
//         int i;
//         char c[sizeof(int)];
//     } x;
//     x.i = 1;
//     return x.c[0] == 1 ? 0 : 1;
// }
// int endian_concision() {
//     int i = 1;
//     return *(char *)&i == 1 ? 0 : 1;
// }
import "C"

type EndianType int

const (
	LITTLE_ENDIAN EndianType = iota
	BIG_ENDIAN
)

func Endian() EndianType {
	var x uint32 = 0x01020304
	if *(*byte)(unsafe.Pointer(&x)) == 0x01 {
		return BIG_ENDIAN
	}
	return LITTLE_ENDIAN
}

func CEndian() EndianType {
	return EndianType(C.endian())
}

func CEndianConcision() EndianType {
	return EndianType(C.endian_concision())
}
