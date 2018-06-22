package cgzip

import "C"

import (
	"hash"
	"unsafe"
)

type adler32Hash struct {
	adler C.uLong
}

func NewAdler32() hash.Hash32 {
	a := &adler32Hash{}
	a.Reset()
	return a
}

func (a *adler32Hash) Write(p []byte) (n int, err error) {
	if len(p) > 0 {
		a.adler = C.adler32(a.adler, (*C.Bytef)(unsafe.Pointer(&p[0])), (C.uInt)(len(p)))
	}
	return len(p), nil
}