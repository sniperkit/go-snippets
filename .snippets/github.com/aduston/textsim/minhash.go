package textsim

import (
	"encoding/binary"
	"hash"
	"math"
	"math/big"
	"math/rand"
	"unicode"
)

const kMersennePrime = (uint64(1) << 61) - 1

var kBigMersennePrime = big.NewInt(0).SetUint64(kMersennePrime)

type Hash64 func([]byte) uint64

// LinearHashFunction is a direct translation of
// http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.121.8215&rep=rep1&type=pdf
// and appears to be the approach most commonly used in practice.
type LinearHashFunction struct {
	A    uint64
	B    uint64
	bigA *big.Int
	bigB *big.Int
}

func (l *LinearHashFunction) Hash(x uint64) uint64 {
	bigX := big.NewInt(0).SetUint64(x)
	bigX.Mul(bigX, l.bigA)
	bigX.Add(bigX, l.bigB)
	bigX.Mod(bigX, kBigMersennePrime)
	return bigX.Uint64()
}

func NewRandomLinearFunction() *LinearHashFunction {
	// TODO: seed
	a := uint64(1 + rand.Int63n(int64(kMersennePrime-1)))
	b := uint64(rand.Int63n(int64(kMersennePrime)))
	return &LinearHashFunction{
		A:    a,
		B:    b,
		bigA: big.NewInt(0).SetUint64(a),
		bigB: big.NewInt(0).SetUint64(b),
	}
}

// Tokenize takes a utf-8 string. Every maximal aphanumeric sequence
// is considered a term and is hashed using wordHash to generate tokens.
func Tokenize(text string, wordHash hash.Hash64) (tokens []uint64) {
	inWord := false
	tokens = make([]uint64, 0, len(text)/6)
	for _, runeValue := range text {
		if unicode.IsLetter(runeValue) || unicode.IsNumber(runeValue) {
			inWord = true
			wordHash.Write(toBytes32(uint32(runeValue)))
		} else if inWord {
			tokens = append(tokens, wordHash.Sum64())
			wordHash.Reset()
			inWord = false
		}
	}
	if inWord {
		tokens = append(tokens, wordHash.Sum64())
	}
	return
}

func toBytes64(num uint64) []byte {
	return []byte{
		byte(num >> 56), byte(num >> 48), byte(num >> 40), byte(num >> 32),
		byte(num >> 24), byte(num >> 16), byte(num >> 8), byte(num)}
}

func toBytes32(num uint32) []byte {
	return []byte{byte(num >> 24), byte(num >> 16), byte(num >> 8), byte(num)}
}

func ConvertToShingles(tokens []uint64, rollingHash RollingHash) []uint64 {
	if len(tokens) < rollingHash.Size() {
		panic("Can't call this function with a number of tokens less than shingleSize.")
	}
	shingles := make([]uint64, len(tokens)-rollingHash.Size()+1)
	for i, token := range tokens {
		rollingHash.Roll(toBytes64(token))
		if i >= rollingHash.Size()-1 {
			shingles[i-rollingHash.Size()+1] = rollingHash.Sum64()
		}
	}
	return shingles
}

// ConvertToMinHashes directly converts a slice of tokens to a slice
// of k minhashes by applying k rolling hashes. In practice, the k
// rolling hashes are Rabin fingerprints using k different irreducible
// polynomials over GF(2) of degree < 64.
func ConvertToMinHashes(tokens []uint64, rollingHashes []RollingHash) []uint64 {
	minimums := make([]uint64, len(rollingHashes))
	for i := range minimums {
		minimums[i] = math.MaxUint64
	}
	for i, token := range tokens {
		for j, rollingHash := range rollingHashes {
			rollingHash.Roll(toBytes64(token))
			if i >= rollingHash.Size()-1 {
				v := minimums[j]
				hv := rollingHash.Sum64()
				if hv < v {
					minimums[j] = hv
				}
			}
		}
	}
	return minimums
}

func CalcMinHashesDGryski(shingles []uint64, h1, h2 Hash64, size int) []uint64 {
	minimums := make([]uint64, size)
	for i := range minimums {
		minimums[i] = math.MaxUint64
	}
	for _, shingle := range shingles {
		shingleBytes := toBytes64(shingle)
		v1 := h1(shingleBytes)
		v2 := h2(shingleBytes)

		for i, v := range minimums {
			hv := v1 + uint64(i)*v2
			if hv < v {
				minimums[i] = hv
			}
		}
	}
	return minimums
}

func MakePermHashes(hash1, hash2 hash.Hash64) (h1, h2 Hash64) {
	// TODO: make truly random
	r := rand.New(rand.NewSource(int64(42)))
	b := binary.LittleEndian
	b1 := make([]byte, 8)
	b2 := make([]byte, 8)
	b.PutUint64(b1, uint64(r.Int63()))
	b.PutUint64(b2, uint64(r.Int63()))
	h1 = func(b []byte) uint64 {
		hash1.Reset()
		hash1.Write(b1)
		hash1.Write(b)
		return hash1.Sum64()
	}
	h2 = func(b []byte) uint64 {
		hash2.Reset()
		hash2.Write(b2)
		hash2.Write(b)
		return hash2.Sum64()
	}
	return
}

func CalcMinHashesLinear(shingles []uint64, linearHashes []*LinearHashFunction) []uint64 {
	minimums := make([]uint64, len(linearHashes))
	for i := range minimums {
		minimums[i] = math.MaxUint64
	}
	for _, shingle := range shingles {
		for i, v := range minimums {
			hv := linearHashes[i].Hash(shingle)
			if hv < v {
				minimums[i] = hv
			}
		}
	}
	return minimums
}

func GenerateLinearMinHashParms(num int) []*LinearHashFunction {
	linearHashFunctions := make([]*LinearHashFunction, num)
	for i := 0; i < num; i++ {
		linearHashFunctions[i] = NewRandomLinearFunction()
	}
	return linearHashFunctions
}
