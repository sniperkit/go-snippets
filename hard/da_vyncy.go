package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var h map[string][]int

func init() {
	h = make(map[string][]int)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func overlap(a, b string) []int {
	k := a + ";" + b
	if _, found := h[k]; found {
		return h[k]
	}
	if len(a) < len(b) {
		o := overlap(b, a)
		return []int{o[0], o[2], o[1]}
	}
	maxfit := make([]int, 3)
	for i := 0; i < len(a)-1; i++ {
		l := min(len(a)-i, len(b))
		if a[i:i+l-1] == b[:l-1] {
			maxfit = []int{l, i, 0}
			break
		}
	}
	for i := 0; i < len(b)-1; i++ {
		l := len(b) - i
		if maxfit[0] >= l {
			break
		}
		if a[:l-1] == b[i:i+l-1] {
			maxfit = []int{l, 0, i}
			break
		}
	}
	h[k] = maxfit
	return maxfit
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ";")
		for len(s) > 1 {
			maxoverlap := make([]int, 5)
			for ix, i := range s {
				for jx, j := range s {
					if ix == jx {
						continue
					}
					o := overlap(i, j)
					if o[0] > maxoverlap[0] {
						maxoverlap = []int{o[0], ix, jx, o[1], o[2]}
					}
				}
			}
			if maxoverlap[4] == 0 {
				s[maxoverlap[1]] = s[maxoverlap[1]] + s[maxoverlap[2]][maxoverlap[0]:]
			} else if maxoverlap[3] == 0 {
				s[maxoverlap[1]] = s[maxoverlap[2]] + s[maxoverlap[1]][maxoverlap[0]:]
			}
			s[maxoverlap[2]], s[len(s)-1], s = s[len(s)-1], "", s[:len(s)-1]
		}
		fmt.Println(s[0])
	}
}
