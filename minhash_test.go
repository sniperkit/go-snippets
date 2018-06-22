package textsim

import (
	"bufio"
	"fmt"
	"hash"
	"hash/fnv"
	"math/rand"
	"os"
	"testing"

	"github.com/aduston/rabin"
	spooky "github.com/dgryski/go-spooky"
)

func tokenize(s string, wordHash hash.Hash64) uint64 {
	wordHash.Reset()
	for _, r := range s {
		wordHash.Write(toBytes32(uint32(r)))
	}
	return wordHash.Sum64()
}

func TestTokenize(t *testing.T) {
	text := "well, 日本語 is an àwesome language!"
	wordHash := fnv.New64a()
	tokens := Tokenize(text, wordHash)
	if len(tokens) != 6 {
		t.Errorf("Expected 6 tokens but got %d", len(tokens))
	}
	wordTokens := []string{"well", "日本語", "is", "an", "àwesome", "language"}
	for i, wordToken := range wordTokens {
		if tokens[i] != tokenize(wordToken, wordHash) {
			t.Errorf("Expected word %d to be correct hash for %v", i, wordToken)
		}
	}
}

func testShingles(rollingHash RollingHash, t *testing.T) {
	text := "the brown fox jumped over the brown fox jumped over the brown fox"
	tokens := Tokenize(text, fnv.New64a())
	shingles := ConvertToShingles(tokens, rollingHash)
	const numShingles = 10
	if len(shingles) != numShingles {
		t.Errorf("Expected %d shingles, got %d", numShingles, len(shingles))
	}
	for i := 0; i < len(shingles); i++ {
		for j := i + 1; j < len(shingles); j++ {
			if j != i+5 && shingles[i] == shingles[j] {
				t.Errorf("Did not expect shingles %d and %d to equal each other", i, j)
			} else if j == i+5 && shingles[i] != shingles[j] {
				t.Errorf("Expected shingles %d and %d to equal each other", i, j)
			}
		}
	}
}

func TestShinglesReg(t *testing.T) {
	testShingles(NewRegHashRollingHash(fnv.New64a(), 4), t)
}

func TestShinglesRabin(t *testing.T) {
	testShingles(NewRabinRollingHash(rabin.NewRolling(4*8), 4), t)
}

func readText(fileNo int) ([]string, error) {
	f, err := os.Open(fmt.Sprintf("testdata/text%d.txt", fileNo))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func benchmarkTokenize(b *testing.B, wordHash hash.Hash64) {
	lines, err := readText(0)
	if err != nil {
		b.Errorf("Couldn't read file: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, line := range lines {
			wordHash.Reset()
			Tokenize(line, wordHash)
		}
	}
}

func benchmarkConvertToShingles(b *testing.B, rollingHash RollingHash) {
	lines, err := readText(0)
	if err != nil {
		b.Errorf("Couldn't read file: %v", err)
	}
	tokens := make([][]uint64, len(lines))
	for i, line := range lines {
		tokens[i] = Tokenize(line, fnv.New64a())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rollingHash.Reset()
		for _, lineTokens := range tokens {
			ConvertToShingles(lineTokens, rollingHash)
		}
	}
}

func BenchmarkTokenizeRabin(b *testing.B) {
	benchmarkTokenize(b, rabin.New64())
}

func BenchmarkTokenizeFnv(b *testing.B) {
	benchmarkTokenize(b, fnv.New64a())
}

func BenchmarkTokenizeSpooky(b *testing.B) {
	benchmarkTokenize(b, spooky.New(uint64(rand.Int63()), uint64(rand.Int63())))
}

func BenchmarkConvertToShinglesRabin(b *testing.B) {
	benchmarkConvertToShingles(b, NewRabinRollingHash(rabin.NewRolling(4*8), 4))
}

func BenchmarkConvertToShinglesFnv(b *testing.B) {
	benchmarkConvertToShingles(b, NewRegHashRollingHash(fnv.New64a(), 4))
}

func makeShingles(b *testing.B, fileNo int) [][]uint64 {
	rollingHash := NewRegHashRollingHash(fnv.New64a(), 4)
	lines, err := readText(0)
	if err != nil {
		b.Errorf("Couldn't read file: %v", err)
	}
	tokens := make([][]uint64, len(lines))
	for i, line := range lines {
		tokens[i] = Tokenize(line, fnv.New64a())
	}
	shingles := make([][]uint64, len(lines))
	for i, lineTokens := range tokens {
		shingles[i] = ConvertToShingles(lineTokens, rollingHash)
	}
	return shingles
}

func BenchmarkPermutationFnv(b *testing.B) {
	shingles := makeShingles(b, 0)
	h1, h2 := MakePermHashes(fnv.New64a(), fnv.New64a())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, lineShingles := range shingles {
			CalcMinHashesDGryski(lineShingles, h1, h2, 1)
		}
	}
}

func BenchmarkPermutationLinear(b *testing.B) {
	shingles := makeShingles(b, 0)
	hashFuncs := GenerateLinearMinHashParms(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, lineShingles := range shingles {
			CalcMinHashesLinear(lineShingles, hashFuncs)
		}
	}
}
