/*
Copyright 2013 Zack Pierce.
Use of this source code is governed by a MIT-style
license that can be found in the LICENSE file.
*/
/*
Package stralgo/runewise implements various string algorithms
with an emphasis on similarity metrics, implemented with
support for multi-byte runes.
*/
package runewise

import (
	"errors"
	"fmt"
	"sort"
	"unicode"
)

const (
	WinklerBoostThreshold  = 0.7 // JaroWinklerSimilarity suggested parameter. If the JaroSimilarity for the compared strings is above this value, add an additional boost factor based on the shared prefix length and prefix scale.
	WinklerMaxPrefixLength = 4   // JaroWinklerSimilarity suggested parameter. Used to control the maximum size of identical prefixes used in the prefix boost factor.
	WinklerPrefixScale     = 0.1 // JaroWinklerSimilarity suggested parameter. Used to control the scale of bonus added for a pair having a JaroSimilarity above the threshold and with shared string prefixes.
)

// HammingDistance calculates the Hamming distance between
// two strings containing equal numbers of runes.
//
// The Hamming distance is the total number of runes
// that differ at the same index within the two strings.
//
// The higher the result, the more different the strings.

// See: http://en.wikipedia.org/wiki/Hamming_distance
//
// Returns an error if the string rune counts are not equal.
func HammingDistance(a, b []rune) (uint, error) {
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
// strings per the Sorensen-Dice coefficient, runewise.
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
// See: http://en.wikipedia.org/wiki/Dice_coefficient
//
// Returns an error if both of the input strings contain
// less than two runes.
func DiceCoefficient(a, b []rune) (float64, error) {
	aLen := len(a)
	bLen := len(b)
	aLimit := aLen - 1
	bLimit := bLen - 1
	if aLimit < 1 && bLimit < 1 {
		return 0, errors.New("At least one of the input strings must contain 2 or more runes for the bigram-based DiceCoefficient to be calculated.")
	}
	aSet := make(map[runeBigram]bool, aLimit)
	totalBigrams := 0.0
	var bigram runeBigram
	for i := 0; i < aLimit; i++ {
		bigram = runeBigram{rA: a[i], rB: a[i+1]}
		if !aSet[bigram] {
			totalBigrams++
			aSet[bigram] = true
		}
	}

	bSet := make(map[runeBigram]bool, bLimit)
	sharedBigrams := 0.0
	for i := 0; i < bLimit; i++ {
		bigram = runeBigram{rA: b[i], rB: b[i+1]}
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
// Coefficient algorithm.
//
// The resulting value is scaled between 0 and 1.0,
// and a higher value means a higher similarity.
//
// WhiteSimilarity differs from DiceCoefficient in that
// it disregards bigrams that include whitespace, and
// applies an upper-case filter, and accounts for bigram
// frequency.
//
// See: http://www.catalysoft.com/articles/strikeamatch.html
//
// Returns an error if neither of the input strings
// contains at least one rune bigram without whitespace.
func WhiteSimilarity(a, b []rune) (float64, error) {
	aPairs, aLen := upperWordLetterPairs(a)
	bPairs, bLen := upperWordLetterPairs(b)
	union := aLen + bLen
	if union == 0 {
		return 0.0, errors.New("At least one of the input strings must contain two or more non-whitespace rune bigrams in order to calculate the White Similarity.")
	}
	intersection := 0.0
	for _, aBigram := range aPairs {
		for j, bBigram := range bPairs {
			if aBigram == bBigram {
				intersection++
				bPairs[j] = runeBigram{}
				break
			}
		}
	}
	return 2 * intersection / float64(union), nil
}

func upperWordLetterPairs(runes []rune) ([]runeBigram, int) {
	limit := len(runes) - 1
	if limit < 1 {
		return make([]runeBigram, 0), 0
	}
	bigrams := make([]runeBigram, limit)
	var a rune
	var b rune
	numPairs := 0
	for i := 0; i < limit; i++ {
		a = runes[i]
		b = runes[i+1]
		if unicode.IsSpace(b) {
			i++
			continue
		}
		if unicode.IsSpace(a) {
			continue
		}
		bigrams[numPairs] = runeBigram{rA: unicode.ToUpper(a), rB: unicode.ToUpper(b)}
		numPairs++
	}
	bigrams = bigrams[0:numPairs]
	return bigrams, numPairs
}

// LevenshteinDistance calculates the magnitude of
// difference between two strings using the
// Levenshtein Distance metric.
//
// This edit distance is the minimum number of single-rune
// edits (insertions, deletions, or substitutions) needed
// to transform one string into the other.
//
// The larger the result, the more different the strings.
//
// See: http://en.wikipedia.org/wiki/Levenshtein_distance
func LevenshteinDistance(a, b []rune) (int, error) {
	aLen := len(a)
	bLen := len(b)
	if aLen == 0 {
		return bLen, nil
	}
	if bLen == 0 {
		return aLen, nil
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
			if a[i] == b[j] {
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
// runewise.
//
// This edit distance is the minimum number of single-rune
// edits (insertions, deletions, substitutions, or
// transpositions) to transform one string into the other.
// DamerauLevenshtein differs from Levenshtein primarily
// in that DamerauLevenshtein considers adjacent-rune transpositions.
//
// The larger the result, the more different the strings.
//
// See: http://en.wikipedia.org/wiki/Damerau-Levenshtein_distance
func DamerauLevenshteinDistance(a, b []rune) (int, error) {
	aLen := len(a)
	bLen := len(b)
	if aLen == 0 {
		return bLen, nil
	} else if bLen == 0 {
		return aLen, nil
	}

	// Swap to ensure a contains the shorter slice
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
	var prevB rune
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

		var prevA rune
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

// JaroSimilarity calculates the similarity between two strings
// using the original Jaro distance formula.
//
// The result is between 0 and 1.0, and the higher the score,
// the more similar the two strings are. 1.0 is a perfect match.
//
// If either input argument is empty ([]rune("")) or nil, the result
// will be 0.0. This is due to a quirk in the formal definition of
// the algorithm which counts the number of matching characters.
// In the empty or nil cases, no matches may be found at all.
//
// See (the first half of) : http://en.wikipedia.org/wiki/Jaro-Winkler_distance
//
// See also : http://alias-i.com/lingpipe/docs/api/com/aliasi/spell/JaroWinklerDistance.html
func JaroSimilarity(a, b []rune) float64 {
	matches, transpositions := jaroMatchesAndHalfTranspositions(a, b)

	if matches == 0 {
		return 0.0
	}

	matchFloat := float64(matches)
	return (1.0 / 3.0) * (matchFloat/float64(len(a)) + matchFloat/float64(len(b)) + (matchFloat-float64(transpositions/2))/matchFloat)
}

// jaroMatchesAndHalfTranspositions calculates the number of
// matches and half-transpositions defined by the Jaro distance
// formula.
func jaroMatchesAndHalfTranspositions(a, b []rune) (int, int) {
	aLen := len(a)
	bLen := len(b)
	if aLen == 0 || bLen == 0 {
		return 0, 0
	}
	if aLen < bLen {
		a, aLen, b, bLen = b, bLen, a, aLen
	}
	matchMax := (aLen / 2) - 1
	if matchMax < 0 {
		matchMax = 0
	}
	aCommon := make([]rune, aLen, aLen)
	numAMatched := 0
	bMatchedIndices := make(map[int]bool, aLen)
	for i, aRune := range a {
		from := i - matchMax
		if from < 0 {
			from = 0
		}
		to := i + matchMax
		if to >= bLen {
			to = bLen - 1
		}
		aMatched := false
		for j := from; j <= to; j++ {
			if aRune != b[j] {
				continue
			}
			if !aMatched {
				aCommon[numAMatched] = aRune
				aMatched = true
				numAMatched++
			}
			if _, ok := bMatchedIndices[j]; !ok {
				bMatchedIndices[j] = true
			}
		}
	}

	bIndices := make([]int, numAMatched, numAMatched)
	c := 0
	for s, _ := range bMatchedIndices {
		bIndices[c] = s
		c++
	}
	sort.Ints(bIndices)

	transCount := 0
	for k := 0; k < numAMatched; k++ {
		if aCommon[k] != b[bIndices[k]] {
			transCount++
		}
	}
	return numAMatched, transCount
}

// JaroWinklerSimilarity calculates the similarity between
// two input strings using the Jaro-Winkler distance formula.
//
// Winkler's suggested constants for max considered common prefix
// length (4), common prefix scaling factor (0.1), and boost
// threshold (0.7) are used.
//
// The result is between 0 and 1.0, and the higher the score,
// the more similar the two strings are. 1.0 is a perfect match.
//
// If either input argument is empty ([]rune("")) or nil, the result
// will be 0.0. This is due to a quirk in the formal definition of
// the algorithm which counts the number of matching characters.
// In the empty or nil cases, no matches may be found at all.
//
// See : http://en.wikipedia.org/wiki/Jaro-Winkler_distance
//
// See : http://alias-i.com/lingpipe/docs/api/com/aliasi/spell/JaroWinklerDistance.html
//
// Note that the wikipedia article does not include a description
// of Winkler's boost threshold, an explanation of which can be
// found in the lingpipe documentation (linked above), and is
// demonstrated in Winkler's original code.
//
// In short, the boost threshold has the following effect:
//
//    if calculatedJaroSimilarity < WinklerBoostThreshold {
//    	return calculatedJaroSimilarity
//    } else {
//    	return calculatedJaroSimilarity + prefixSimilarityBonus
//    }
//
// The prefixSimilarityBonus is the modification to the original
// Jaro formula described on the Wikipedia article, and is equivalent to:
//
//    Min(calculatedLengthOfCommonPrefix, WinklerMaxPrefixLength)*WinklerPrefixScale*(1 - calculatedJaroSimilarity)
func JaroWinklerSimilarity(a, b []rune) float64 {
	return JaroWinklerSimilarityParametric(a, b, WinklerPrefixScale, WinklerMaxPrefixLength, WinklerBoostThreshold)
}

// JaroWinklerSimilarityParametric calculates similarity between
// two input strings using the Jaro-Winkler distance formula.
//
// The product of prefixScale and maxPrefixLength should be between 0.0 and 1.0.
// Assuming this is true, the result will be between 0 and 1.0.
//
// The higher the score, the more similar the two strings are.
// 1.0 is a perfect match.
//
// See : http://en.wikipedia.org/wiki/Jaro-Winkler_distance
//
// See : http://alias-i.com/lingpipe/docs/api/com/aliasi/spell/JaroWinklerDistance.html
//
// Note that the wikipedia article does not include a description
// of Winkler's boost threshold, an explanation of which can be
// found in the lingpipe documentation (linked above), and is
// demonstrated in Winkler's original code.
//
// In short, the boost threshold has the following effect:
//
//    if calculatedJaroSimilarity < boostThreshold {
//    	return calculatedJaroSimilarity
//    } else {
//    	return calculatedJaroSimilarity + prefixSimilarityBonus
//    }
//
// The prefixSimilarityBonus is the modification to the original
// Jaro formula described on the Wikipedia article, and is equivalent to:
//
//    Min(calculatedLengthOfCommonPrefix, maxPrefixLength)*prefixScale*(1 - calculatedJaroSimilarity)/
func JaroWinklerSimilarityParametric(a, b []rune, prefixScale float64, maxPrefixLength int, boostThreshold float64) float64 {
	j := JaroSimilarity(a, b)
	if j < boostThreshold {
		return j
	}
	return j + float64(clampedSharedPrefixLength(a, b, maxPrefixLength))*prefixScale*(1.0-j)
}

func clampedSharedPrefixLength(a, b []rune, maxPrefixLength int) int {
	minLen := min(len(a), len(b), maxPrefixLength)
	i := 0
	for ; i < minLen; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return i
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

type runeBigram struct {
	rA, rB rune
}

func (r runeBigram) String() string {
	return fmt.Sprintf("{%q, %q}", r.rA, r.rB)
}
