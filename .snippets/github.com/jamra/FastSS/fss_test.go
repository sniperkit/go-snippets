package fss

import (
	"bufio"
	"os"
	"testing"
)

var fss *Fss

func TestBuildFSS(t *testing.T) {
	fss = NewFss(2)
	filename := "/usr/share/dict/words"
	f, err := os.Open(filename)
	if err != nil {
		t.Error("Unable to open file %s", filename)
	}

	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		count++
		if count == 2000 {
			break
		}
		fss.Insert(scanner.Text())
	}

	if scanner.Err() != nil {
		t.Errorf("Could not read lines from file `%s`", filename)
	}
}

//I'm not sure if this is needed as I am using a third party library that is supposed to be tested by industry standards
func TestHashWord(t *testing.T) {

}

func TestSplitWord(t *testing.T) {
	//TODO: Add more words with different distances.
	table := []struct {
		word     string
		distance int
		expected []string
	}{
		{word: "table", distance: 1, expected: []string{"tabl", "tabe", "tale", "tble", "able"}},
	}

	for _, tab := range table {
		perm := getpermutations(tab.word, tab.distance+1)

		for _, str := range tab.expected {
			if !contains(perm, str) {
				t.Errorf("Expected string %s.", str)
			}
		}
	}
}

func qcontains(strs []QueryResult, s string) bool {
	for _, str := range strs {
		if str.S == s {
			return true
		}
	}
	return false
}

func contains(strs []string, s string) bool {
	for _, str := range strs {
		if str == s {
			return true
		}
	}
	return false
}

func TestLookup(t *testing.T) {
	//TODO: This should build an FSS for each query item. This would make sure I can search
	//different sizes of deletion neighborhoods
	table := []struct {
		query    string
		matches  []string
		distance int
	}{
		{query: "abaft", matches: []string{"abaft"}, distance: 0},
		{query: "abaft", matches: []string{"abaft", "abaff"}, distance: 1},
	}

	for _, tabl := range table {
		results := fss.Search(tabl.query)
		for _, w := range tabl.matches {
			if !qcontains(results, w) {
				t.Errorf("Expected %s", w)
			}
		}
	}
}

func TestSplitQuery(t *testing.T) {
	//Because splitting a query is not the same as splitting a full word, the result is multiple splits
}
