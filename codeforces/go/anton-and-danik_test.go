package main

import (
	"fmt"
	"strings"
)

func do(s string) string {
	a := strings.Count(s, "A")
	d := strings.Count(s, "D")
	if a > d {
		return "Anton"
	}
	if a < d {
		return "Danik"
	}
	return "Friendship"
}

// put the test case here.
// comment this func before
// submitting to codeforces
// func Test(t *testing.T) {
// 	a := assert.New(t)
// 	a.Equal("Anton", do("ADAAAA"))
// 	a.Equal("Danik", do("DDDAADA"))
// 	a.Equal("Friendship", do("DADADA"))
// }

// http://codeforces.com/problemset/problem/734/A
func main() {
	var s string
	var n int
	fmt.Scan(&n, &s)
	fmt.Println(do(s))
}
