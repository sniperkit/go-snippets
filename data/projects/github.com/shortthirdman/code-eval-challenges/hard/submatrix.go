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
		m, ps     [][]int
		maxmat, x int
	)
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		var r []int
		t := strings.Fields(scanner.Text())
		for _, i := range t {
			fmt.Sscan(i, &x)
			r = append(r, x)
		}
		m = append(m, r)
	}
	n := len(m)
	for i := 0; i < n; i++ {
		r := make([]int, n)
		ps = append(ps, r)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			ps[j][i] = m[j][i]
			if j > 0 {
				ps[j][i] += ps[j-1][i]
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			var minmat, submat int
			for k := 0; k < n; k++ {
				submat += ps[j][k]
				if i > 0 {
					submat -= ps[i-1][k]
				}
				if submat < minmat {
					minmat = submat
				}
				if submat-minmat > maxmat {
					maxmat = submat - minmat
				}
			}
		}
	}
	fmt.Println(maxmat)
}
