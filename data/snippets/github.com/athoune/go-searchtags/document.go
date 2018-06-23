package main

import (
	"github.com/pmylund/go-bitset"
)

type documents struct {
	tags  []*bitset.Bitset64
	size  uint64
	bonus map[uint32]*bitset.Bitset64
}

func (self *documents) Set(pos uint32, tag *bitset.Bitset64) {
	self.tags[pos] = tag
}

func NewDocuments(docs uint32, tags uint64) documents {
	return documents{
		make([]*bitset.Bitset64, docs),
		tags,
		make(map[uint32]*bitset.Bitset64, 5)}
}

func NewDocument(ids []uint64) *bitset.Bitset64 {
	b := bitset.New64(0)
	for _, id := range ids {
		b.Set(id)
	}
	return b
}
