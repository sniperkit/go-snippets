package basic

import (
	"log"
	"regexp"
	"testing"
)

func TestPatternMatch(t *testing.T) {
	pattern := "^\\d+$"
	re, err := regexp.Compile(pattern)
	log.Println(re, err)

	// match
	log.Println("matched: ", re.Match([]byte("1234")))

	// unmatch
	log.Println("unmatched: ", re.Match([]byte("a1234")))
}
