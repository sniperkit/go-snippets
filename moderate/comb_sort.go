package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const shrink = 1.25

func combSort(q string) uint {
	var (
		n, l uint
		k, g int
		c    bool
	)
	u := strings.Fields(q)
	g, c = len(u), false
	v := make([]int, len(u))
	for ix, i := range u {
		fmt.Sscan(i, &k)
		v[ix] = k
	}
	for g > 1 || c {
		n, c = n+1, false
		if g > 1 {
			g = int(float32(g) / shrink)
		}
		for j := 0; j < len(u)-g; j++ {
			if u[j] > u[j+g] {
				u[j], u[j+g], c, l = u[j+g], u[j], true, n
			}
		}
	}
	return l
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(combSort(scanner.Text()))
	}
}
