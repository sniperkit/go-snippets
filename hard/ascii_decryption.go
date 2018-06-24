package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	var (
		n int
		f bool
	)
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		var (
			minc, offset int
			m            [][]int
		)
		s := strings.Split(scanner.Text(), "|")
		fmt.Sscan(strings.TrimSpace(s[0]), &n)
		c := strings.TrimSpace(s[1])[0]
		t := strings.Fields(strings.TrimSpace(s[2]))
		u := make([]int, len(t))
		for ix, i := range t {
			fmt.Sscan(i, &u[ix])
			if ix == 0 {
				minc = u[0]
			} else if u[ix] < minc {
				minc = u[ix]
			}
		}
		for i := 0; i < len(t)-2*n-1; i++ {
			for j := i + n + 1; j < len(t)-n; j++ {
				f = true
				for k := 0; k < n; k++ {
					if u[i+k] != u[j+k] {
						f = false
						break
					}
				}
				if f {
					m = append(m, []int{i, j})
				}
			}
		}
		f = true
		for len(m) > 1 && f {
			f = false
			for i := len(m) - 1; i >= 0; i-- {
				for j := m[i][0]; j < m[i][0]+n; j++ {
					if u[j] == minc {
						m[i], m[len(m)-1], m = m[len(m)-1], nil, m[:len(m)-1]
						f = true
						break
					}
				}
			}
		}
		if len(m) == 0 {
			offset = 32 - minc
		} else {
			offset = int(c) - u[m[0][0]+n-1]
		}
		for _, i := range u {
			fmt.Printf("%c", i+offset)
		}
		fmt.Println()
	}
}
