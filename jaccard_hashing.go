package sticker

import (
	"encoding/gob"
	"fmt"
	"io"
	"math/bits"
	"math/rand"
)

// JaccardHashing is the optimal Densified One Permutation Hashing (DOPH) for estimating Jaccard similarity (Wang+ 2017).
//
// References:
//
// (Wang+ 2017) Y. Wang, A. Shrivastava, and J. Ryu. "FLASH: Randomized Algorithms Accelerated over CPU-GPU for Ultra-High Dimensional Similarity Search." arXiv preprint arXiv:1709.01190, 2017.
type JaccardHashing struct {
	_R        uint
	tables    [][][]uint32
	tableHits [][]uint32
	rng       *rand.Rand
}

// NewJaccardHashing returns an new JaccardHashing.
// Recommended K, L and R are 64, 16 and 64, respectively.
func NewJaccardHashing(K, L, R uint) *JaccardHashing {
	rng := rand.New(rand.NewSource(0))
	tables, tableHits := make([][][]uint32, K), make([][]uint32, K)
	for k := range tables {
		tables[k], tableHits[k] = make([][]uint32, 1<<L), make([]uint32, 1<<L)
	}
	hashing := &JaccardHashing{
		_R:        R,
		tables:    tables,
		tableHits: tableHits,
		rng:       rng,
	}
	return hashing
}

// Add adds the given feature vector as i-th index to the hash tables.
func (hashing *JaccardHashing) Add(vec FeatureVector, i uint32) {
	hvec := hashing.Hash(vec)
	for k, hveck := range hvec {
		if uint(len(hashing.tables[k][hveck])) < hashing._R {
			hashing.tables[k][hveck] = append(hashing.tables[k][hveck], uint32(i))
		} else {
			j := hashing.rng.Intn(int(hashing.tableHits[k][hveck]))
			if uint(j) < hashing._R {
				hashing.tables[k][hveck][j] = i
			}
		}
		hashing.tableHits[k][hveck]++
	}
}

// DecodeJaccardHashingWithGobDecoder decodes JaccardHashing using decoder.
//
// This function returns an error in decoding.
func DecodeJaccardHashingWithGobDecoder(hashing *JaccardHashing, decoder *gob.Decoder) error {
	if err := decoder.Decode(&hashing._R); err != nil {
		return fmt.Errorf("DecodeJaccardHashing: _R: %s", err)
	}
	var K uint
	if err := decoder.Decode(&K); err != nil {
		return fmt.Errorf("DecodeJaccardHashing: K: %s", err)
	}
	hashing.tables = make([][][]uint32, K)
	for k := range hashing.tables {
		if err := decoder.Decode(&hashing.tables[k]); err != nil {
			return fmt.Errorf("DecodeJaccardHashing: #%d table: %s", k, err)
		}
	}
	hashing.tableHits = make([][]uint32, K)
	if err := decoder.Decode(&hashing.tableHits); err != nil {
		return fmt.Errorf("DecodeJaccardHashing: tableHits: %s", err)
	}
	hashing.ResetRng()
	return nil
}

// DecodeJaccardHashing decodes Dataset from r.
// Directly passing *os.File used by a gob.Decoder to this function causes mysterious errors.
// Thus, if users use gob.Decoder, then they should call DecodeJaccardHashingWithGobDecoder.
//
// This function returns an error in decoding.
func DecodeJaccardHashing(hashing *JaccardHashing, r io.Reader) error {
	return DecodeJaccardHashingWithGobDecoder(hashing, gob.NewDecoder(r))
}

// EncodeJaccardHashingWithGobEncoder decodes JaccardHashing using encoder.
//
// This function returns an error in decoding.
func EncodeJaccardHashingWithGobEncoder(hashing *JaccardHashing, encoder *gob.Encoder) error {
	if err := encoder.Encode(hashing._R); err != nil {
		return fmt.Errorf("EncodeJaccardHashing: _R: %s", err)
	}
	if err := encoder.Encode(hashing.K()); err != nil {
		return fmt.Errorf("EncodeJaccardHashing: K: %s", err)
	}
	for k, table := range hashing.tables {
		if err := encoder.Encode(table); err != nil {
			return fmt.Errorf("EncodeJaccardHashing: #%d table: %s", k, err)
		}
	}
	if err := encoder.Encode(hashing.tableHits); err != nil {
		return fmt.Errorf("EncodeJaccardHashing: tableHits: %s", err)
	}
	hashing.ResetRng()
	return nil
}

// EncodeJaccardHashing encodes JaccardHashing to w.
// Directly passing *os.File used by a gob.Encoder to this function causes mysterious errors.
// Thus, if users use gob.Encoder, then they should call EncodeJaccardHashingWithGobEncoder.
//
// This function returns an error in encoding.
func EncodeJaccardHashing(hashing *JaccardHashing, w io.Writer) error {
	return EncodeJaccardHashingWithGobEncoder(hashing, gob.NewEncoder(w))
}

// GobEncode returns the error always, because users should encode large JaccardHashing objects with EncodeJaccardHashing.
func (hashing *JaccardHashing) GobEncode() ([]byte, error) {
	return nil, fmt.Errorf("JaccardHashing should be encoded with EncodeJaccardHashing")
}

// FindNears returns the histogram of the neighbors from the given feature vector.
func (hashing *JaccardHashing) FindNears(vec FeatureVector) KeyCountMap32 {
	hvec := hashing.Hash(vec)
	nears := NewKeyCountMap32(hashing.K() * hashing.R())
	for k, hveck := range hvec {
		for _, idx := range hashing.tables[k][hveck] {
			nears.Inc(idx)
		}
	}
	return nears
}

// Hash returns the K hashed values of the given feature vector.
func (hashing *JaccardHashing) Hash(vec FeatureVector) []uint32 {
	K, L := hashing.K(), hashing.L()
	if len(vec) == 0 {
		// Adhoc solution.
		return make([]uint32, K)
	}
	shiftL := 32 - L
	binSize := uint32(((1 << L) + (K - 1)) / K)
	H0 := make([]uint32, K)
	for k := range H0 {
		H0[k] = ^uint32(0)
	}
	for _, vecpair := range vec {
		key := vecpair.Key
		h := HashUint32(key) >> shiftL
		id := h / binSize
		if H0[id] > h {
			H0[id] = h
		}
	}
	H := make([]uint32, K)
	for k := range H {
		next := H0[k]
		if next == ^uint32(0) {
			for count := uint32(0); next == ^uint32(0); count++ {
				next = H0[(HashUint32((uint32(k)<<8)+count)>>shiftL)/binSize]
			}
		}
		H[k] = next
	}
	return H
}

// K returns the number of the hash tables.
func (hashing *JaccardHashing) K() uint {
	return uint(len(hashing.tables))
}

// L returns the bit-width for backet indices in each hash table.
func (hashing *JaccardHashing) L() uint {
	return uint(bits.Len(uint(len(hashing.tables[0]))) - 1)
}

// R returns the size of a reservoir of each backet.
func (hashing *JaccardHashing) R() uint {
	return hashing._R
}

// ResetRng resets the internal random number generator.
func (hashing *JaccardHashing) ResetRng() {
	hashing.rng = rand.New(rand.NewSource(0))
}

// Summary returns the slice of bucket usage and the bucket size (size of each reservoir) histogram.
func (hashing *JaccardHashing) Summary() (backetUsage []int, backetHist []int) {
	backetUsage = make([]int, len(hashing.tables))
	backetHist = make([]int, hashing._R)
	for t, tblk := range hashing.tables {
		count := 0
		for _, entry := range tblk {
			if len(entry) > 0 {
				backetHist[len(entry)-1]++
				count++
			}
		}
		backetUsage[t] = count
	}
	return
}
