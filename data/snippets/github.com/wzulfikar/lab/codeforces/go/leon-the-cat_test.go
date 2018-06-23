package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

/**
 * If the inner volume of a pyramid is
 * at least 1000 times bigger than the volume of his body,
 * then he considers that pyramid big enough for him to play
 * all day long and be happy.
 *
 * you need to figure out if a pyramid
 * will make Leon happy or not given the height,
 * width and length of Leon and the inner volume of a pyramid.
 *
 * Input
* The first line of input will consists of four integers
* h (1 ≤ h ≤ 100), w (1 ≤ w ≤ 100), l (1 ≤ l ≤ 1000) and v (1 ≤ v ≤ 1010).
*
* Output
* Print "YES" if the given pyramid makes Leon happy else
* "NO" without the qoutes.
*/

func do(s1 string) string {
	s1.

	if s1 < s2 {
		return -1
	} else if s1 > s2 {
		return 1
	}
	return 0
}

func Test(t *testing.T) {
	a := assert.New(t)

	a.Equal("NO", do("1 1 1 1"))
	a.Equal("NO", do("100 100 1000 1000"))
}

func main() {
	var s1, s2 string
	fmt.Scan(&s1, &s2)
	fmt.Println(do(s1, s2))
}
