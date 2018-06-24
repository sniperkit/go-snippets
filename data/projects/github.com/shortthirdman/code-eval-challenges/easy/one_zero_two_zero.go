package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func xz(n, a int) bool {
	for a > 0 && n >= 0 {
		if a&1 == 0 {
			n--
		}
		a >>= 1
	}
	return n == 0
}

func main() {
	var n, a int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		var r int
		fmt.Sscanf(scanner.Text(), "%d %d", &n, &a)
		for i := 1 << uint(n); i <= a; i++ {
			if xz(n, i) {
				r++
			}
		}
		fmt.Println(r)
	}
}
