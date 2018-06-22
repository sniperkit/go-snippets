package textsim

import (
	"hash"

	"github.com/aduston/rabin"
)

type RollingHash interface {
	// Roll rolls in newData and rolls out old data of same byte length
	// as newData. For a given RollingHash instance, this must always be
	// called with the correct number of bytes.
	Roll(newData []byte)
	IsFull() bool
	Sum64() uint64
	Size() int
	Reset()
}

type regHashRollingHash struct {
	regHash hash.Hash64
	buf     *circularBuffer
	size    int
}

func NewRegHashRollingHash(regHash hash.Hash64, size int) RollingHash {
	return &regHashRollingHash{
		regHash: regHash,
		buf:     newCircularBuffer(size),
		size:    size,
	}
}

func (r *regHashRollingHash) Roll(newData []byte) {
	r.buf.addElem(newData)
	if r.buf.isFull() {
		r.regHash.Reset()
		r.buf.write(r.regHash)
	}
}

func (r *regHashRollingHash) IsFull() bool {
	return r.buf.isFull()
}

func (r *regHashRollingHash) Sum64() uint64 {
	if !r.buf.isFull() {
		panic("Cannot ask for sum until full")
	}
	return r.regHash.Sum64()
}

func (r *regHashRollingHash) Size() int {
	return r.size
}

func (r *regHashRollingHash) Reset() {
	r.regHash.Reset()
	r.buf = newCircularBuffer(r.size)
}

type rabinRollingHash struct {
	rabinHash rabin.RollingHash
	buf       *circularBuffer
	size      int
}

func NewRabinRollingHash(rabinHash rabin.RollingHash, size int) RollingHash {
	return &rabinRollingHash{
		rabinHash: rabinHash,
		buf:       newCircularBuffer(size),
		size:      size,
	}
}

func (r *rabinRollingHash) Roll(newData []byte) {
	displaced := r.buf.addElem(newData)
	if displaced != nil {
		r.rabinHash.Roll(displaced, newData)
	} else {
		r.rabinHash.Write(newData)
	}
}

func (r *rabinRollingHash) IsFull() bool {
	return r.buf.isFull()
}

func (r *rabinRollingHash) Sum64() uint64 {
	return r.rabinHash.Sum64()
}

func (r *rabinRollingHash) Size() int {
	return r.size
}

func (r *rabinRollingHash) Reset() {
	r.rabinHash.Reset()
	r.buf = newCircularBuffer(r.size)
}
