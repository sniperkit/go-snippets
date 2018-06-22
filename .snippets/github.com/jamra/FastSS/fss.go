package fss

import (
	"hash"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/OneOfOne/xxhash/native"
)

type Fss struct {
	count int
	hf    hash.Hash32
	ht    map[uint32][]int
	d     int
	fi    map[int]string
}

func NewFss(d int) *Fss {
	f := &Fss{
		hf:    xxhash.New32(),
		ht:    make(map[uint32][]int),
		d:     d,
		count: 0,
		fi:    make(map[int]string),
	}
	return f
}

func (f *Fss) Insert(w string) {
	w = strings.ToLower(w)
	f.count++
	id := f.count
	f.fi[id] = w
	permutations := getpermutations(w, f.d)
	for _, p := range permutations {
		hash := f.gethash(p)
		f.ht[hash] = append(f.ht[hash], id)
	}
}

type QueryResult struct {
	S    string
	Rank float32
}

func (f *Fss) Search(q string) []QueryResult {
	q = strings.ToLower(q)
	results := make([]QueryResult, 0)

	permutations := getpermutations(q, f.d)

	for _, p := range permutations {
		hash := f.gethash(p)
		vals, ok := f.ht[hash]
		if ok {
			for _, val := range vals {
				str, _ := f.fi[val]
				qr := QueryResult{
					S: str,
				}
				results = append(results, qr)
			}
		}
	}
	return results
}

func (f *Fss) gethash(w string) uint32 {
	f.hf.Reset()
	r := strings.NewReader(w)
	io.Copy(f.hf, r)
	hash := f.hf.Sum32()
	return hash
}

type loc struct {
	index int
	width int
}

//getpermutations takes a word and an edit distance and returns all strings in the deletion neighborhood
//specified by the edit distance
func getpermutations(word string, del int) []string {
	locs := make([]loc, 0)

	for i, w := 0, 0; i < len(word); i += w {
		_, width := utf8.DecodeRuneInString(word[i:])
		locs = append(locs, loc{
			index: i,
			width: width,
		})
		w = width
	}

	perms := permutations(locs, del)
	results := stringsFromPermutations(word, locs, perms)

	return results
}

func permutations(locs []loc, del int) [][]loc {
	results := make([][]loc, 0)

	for d := 1; d <= del; d++ {
		results = append(results, permutationsR(locs, len(locs)-d)...)
	}

	return results
}

//permutationsR takes in a length and a deletion distance. It then recursively
//calls itself, assembling every combination of len(loc) choose length
func permutationsR(locs []loc, length int) [][]loc {
	results := make([][]loc, 0)

	for i, cur := range locs {
		if length == 1 {
			results = append(results, []loc{cur})
		} else {
			items := permutationsR(locs[i+1:], length-1)
			for _, next := range items {
				tmp := make([]loc, 0)
				tmp = append(tmp, cur)
				tmp = append(tmp, next...)
				results = append(results, tmp)
			}
		}
	}

	return results
}

//stringsFromPermutations takes an original string, a list of its rune start indices and widths, and
//an array of inclusions (perm). It goes through each deletion item and builds a string from the original.
//The return value is an array of strings that make up each version of the original string with associated deletions.
func stringsFromPermutations(w string, locs []loc, perms [][]loc) []string {
	results := make([]string, 0)
	results = append(results, w)

	for _, v := range perms {
		s := ""
		for _, l := range v {
			s += w[l.index : l.index+l.width]
		}
		results = append(results, s)
	}
	return results
}
