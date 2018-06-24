package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var h map[int]int

func init() {
	h = make(map[int]int)
}

func floors(e, s int) int {
	switch {
	case e == 0 || s == 0:
		return 0
	case e == 1:
		return s
	case s == 1:
		return 1
	case e > s:
		return floors(s, s)
	}
	if hf, f := h[e<<16+s]; f {
		return hf
	}
	h[e<<16+s] = floors(e-1, s-1) + floors(e, s-1) + 1
	return h[e<<16+s]
}

func main() {
	var e, s int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		var n int
		fmt.Sscanf(scanner.Text(), "%d %d", &e, &s)
		for floors(e, n) < s {
			n++
		}
		fmt.Println(n)
	}
}
