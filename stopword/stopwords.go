package stopword

import (
	"io/ioutil"
	"regexp"
	"strings"
)

var stopWords map[string]bool
var file = "stopword/stopwords.txt"

func loadStopWords() {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	// Replace all newlines with '|'
	newlineRegex := regexp.MustCompile("\r?\n")
	stopWordsString := newlineRegex.ReplaceAllString(string(data), " ")
	stopWordsArr := strings.Split(stopWordsString, " ")

	// Make set
	stopWords = make(map[string]bool)
	for _, word := range stopWordsArr {
		stopWords[word] = true
	}
}

// IsStopWord checks if a word is stopword or not
func IsStopWord(s string) bool {
	if stopWords == nil {
		loadStopWords()
	}
	return stopWords[s]
}
