// +build OMIT
// The comments in these files for the presentation tool used to generate the slides.

package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Result string

// START0 OMIT
func Google(query string) []Result {
	var results []Result
	resultsChan := make(chan Result)

	go func() { resultsChan <- Web(query) }()
	go func() { resultsChan <- Image(query) }()
	go func() { resultsChan <- Video(query) }()

	for i := 0; i < 3; i++ {
		results = append(results, <-resultsChan)
	}

	return results
}

// STOP0 OMIT

var (
	Web   = fakeSearch("web")
	Image = fakeSearch("image")
	Video = fakeSearch("video")
)

type Search func(query string) Result

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := Google("golang")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)

}
