package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func do(s string) string {
	// your code here
	return s
}

// put the test case here.
// comment this func before
// submitting to codeforces
func Test(t *testing.T) {
	a := assert.New(t)
	a.Equal("yo", do("yo"))
}

func main() {
	var s string
	fmt.Scan(&s)
	fmt.Println(do(s))
}
