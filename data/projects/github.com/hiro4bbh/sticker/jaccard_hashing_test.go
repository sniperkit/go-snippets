package sticker

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/hiro4bbh/go-assert"
)

func TestJaccardHashingHash(t *testing.T) {
	countMatches := func(x, y []uint32) int {
		n := len(x)
		if n > len(y) {
			n = len(y)
		}
		count := 0
		for i := 0; i < n; i++ {
			if x[i] == y[i] {
				count++
			}
		}
		return count
	}
	vecs := FeatureVectors{
		FeatureVector{}, FeatureVector{KeyValue32{0, 1.0}},
		FeatureVector{KeyValue32{1, 1.0}}, FeatureVector{KeyValue32{1, 1.0}, KeyValue32{3, 2.0}}, FeatureVector{KeyValue32{0, 4.0}, KeyValue32{1, 1.0}, KeyValue32{2, 3.0}, KeyValue32{3, 2.0}},
	}
	assertAllMatches := func(K uint, indexPairJaccards [][3]int) {
		hashing := NewJaccardHashing(K, 10, 8)
		hashes := make([][]uint32, len(vecs))
		for i, veci := range vecs {
			hashes[i] = hashing.Hash(veci)
		}
		for _, indexPairJaccard := range indexPairJaccards {
			i1, i2, J := indexPairJaccard[0], indexPairJaccard[1], indexPairJaccard[2]
			goassert.New(t, [3]int{i1, i2, J}).Equal([3]int{i1, i2, countMatches(hashes[i1], hashes[i2])})
		}
	}
	assertAllMatches(16, [][3]int{
		{0, 1, 0},  // Answer: 0/1
		{2, 3, 8},  // Answer: 1/2
		{2, 4, 4},  // Answer: 1/4
		{3, 4, 10}, // Answer: 2/4
	})
	assertAllMatches(256, [][3]int{
		{0, 1, 0},
		{2, 3, 132},
		{2, 4, 65},
		{3, 4, 120},
	})
	assertAllMatches(4096, [][3]int{
		{0, 1, 0},
		{2, 3, 1991},
		{2, 4, 971},
		{3, 4, 2036},
	})
}

func TestJaccardHashing(t *testing.T) {
	x1 := FeatureVector{KeyValue32{1, 1.0}, KeyValue32{2, 1.0}}
	x2 := FeatureVector{KeyValue32{1, 1.0}, KeyValue32{2, 1.0}, KeyValue32{5, 1.0}, KeyValue32{6, 1.0}}
	hashing := NewJaccardHashing(16, 10, 8)
	goassert.New(t, uint(16)).Equal(hashing.K())
	goassert.New(t, uint(10)).Equal(hashing.L())
	goassert.New(t, uint(8)).Equal(hashing.R())
	hashing.Add(x1, 0)
	hashing.Add(x2, 1)
	goassert.New(t, map[uint32]uint32{0: 16, 1: 6}).Equal(hashing.FindNears(x1).Map())
	goassert.New(t, map[uint32]uint32{0: 6, 1: 16}).Equal(hashing.FindNears(x2).Map())
	goassert.New(t, []int{2, 2, 2, 1, 1, 2, 2, 1, 1, 2, 1, 2, 2, 2, 1, 2}, []int{20, 6, 0, 0, 0, 0, 0, 0}).Equal(hashing.Summary())
	// Check the reservoir boundedness.
	for t := uint32(2); t < 128; t++ {
		hashing.Add(x1, t)
	}
	// FIXME: Heavy conflicts would decrease the estimated Jaccard similarity due to reservoir sampling.
	//goassert.New(t, map[uint32]uint32{0: 16, 1: 7}).Equal(hashing.FindNears(x1))
	goassert.New(t, []int{2, 2, 2, 1, 1, 2, 2, 1, 1, 2, 1, 2, 2, 2, 1, 2}, []int{10, 0, 0, 0, 0, 0, 0, 16}).Equal(hashing.Summary())
	// Check gob encoding/decoding.
	var buf bytes.Buffer
	goassert.New(t).SucceedWithoutError(EncodeJaccardHashing(hashing, &buf))
	var decodedHashing JaccardHashing
	goassert.New(t).SucceedWithoutError(DecodeJaccardHashing(&decodedHashing, &buf))
	goassert.New(t, hashing).Equal(&decodedHashing)
	// gob.Decoder.Decode won't call JaccardHashing.GobDecode, because the encoder did not encode JaccardHashing.
	goassert.New(t, "JaccardHashing should be encoded with EncodeJaccardHashing").ExpectError(gob.NewEncoder(&buf).Encode(&decodedHashing))
	// Check that an empty FeatureVector can be hashed safely.
	hashing0 := NewJaccardHashing(16, 10, 8)
	hashing0.Add(FeatureVector{}, 0)
	goassert.New(t, []uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}).Equal(hashing0.Hash(FeatureVector{}))
	goassert.New(t, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, []int{16, 0, 0, 0, 0, 0, 0, 0}).Equal(hashing0.Summary())
}
