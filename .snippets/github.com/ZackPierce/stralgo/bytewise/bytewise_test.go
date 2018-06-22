package bytewise

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func Test_HammingDistance(t *testing.T) {
	d, err := HammingDistance("toned", "roses")
	assert.Nil(t, err, "HammingDistance should not return an error for two strings of equivalent length.")
	assert.Equal(t, uint(3), d, "HammingDistance between 'toned' and 'roses' should be 3")

	d, err = HammingDistance("", "")
	assert.Nil(t, err, "HammingDistance should not return an error for two empty strings")
	assert.Equal(t, uint(0), d, "HammingDistance between two empty strings should be 0")

	d, err = HammingDistance("spam", "spam")
	assert.Nil(t, err, "HammingDistance should not return an error for two equivalent non-empty strings.")
	assert.Equal(t, uint(0), d, "HammingDistance between 'spam' and 'spam' should be 0, but was %i", d)

	d, err = HammingDistance("green eggs", "ham")
	assert.NotNil(t, err, "HammingDistance between 'green eggs' and 'ham' should produce an error due to unequal lengths")

	d, err = HammingDistance("日本語", "日本ゴ")
	assert.Nil(t, err)
	assert.Equal(t, uint(3), d)

	d, err = HammingDistance("日本語", "日本g")
	assert.NotNil(t, err, "HammingDistance bytewise produces an error due to the difference in total bytes in the two strings, even though the rune-counts are equivalent.")
	assert.Equal(t, uint(0), d)

	d, err = HammingDistance("日本語", "日本gon")
	assert.Nil(t, err, "HammingDistance bytewise does not produce an error when comparing strings of equal byte-lengths, even though the rune-counts are different.")
	assert.Equal(t, uint(3), d)
}

func Test_DiceCoefficient(t *testing.T) {
	c, err := DiceCoefficient("night", "nacht")
	assert.Nil(t, err)
	assert.Equal(t, 1.0/4.0, c)

	c, err = DiceCoefficient("GGGG", "GGGG")
	assert.Nil(t, err)
	assert.Equal(t, 1.0, c)

	c, err = DiceCoefficient("", "")
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, c)

	c, err = DiceCoefficient("a", "b")
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, c)

	c, err = DiceCoefficient("GG", "GGGG")
	assert.Nil(t, err)
	assert.Equal(t, 1.0, c, "Naive Dice coefficient does not account for differences in occurrence-count for bigrams.")

	c, err = DiceCoefficient("日", "本")
	assert.Nil(t, err, "DiceCoefficient bytewise does not report an error about lack of bigrams for this case, because the runes involved have a width of 2 or greater.")
	assert.Equal(t, 0.0, c)

	// [230 151 165] and [230, 151, 168]
	c, err = DiceCoefficient("日", "旨")
	assert.Nil(t, err, "DiceCoefficient bytewise does not report an error about lack of bigrams for this case, because the runes involved have a width of 2 or greater.")
	assert.Equal(t, 2.0/4.0, c)
}

func Test_WhiteSimilarity(t *testing.T) {
	c, err := WhiteSimilarity("Healed", "Healed")
	assert.Nil(t, err)
	assert.Equal(t, 1.0, c)

	c, err = WhiteSimilarity("Healed", "Sealed")
	assert.Nil(t, err)
	EqualWithin(t, 0.8, c, 0.01, "Sealed")

	c, err = WhiteSimilarity("Healed", "Healthy")
	assert.Nil(t, err)
	EqualWithin(t, 0.55, c, 0.01, "Healthy")

	c, err = WhiteSimilarity("Healed", "Heard")
	assert.Nil(t, err)
	EqualWithin(t, 0.44, c, 0.01, "Heard")

	c, err = WhiteSimilarity("Healed", "Herded")
	assert.Nil(t, err)
	EqualWithin(t, 0.40, c, 0.01, "Herded")

	c, err = WhiteSimilarity("Healed", "Help")
	assert.Nil(t, err)
	EqualWithin(t, 0.25, c, 0.01, "Help")

	c, err = WhiteSimilarity("Healed", "Sold")
	assert.Nil(t, err)
	EqualWithin(t, 0.0, c, 0.01, "Sold")

	c, err = WhiteSimilarity("Healed ", "HEALed")
	assert.Nil(t, err)
	assert.Equal(t, 1.0, c)

	c, err = WhiteSimilarity("GGGG", "GGGG")
	assert.Nil(t, err)
	assert.Equal(t, 1.0, c)

	c, err = WhiteSimilarity("REPUBLIC OF FRANCE", "FRANCE")
	assert.Nil(t, err)
	EqualWithin(t, 0.56, c, 0.01)

	c, err = WhiteSimilarity("FRANCE", "QUEBEC")
	assert.Nil(t, err)
	EqualWithin(t, 0.0, c, 0.001)

	c, err = WhiteSimilarity("FRENCH REPUBLIC", "REPUBLIC OF FRANCE")
	assert.Nil(t, err)
	EqualWithin(t, 0.72, c, 0.01)

	c, err = WhiteSimilarity("FRENCH REPUBLIC", "REPUBLIC OF CUBA")
	assert.Nil(t, err)
	EqualWithin(t, 0.61, c, 0.01)

	c, err = WhiteSimilarity("", "")
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, c)

	c, err = WhiteSimilarity("a", "b")
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, c)

	c, err = WhiteSimilarity("GG", "GGGGG")
	assert.Nil(t, err)
	EqualWithin(t, 0.4, c, 0.01)

	c, err = WhiteSimilarity("GGGGG", "GG")
	assert.Nil(t, err)
	EqualWithin(t, 0.4, c, 0.01)
}

func Test_LevenshteinDistance_Easy(t *testing.T) {
	d, err := LevenshteinDistance("test", "test")
	assert.Nil(t, err)
	assert.Equal(t, 0, d)

	d, err = LevenshteinDistance("test", "tent")
	assert.Nil(t, err)
	assert.Equal(t, 1, d)

	d, err = LevenshteinDistance("gumbo", "gambol")
	assert.Nil(t, err)
	assert.Equal(t, 2, d)

	d, err = LevenshteinDistance("kitten", "sitting")
	assert.Nil(t, err)
	assert.Equal(t, 3, d)

	d, err = LevenshteinDistance("foo", "")
	assert.Nil(t, err)
	assert.Equal(t, 3, d)

	d, err = LevenshteinDistance("", "foo")
	assert.Nil(t, err)
	assert.Equal(t, 3, d)

	d, err = LevenshteinDistance("", "")
	assert.Nil(t, err)
	assert.Equal(t, 0, d)

	d, err = LevenshteinDistance("a", "")
	assert.Nil(t, err)
	assert.Equal(t, 1, d)

}

func Test_DamerauLevenshteinDistance(t *testing.T) {
	d, err := DamerauLevenshteinDistance("azertyuiop", "aeryuop")
	assert.Nil(t, err)
	assert.Equal(t, 3, d)

	d, err = DamerauLevenshteinDistance("aeryuop", "azertyuiop")
	assert.Nil(t, err)
	assert.Equal(t, 3, d)

	d, err = DamerauLevenshteinDistance("azertyuiopqsdfghjklmwxcvbn", "qwertyuiopasdfghjkl;zxcvbnm")
	assert.Nil(t, err)
	assert.Equal(t, 6, d)

	d, err = DamerauLevenshteinDistance("1234567890", "1324576809")
	assert.Nil(t, err)
	assert.Equal(t, 3, d)

	d, err = DamerauLevenshteinDistance("ab", "ab")
	assert.Nil(t, err)
	assert.Equal(t, 0, d)

	d, err = DamerauLevenshteinDistance("", "ab")
	assert.Nil(t, err)
	assert.Equal(t, 2, d)

	d, err = DamerauLevenshteinDistance("ab", "")
	assert.Nil(t, err)
	assert.Equal(t, 2, d)

	d, err = DamerauLevenshteinDistance("Cedarinia scabra Sjöstedt 1921", "Cedarinia scabra Söjstedt 1921")
	assert.Nil(t, err)
	assert.Equal(t, 2, d, "Note that this requires two edits, despite the fact that only two adjacent runes have been transposed, due to the byte-wise handling approach")
}

func Benchmark_LevenshteinDistance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LevenshteinDistance("kitten", "sitting")
		LevenshteinDistance("gumbo", "gambol")
	}
}

func EqualWithin(t *testing.T, a, b, delta float64, msgAndArgs ...interface{}) bool {
	if math.Abs(a-b) > delta {
		return assert.Fail(t, fmt.Sprintf("Not within delta: Abs(%#v - %#v) > %#v", a, b, delta), msgAndArgs...)
	}

	return true
}
