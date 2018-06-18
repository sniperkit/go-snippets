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
	a.Equal("house", do("HoUse"))
	a.Equal("VIP", do("ViP"))
	a.Equal("matrix", do("maTRIx"))
}

// http://codeforces.com/problemset/problem/59/A
func main() {
	var in string
	fmt.Scan(&in)
}
