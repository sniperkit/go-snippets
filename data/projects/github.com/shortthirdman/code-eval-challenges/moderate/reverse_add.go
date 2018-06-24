package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func rev(n int) (r int) {
	for ; n > 0; n /= 10 {
		r = 10*r + n%10
	}
	return r
}

func revAdd(n int) string {
	var c, r int
	for c < 100 {
		r = rev(n)
		if r == n {
			return fmt.Sprint(c, n)
		}
		n, c = n+r, c+1
	}
	return "not found"
}

func main() {
	var n int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscan(scanner.Text(), &n)
		fmt.Println(revAdd(n))
	}
}
