package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	var (
		n    int
		c, d []int
	)
	for i := 0; scanner.Scan(); i++ {
		if n == 0 {
			n = len(scanner.Text())
			d = []int{-n - 1, -n, 1 - n, -1, 0, 1, n - 1, n, n + 1}
		}
		for jx, j := range scanner.Text() {
			if j == '*' {
				c = append(c, i*n+jx)
			}
		}
	}

	for r := 0; r < 10; r++ {
		var c2 []int
		m := make([]int, n*n)
		for _, i := range c {
			for jx, j := range d {
				if j != 0 && i+j >= 0 && i+j < n*n && ((i%n > 0) || (jx%3 > 0)) && ((i%n < (n - 1)) || (jx%3 < 2)) {
					m[i+j]++
				}
			}
		}

		for i, j := 0, 0; i < n*n; i++ {
			switch m[i] {
			case 2:
				for j < len(c) && c[j] < i {
					j++
				}
				if j < len(c) && c[j] == i {
					c2 = append(c2, i)
				}
			case 3:
				c2 = append(c2, i)
			}
		}
		c = c2
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if len(c) > 0 && c[0] == i*n+j {
				fmt.Print("*")
				copy(c, c[1:])
				c = c[:len(c)-1]
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
