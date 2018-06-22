package runewise

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func Test_HammingDistance(t *testing.T) {
	d, err := HammingDistance([]rune("toned"), []rune("roses"))
	assert.Nil(t, err, "HammingDistance should not return an error for two strings of equivalent length.")
	assert.Equal(t, uint(3), d, "HammingDistance between 'toned' and 'roses' should be 3")

	d, err = HammingDistance([]rune(""), []rune(""))
	assert.Nil(t, err, "HammingDistance should not return an error for two empty strings")
	assert.Equal(t, uint(0), d, "HammingDistance between two empty strings should be 0")

	d, err = HammingDistance([]rune("spam"), []rune("spam"))
	assert.Nil(t, err, "HammingDistance should not return an error for two equivalent non-empty strings.")
	assert.Equal(t, uint(0), d, "HammingDistance between 'spam' and 'spam' should be 0, but was %i", d)

	d, err = HammingDistance([]rune("green eggs"), []rune("ham"))
	assert.NotNil(t, err, "HammingDistance between 'green eggs' and 'ham' should produce an error due to unequal lengths")
	assert.Equal(t, uint(0), d)

	d, err = HammingDistance([]rune("日本語"), []rune("日本ゴ"))
	assert.Nil(t, err)
	assert.Equal(t, uint(1), d)

	d, err = HammingDistance([]rune("日本語"), []rune("日本g"))
	assert.Nil(t, err)
	assert.Equal(t, uint(1), d)

	d, err = HammingDistance([]rune("日本語"), []rune("日本go"))
	assert.NotNil(t, err)
	assert.Equal(t, uint(0), d)
}

func Test_DiceCoefficient(t *testing.T) {
	c, err := DiceCoefficient([]rune("night"), []rune("nacht"))
	assert.Nil(t, err)
	assert.Equal(t, 1.0/4.0, c)

	c, err = DiceCoefficient([]rune("GGGG"), []rune("GGGG"))
	assert.Nil(t, err)
	assert.Equal(t, 1.0, c)

	c, err = DiceCoefficient([]rune(""), []rune(""))
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, c)

	c, err = DiceCoefficient([]rune("a"), []rune("b"))
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, c)

	c, err = DiceCoefficient([]rune("GG"), []rune("GGGG"))
	assert.Nil(t, err)
	assert.Equal(t, 1.0, c, "Naive Dice coefficient does not account for differences in occurrence-count for bigrams.")

	c, err = DiceCoefficient([]rune("日本語"), []rune("日本語"))
	assert.Nil(t, err)
	assert.Equal(t, 1.0, c)

	c, err = DiceCoefficient([]rune("日本語"), []rune("日本ゴ"))
	assert.Nil(t, err)
	assert.Equal(t, 2.0/4.0, c)

	c, err = DiceCoefficient([]rune("日本語"), []rune("日本g"))
	assert.Nil(t, err)
	assert.Equal(t, 2.0/4.0, c)

	c, err = DiceCoefficient([]rune("日"), []rune("本"))
	assert.NotNil(t, err, "DiceCoefficient runewise reports an error due to lack of rune-bigrams.")
	assert.Equal(t, 0.0, c)

	// [230 151 165] and [230, 151, 168]
	c, err = DiceCoefficient([]rune("日"), []rune("旨"))
	assert.NotNil(t, err, "DiceCoefficient runewise reports an error due to lack of rune-bigrams.")
	assert.Equal(t, 0.0, c)
}

func Test_WhiteSimilarity(t *testing.T) {
	c, err := WhiteSimilarity([]rune("Healed"), []rune("Healed"))
	assert.Nil(t, err)
	assert.Equal(t, 1.0, c)

	c, err = WhiteSimilarity([]rune("Healed"), []rune("Sealed"))
	assert.Nil(t, err)
	EqualWithin(t, 0.8, c, 0.01, "Sealed")

	c, err = WhiteSimilarity([]rune("Healed"), []rune("Healthy"))
	assert.Nil(t, err)
	EqualWithin(t, 0.55, c, 0.01, "Healthy")

	c, err = WhiteSimilarity([]rune("Healed"), []rune("Heard"))
	assert.Nil(t, err)
	EqualWithin(t, 0.44, c, 0.01, "Heard")

	c, err = WhiteSimilarity([]rune("Healed"), []rune("Herded"))
	assert.Nil(t, err)
	EqualWithin(t, 0.40, c, 0.01, "Herded")

	c, err = WhiteSimilarity([]rune("Healed"), []rune("Help"))
	assert.Nil(t, err)
	EqualWithin(t, 0.25, c, 0.01, "Help")

	c, err = WhiteSimilarity([]rune("Healed"), []rune("Sold"))
	assert.Nil(t, err)
	EqualWithin(t, 0.0, c, 0.01, "Sold")

	c, err = WhiteSimilarity([]rune("Healed "), []rune("HEALed"))
	assert.Nil(t, err)
	assert.Equal(t, 1.0, c)

	c, err = WhiteSimilarity([]rune("GGGG"), []rune("GGGG"))
	assert.Nil(t, err)
	assert.Equal(t, 1.0, c)

	c, err = WhiteSimilarity([]rune("REPUBLIC OF FRANCE"), []rune("FRANCE"))
	assert.Nil(t, err)
	EqualWithin(t, 0.56, c, 0.01)

	c, err = WhiteSimilarity([]rune("FRANCE"), []rune("QUEBEC"))
	assert.Nil(t, err)
	EqualWithin(t, 0.0, c, 0.001)

	c, err = WhiteSimilarity([]rune("FRENCH REPUBLIC"), []rune("REPUBLIC OF FRANCE"))
	assert.Nil(t, err)
	EqualWithin(t, 0.72, c, 0.01)

	c, err = WhiteSimilarity([]rune("FRENCH REPUBLIC"), []rune("REPUBLIC OF CUBA"))
	assert.Nil(t, err)
	EqualWithin(t, 0.61, c, 0.01)

	c, err = WhiteSimilarity([]rune(""), []rune(""))
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, c)

	c, err = WhiteSimilarity([]rune("a"), []rune("b"))
	assert.NotNil(t, err)
	assert.Equal(t, 0.0, c)

	c, err = WhiteSimilarity([]rune("GG"), []rune("GGGGG"))
	assert.Nil(t, err)
	EqualWithin(t, 0.4, c, 0.01)

	c, err = WhiteSimilarity([]rune("GGGGG"), []rune("GG"))
	assert.Nil(t, err)
	EqualWithin(t, 0.4, c, 0.01)
}

func Test_LevenshteinDistance_Easy(t *testing.T) {
	d, err := LevenshteinDistance([]rune("test"), []rune("test"))
	assert.Nil(t, err)
	assert.Equal(t, 0, d)

	d, err = LevenshteinDistance([]rune("test"), []rune("tent"))
	assert.Nil(t, err)
	assert.Equal(t, 1, d)

	d, err = LevenshteinDistance([]rune("gumbo"), []rune("gambol"))
	assert.Nil(t, err)
	assert.Equal(t, 2, d)

	d, err = LevenshteinDistance([]rune("kitten"), []rune("sitting"))
	assert.Nil(t, err)
	assert.Equal(t, 3, d)

	d, err = LevenshteinDistance([]rune("foo"), []rune(""))
	assert.Nil(t, err)
	assert.Equal(t, 3, d)

	d, err = LevenshteinDistance([]rune(""), []rune("foo"))
	assert.Nil(t, err)
	assert.Equal(t, 3, d)

	d, err = LevenshteinDistance([]rune(""), []rune(""))
	assert.Nil(t, err)
	assert.Equal(t, 0, d)

	d, err = LevenshteinDistance([]rune("a"), []rune(""))
	assert.Nil(t, err)
	assert.Equal(t, 1, d)

}

func Test_DamerauLevenshteinDistance(t *testing.T) {
	d, err := DamerauLevenshteinDistance([]rune("azertyuiop"), []rune("aeryuop"))
	assert.Nil(t, err)
	assert.Equal(t, 3, d)

	d, err = DamerauLevenshteinDistance([]rune("aeryuop"), []rune("azertyuiop"))
	assert.Nil(t, err)
	assert.Equal(t, 3, d)

	d, err = DamerauLevenshteinDistance([]rune("azertyuiopqsdfghjklmwxcvbn"), []rune("qwertyuiopasdfghjkl;zxcvbnm"))
	assert.Nil(t, err)
	assert.Equal(t, 6, d)

	d, err = DamerauLevenshteinDistance([]rune("1234567890"), []rune("1324576809"))
	assert.Nil(t, err)
	assert.Equal(t, 3, d)

	d, err = DamerauLevenshteinDistance([]rune("ab"), []rune("ab"))
	assert.Nil(t, err)
	assert.Equal(t, 0, d)

	d, err = DamerauLevenshteinDistance([]rune(""), []rune("ab"))
	assert.Nil(t, err)
	assert.Equal(t, 2, d)

	d, err = DamerauLevenshteinDistance([]rune("ab"), []rune(""))
	assert.Nil(t, err)
	assert.Equal(t, 2, d)

	d, err = DamerauLevenshteinDistance([]rune("Cedarinia scabra Sjöstedt 1921"), []rune("Cedarinia scabra Söjstedt 1921"))
	assert.Nil(t, err)
	assert.Equal(t, 1, d)
}

func Test_Jaro_Empty(t *testing.T) {
	c := JaroSimilarity([]rune(""), []rune(""))
	assert.Equal(t, 0.0, c, "Empty strings should produce 0.0 for Jaro")
	c = JaroSimilarity([]rune(""), []rune("a"))
	assert.Equal(t, 0.0, c, "First empty string should produce 0.0 for Jaro")
	c = JaroSimilarity([]rune("b"), []rune(""))
	assert.Equal(t, 0.0, c, "Second empty string should produce 0.0 for Jaro")
	c = JaroSimilarity(nil, nil)
	assert.Equal(t, 0.0, c, "Nil strings should produce 0.0 for Jaro")
	c = JaroSimilarity(nil, []rune("a"))
	assert.Equal(t, 0.0, c, "First nil should produce 0.0 for Jaro")
	c = JaroSimilarity([]rune("a"), nil)
	assert.Equal(t, 0.0, c, "Second nil should produce 0.0 for Jaro")
}

func Test_Jaro_SimpleEquality(t *testing.T) {
	c := JaroSimilarity([]rune("a"), []rune("a"))
	assert.Equal(t, 1.0, c, "Equal strings should produce 1.0 for Jaro")

	c = JaroSimilarity([]rune("abc"), []rune("abc"))
	assert.Equal(t, 1.0, c, "Equal strings should produce 1.0 for Jaro")

	c = JaroSimilarity([]rune("abc"), []rune("123"))
	assert.Equal(t, 0.0, c, "Completely different strings should produce 0.0 for Jaro")
}

func Test_Jaro_Unequal(t *testing.T) {
	c := JaroSimilarity([]rune("abcvwxyz"), []rune("cabvwxyz"))
	EqualWithin(t, 0.958, c, 0.001)

	c = JaroSimilarity([]rune("abcduvwxyz"), []rune("dabcuvwxyz"))
	EqualWithin(t, (1.0/3.0)*(2.0+(10.0-2.0)/10.0), c, 0.0001)

	c = JaroSimilarity([]rune("abcduvwxyz"), []rune("dbacuvwxyz"))
	EqualWithin(t, (1.0/3.0)*(2.0+(10.0-1.0)/10.0), c, 0.0001)

	c = JaroSimilarity([]rune("martha"), []rune("marhta"))
	EqualWithin(t, 0.9444444, c, 0.0001, "martha and marhta")

	c = JaroSimilarity([]rune("dwayne"), []rune("duane"))
	EqualWithin(t, 0.8222222, c, 0.0001, "dwayne and duane")

	c = JaroSimilarity([]rune("dixon"), []rune("dicksonx"))
	EqualWithin(t, 0.7666666, c, 0.0001)

	c = JaroSimilarity([]rune("abcd"), []rune("qrsd"))
	EqualWithin(t, (1.0/3.0)*(1.0/4.0+1.0/4.0+1.0/1.0), c, 0.0001)
}

func Test_JaroWinkler_Empty(t *testing.T) {
	c := JaroWinklerSimilarity([]rune(""), []rune(""))
	assert.Equal(t, 0.0, c, "Empty strings should produce 0.0 for JaroWinkler")
	c = JaroWinklerSimilarity([]rune(""), []rune("a"))
	assert.Equal(t, 0.0, c, "First empty string should produce 0.0 for JaroWinkler")
	c = JaroWinklerSimilarity([]rune("b"), []rune(""))
	assert.Equal(t, 0.0, c, "Second empty string should produce 0.0 for JaroWinkler")
	c = JaroWinklerSimilarity(nil, nil)
	assert.Equal(t, 0.0, c, "Nil strings should produce 0.0 for JaroWinkler")
	c = JaroWinklerSimilarity(nil, []rune("a"))
	assert.Equal(t, 0.0, c, "First nil should produce 0.0 for JaroWinkler")
	c = JaroWinklerSimilarity([]rune("a"), nil)
	assert.Equal(t, 0.0, c, "Second nil should produce 0.0 for JaroWinkler")
}

func Test_JaroWinkler_SimpleEquality(t *testing.T) {
	c := JaroWinklerSimilarity([]rune("a"), []rune("a"))
	assert.Equal(t, 1.0, c, "Equal strings should produce 1.0 for JaroWinkler")

	c = JaroWinklerSimilarity([]rune("abc"), []rune("abc"))
	assert.Equal(t, 1.0, c, "Equal strings should produce 1.0 for JaroWinkler")

	c = JaroWinklerSimilarity([]rune("abc"), []rune("123"))
	assert.Equal(t, 0.0, c, "Completely different strings should produce 0.0 for JaroWinkler")
}

func Test_JaroWinkler_Unequal(t *testing.T) {
	c := JaroWinklerSimilarity([]rune("abcduvwxyz"), []rune("dabcuvwxyz"))
	EqualWithin(t, (1.0/3.0)*(2.0+(10.0-2.0)/10.0)+0.0, c, 0.0001)

	c = JaroWinklerSimilarity([]rune("abcduvwxyz"), []rune("dbacuvwxyz"))
	EqualWithin(t, (1.0/3.0)*(2.0+(10.0-1.0)/10.0)+0.0, c, 0.0001)

	c = JaroWinklerSimilarity([]rune("martha"), []rune("marhta"))
	EqualWithin(t, 0.9444444+(0.1*3*(1-0.944444)), c, 0.0001, "martha and marhta")

	c = JaroWinklerSimilarity([]rune("dwayne"), []rune("duane"))
	EqualWithin(t, 0.8222222+(0.1*1*(1-0.822222)), c, 0.0001, "dwayne and duane")

	c = JaroWinklerSimilarity([]rune("dixon"), []rune("dicksonx"))
	EqualWithin(t, 0.7666666+(0.1*2*(1-0.7666666)), c, 0.0001)

	c = JaroWinklerSimilarity([]rune("abcd"), []rune("qrsd"))
	EqualWithin(t, (1.0/3.0)*(1.0/4.0+1.0/4.0+1.0/1.0)+0.0, c, 0.0001)
}

func Benchmark_LevenshteinDistance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LevenshteinDistance([]rune("kitten"), []rune("sitting"))
		LevenshteinDistance([]rune("gumbo"), []rune("gambol"))
	}
}

func EqualWithin(t *testing.T, a, b, delta float64, msgAndArgs ...interface{}) bool {
	if math.Abs(a-b) > delta {
		return assert.Fail(t, fmt.Sprintf("Not within delta: Abs(%#v - %#v) > %#v", a, b, delta), msgAndArgs...)
	}

	return true
}
