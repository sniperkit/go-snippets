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

	results = append(results, Web(query))
	results = append(results, Image(query))
	results = append(results, Video(query))

	return results
}

// STOP0 OMIT

// START1 OMIT
var (
	Web   = fakeSearch("web")
	Image = fakeSearch("image")
	Video = fakeSearch("video")
)

type Search func(query string) Result // HL

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}

// STOP1 OMIT

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := Google("golang") // HL
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)

}
