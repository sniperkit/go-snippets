package main

import (
	"fmt"
	"github.com/pmylund/go-bitset"
	"math/rand"
	"testing"
)

func TestSearch(t *testing.T) {
	var DOCS_SIZE uint32 = 50000
	var TAGS_SIZE uint64 = 1024
	docs := NewDocuments(DOCS_SIZE, TAGS_SIZE)
	var i uint32
	for i = 0; i < DOCS_SIZE; i++ {
		d := bitset.New64(TAGS_SIZE)
		d.Set(1)
		for j := 0; j < 200; j++ {
			d.Set(uint64(rand.Int63n(int64(TAGS_SIZE))))
		}
		docs.Set(i, d)
	}
	d := bitset.New64(TAGS_SIZE)
	d.Set(1)
	for j := 0; j < 200; j++ {
		d.Set(uint64(rand.Int63n(int64(TAGS_SIZE))))
	}
	r := docs.Score(d, 0.2)
	fmt.Println("documents found", len(r))
}
