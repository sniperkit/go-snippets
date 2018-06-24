package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var mat []string

func uniq(x, y, cw, ch int) bool {
	d := make(map[uint8]bool)
	for i := x; i < x+cw; i++ {
		for j := y; j < y+ch; j++ {
			if d[mat[j][i]] {
				return false
			}
			d[mat[j][i]] = true
		}
	}
	return true
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	var (
		w, h int
		res  [][]int
	)
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		h++
		if w == 0 {
			w = len(scanner.Text())
		}
		mat = append(mat, scanner.Text())
	}

	ma, ch, cw := 1, 1, 1
	for x, y := 0, 0; ch <= h; {
		if uniq(x, y, cw, ch) {
			if cw*ch > ma {
				ma = cw * ch
				res = [][]int{[]int{x, y, cw, ch}}
			} else {
				res = append(res, []int{x, y, cw, ch})
			}
			for x+cw < w && uniq(x, y, cw+1, ch) {
				cw++
				ma = cw * ch
				res = [][]int{[]int{x, y, cw, ch}}
			}
		}
		if x+cw < w {
			x++
		} else if y+ch < h {
			x, y = 0, y+1
		} else {
			x, y, cw, ch = 0, 0, 1, ch+1
			for cw*ch < ma {
				cw++
			}
		}
	}

	for _, i := range res {
		for j := i[1]; j < i[1]+i[3]; j++ {
			mat[j] = mat[j][:i[0]] + strings.Repeat("*", i[2]) + mat[j][i[0]+i[2]:]
		}
	}
	for _, i := range mat {
		fmt.Println(i)
	}
}
