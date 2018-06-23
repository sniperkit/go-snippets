/*
Copyright 2013 Zack Pierce.
Use of this source code is governed by a MIT-style
license that can be found in the LICENSE file.
*/
/*
Package stralgo/bytewise implements various string algorithms
with an emphasis on similarity metrics, implemented in per-byte
fashion.

This bytewise approach is suited for speedy comparisons
when the target strings contain no multi-byte runes.
*/
package bytewise

import (
	"errors"
	"unicode"
)

// HammingDistance calculates the Hamming distance between
// two strings of equal length, bytewise.
//
// The Hamming distance is the total number of indices
// at which the corresponding bytes are different.
// The higher the result, the more different the strings.

// See: http://en.wikipedia.org/wiki/Hamming_distance
//
// Returns an error if the string lengths are not equal.
//
// Note that this algorithm implementation operates upon
// individual bytes, and does not account for multibyte
// unicode runes.
func HammingDistance(a, b string) (uint, error) {
	aLen := len(a)
	bLen := len(b)
	if aLen != bLen {
		return 0, errors.New("Hamming distance is undefined between strings of unequal length.")
	}
	var d uint
	for i := 0; i < aLen; i++ {
		if a[i] != b[i] {
			d++
		}
	}
	return d, nil
}

// DiceCoefficent calculates the simiarlity of two
// strings per the Sorensen-Dice coefficient, bytewise.
//
// The resulting value is scaled between 0 and 1.0,
// and a higher value means a higher similarity.
//
// This algorithm is also known as the Sorensen Index, and
// is very close to the White Similarity metric, with the key
// distinctions that DiceCoefficient does not differentiate
// between whitespace and other characters and also does not
// account for bigram frequency count differences between
// the compared strings.
//
// See: http://en.wikipedia.org/wiki/Sorensen-Dice_coefficient
//
// Note that this algorithm implementation operates upon
// individual bytes and does not account for multibyte
// unicode runes.
//
// Returns an error if both of the input strings
// contain less than two bytes.
func DiceCoefficient(a, b string) (float64, error) {
	aLimit := len(a) - 1
	bLimit := len(b) - 1
	if aLimit < 1 && bLimit < 1 {
		return 0, errors.New("At least one of the input strings must have a length of 2 or greater for the bigram-based DiceCoefficient to be calculated.")
	}
	aSet := make(map[string]bool, aLimit)
	totalBigrams := 0.0
	var bigram string
	for i := 0; i < aLimit; i++ {
		bigram = a[i : i+2]
		if !aSet[bigram] {
			totalBigrams++
			aSet[bigram] = true
		}
	}

	bSet := make(map[string]bool, bLimit)
	sharedBigrams := 0.0
	for i := 0; i < bLimit; i++ {
		bigram = b[i : i+2]
		if !bSet[bigram] {
			totalBigrams++
			bSet[bigram] = true
			if aSet[bigram] {
				sharedBigrams++
			}
		}
	}
	return 2 * sharedBigrams / totalBigrams, nil
}

// WhiteSimilarity calculates the similarity of two
// strings through a variation on the Sorensen-Dice
// Coefficient algorithm, bytewise.
//
// The resulting value is scaled between 0 and 1.0,
// and a higher value means a higher similarity.
//
// WhiteSimilarity differs from DiceCoefficient in that
// it disregards bigrams that include (single-byte)
// whitespace, applies an upper-case filter, and
// accounts for bigram frequency.
//
// See: http://www.catalysoft.com/articles/strikeamatch.html
//
// Note that this algorithm implementation operates upon
// individual bytes and does not account for multibyte
// unicode runes.
//
// Returns an error if neither of the input strings
// contains at least one byte bigram without whitespace.
func WhiteSimilarity(a, b string) (float64, error) {
	aPairs, aLen := asciiUpperWordLetterPairs(a)
	bPairs, bLen := asciiUpperWordLetterPairs(b)
	union := aLen + bLen
	if union == 0 {
		return 0.0, errors.New("At least one of the input strings must contain two or more non-whitespace byte bigrams in order to calculate the White Similarity")
	}
	intersection := 0.0
	for _, aBigram := range aPairs {
		for j, bBigram := range bPairs {
			if aBigram == bBigram {
				intersection++
				bPairs[j] = byteBigram{}
				break
			}
		}
	}
	return 2 * intersection / float64(union), nil
}

func asciiUpperWordLetterPairs(s string) ([]byteBigram, int) {
	limit := len(s) - 1
	if limit < 1 {
		return make([]byteBigram, 0), 0
	}
	bigrams := make([]byteBigram, limit)
	var a byte
	var b byte
	var aIsSpace bool
	var bIsSpace bool
	numPairs := 0
	for i := 0; i < limit; i++ {
		b, bIsSpace = asciiUpperOrSpace(s[i+1])
		if bIsSpace {
			i++
			continue
		}
		a, aIsSpace = asciiUpperOrSpace(s[i])
		if aIsSpace {
			continue
		}
		bigrams[numPairs] = byteBigram{a: a, b: b}
		numPairs++
	}
	bigrams = bigrams[0:numPairs]
	return bigrams, numPairs
}

func asciiUpperOrSpace(b byte) (byte, bool) {
	if b <= unicode.MaxASCII {
		if 'a' <= b && b <= 'z' {
			return b - ('a' - 'A'), false
		}
		switch b {
		case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0:
			return b, true
		}
	}
	return b, false
}

// LevenshteinDistance calculates the magnitude of
// difference between two strings using the
// Levenshtein Distance metric, bytewise.
//
// This edit distance is the minimum number of single-byte
// edits (insertions, deletions, or substitutions) needed
// to transform one string into another.
//
// The larger the result, the more different the strings.
//
// See: http://en.wikipedia.org/wiki/Levenshtein_distance
func LevenshteinDistance(a, b string) (int, error) {
	aLen := len(a)
	bLen := len(b)
	if aLen == 0 {
		return bLen, nil
	}
	if bLen == 0 {
		return aLen, nil
	}
	if aLen == bLen && a == b {
		return 0, nil
	}
	rowLen := bLen + 1
	prevRow := make([]int, rowLen, rowLen)
	currRow := make([]int, rowLen, rowLen)
	for h := 0; h < rowLen; h++ {
		prevRow[h] = h
	}
	cost := 0
	for i := 0; i < aLen; i++ {
		currRow[0] = i + 1
		for j := 0; j < bLen; j++ {
			if a[i] == b[i] {
				cost = 0
			} else {
				cost = 1
			}
			currRow[j+1] = min(
				currRow[j]+1,
				prevRow[j+1]+1,
				prevRow[j]+cost)
		}
		prevRow, currRow = currRow, prevRow
	}
	return prevRow[bLen], nil
}

// DamerauLevenshteinDistance calculates the magnitude
// of difference between two strings using the Damerau-
// Levenshtein algorithm with adjacent-only transpositions,
// bytewise.
//
// This edit distance is the minimum number of single-byte
// edits (insertions, deletions, substitutions, or
// transpositions) to transform one string into the other.
// DamerauLevenshtein differs from Levenshtein primarily
// in that DamerauLevenshtein considers adjacent-byte transpositions.
//
// The larger the result, the more different the strings.
//
// See: http://en.wikipedia.org/wiki/Damerau-Levenshtein_distance
func DamerauLevenshteinDistance(a, b string) (int, error) {
	aLen := len(a)
	bLen := len(b)
	if aLen == 0 {
		return bLen, nil
	} else if bLen == 0 {
		return aLen, nil
	}

	// Swap to ensure a contains the shorter string
	if aLen > bLen {
		a, aLen, b, bLen = b, bLen, a, aLen
	}
	rowLen := aLen + 1
	tranRow := make([]int, rowLen, rowLen)
	prevRow := make([]int, rowLen, rowLen)
	currRow := make([]int, rowLen, rowLen)
	for h := 0; h < rowLen; h++ {
		prevRow[h] = h
	}
	var prevB byte
	var cost int
	for i := 1; i <= bLen; i++ {
		currB := b[i-1]
		currRow[0] = i

		start := i - bLen - 1
		if start < 1 {
			start = 1
		}
		end := i + bLen + 1
		if end > aLen {
			end = aLen
		}

		var prevA byte
		for j := start; j <= end; j++ {
			currA := a[j-1]
			if currA == currB {
				cost = 0
			} else {
				cost = 1
			}
			entry := min(
				currRow[j-1]+1,
				prevRow[j]+1,
				prevRow[j-1]+cost)
			if currA == prevB && currB == prevA {
				trans := tranRow[j-2] + cost
				if trans < entry {
					entry = trans
				}
			}
			currRow[j] = entry
			prevA = currA
		}
		prevB = currB
		tranRow, prevRow, currRow = prevRow, currRow, tranRow
	}
	return prevRow[aLen], nil
}

func min(a, b, c int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		return c
	}
	return m
}

type byteBigram struct {
	a, b byte
}
