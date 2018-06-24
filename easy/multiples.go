package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func multiples(x, n int) (r int) {
	var c uint
	r = n
	for (r << 1) <= x {
		r <<= 1
		c++
	}
	for c > 0 {
		c--
		for r+(n<<c) <= x {
			r += n << c
		}
	}
	for r < x {
		r += n
	}
	return r
}

func main() {
	var x, n int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d,%d", &x, &n)
		fmt.Println(multiples(x, n))
	}
}
