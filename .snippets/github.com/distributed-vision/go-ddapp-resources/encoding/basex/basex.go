// Taken from github.com/jbenet/go-base58
// Copyright (c) 2013-2014 Conformal Systems LLC.
// Modified by Mark Holt (mark.holt@distributed.vision)
// to handle multiple radix alphabets

package basex

import (
	"errors"
	"math/big"
)

// alphabet is the modified base58 alphabet used by Bitcoin.
var BTCEncoder = NewEncoder("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
var FlickrEncoder = NewEncoder("123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ")

type Encoder struct {
	alphabet    string
	radix       *big.Int
	alphabetMap map[byte]*big.Int
}

func NewEncoder(alphabet string) *Encoder {
	return &Encoder{
		alphabet:    alphabet,
		radix:       big.NewInt(int64(len(alphabet))),
		alphabetMap: newAlphapetMap(alphabet)}
}

func newAlphapetMap(alphabet string) map[byte]*big.Int {
	alphabetMap := make(map[byte]*big.Int, len(alphabet))

	for i := 0; i < len(alphabet); i++ {
		alphabetMap[alphabet[i]] = big.NewInt(int64(i))
	}

	return alphabetMap
}

var bigZero = big.NewInt(0)

// DecodeAlphabet decodes a modified base58 string to a byte slice, using alphabet.
func (e *Encoder) Decode(b string) ([]byte, error) {
	answer := big.NewInt(0)
	j := big.NewInt(1)

	for i := len(b) - 1; i >= 0; i-- {
		idx, ok := e.alphabetMap[b[i]]

		if !ok {
			return nil, errors.New("Invalid encoding: unexpected character: '" + string(b[i]) + "'")
		}

		tmp1 := big.NewInt(0)
		tmp1.Mul(j, idx)

		answer.Add(answer, tmp1)
		j.Mul(j, e.radix)
	}

	tmpval := answer.Bytes()

	var numZeros int
	for numZeros = 0; numZeros < len(b); numZeros++ {
		if b[numZeros] != e.alphabet[0] {
			break
		}
	}
	flen := numZeros + len(tmpval)
	val := make([]byte, flen, flen)
	copy(val[numZeros:], tmpval)

	return val, nil
}

// Encode encodes a byte slice to a modified base58 string, using alphabet
func (e *Encoder) Encode(b []byte) string {
	x := new(big.Int)
	x.SetBytes(b)

	answer := make([]byte, 0, len(b)*136/100)
	for x.Cmp(bigZero) > 0 {
		mod := new(big.Int)
		x.DivMod(x, e.radix, mod)
		answer = append(answer, e.alphabet[mod.Int64()])
	}

	// leading zero bytes
	for _, i := range b {
		if i != 0 {
			break
		}
		answer = append(answer, e.alphabet[0])
	}

	// reverse
	alen := len(answer)
	for i := 0; i < alen/2; i++ {
		answer[i], answer[alen-1-i] = answer[alen-1-i], answer[i]
	}

	return string(answer)
}
