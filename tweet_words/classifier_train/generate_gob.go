package main

import (
	"github.com/jbrukh/bayesian"
	"io/ioutil"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createMapOfStopWord(fileName string) map[string]bool {
	m := make(map[string]bool)
	stopWords := convFileToListWords(fileName)
	for _, word := range stopWords {
		m[word] = true
	}
	return m
}

func readFromFile(fileName string) (dat string) {
	datB, err := ioutil.ReadFile(fileName)
	check(err)
	dat = string(datB)
	return
}

func convFileToListWords(fileName string) (datList []string) {

	dat := readFromFile(fileName)

	datList = strings.FieldsFunc(string(dat), func(r rune) bool {
		switch r {
		case '\n', ' ', '\t':
			return true
		}
		return false
	})
	return
}

func removeStopWords(stopWords map[string]bool, datList []string) (datListNoStops []string) {

	for _, word := range datList {
		isStopWord := stopWords[word] || strings.HasPrefix(word, "http:") || strings.HasPrefix(word, "https:") || strings.HasPrefix(word, "#")
		if !isStopWord {
			datListNoStops = append(datListNoStops, word)
		}
	}
	return
}

func main() {

	const (
		Positive bayesian.Class = "Positive"
		Negative bayesian.Class = "Negative"
	)

	stopWords := createMapOfStopWord("stopwords.txt")
	classifier := bayesian.NewClassifier(Positive, Negative)
	goodStuff := removeStopWords(stopWords, convFileToListWords("training-1.txt"))
	badStuff := removeStopWords(stopWords, convFileToListWords("training-0.txt"))
	classifier.Learn(goodStuff, Positive)
	classifier.Learn(badStuff, Negative)
	classifier.WriteToFile("classifier.gob")
}
