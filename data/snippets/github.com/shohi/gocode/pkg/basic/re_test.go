package basic

import (
	"log"
	"regexp"
	"testing"
)

func TestReMatch(t *testing.T) {
	ptn := "reading"
	str := "Reading book is good"

	re := regexp.MustCompile("(?i).*" + ptn + ".*")
	log.Println(re.Match([]byte(str)))
}
