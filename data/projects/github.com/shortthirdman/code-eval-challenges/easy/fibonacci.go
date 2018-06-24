package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func fibo(n uint64) uint64 {
	if n < 2 {
		return n
	}
	var (
		r uint64 = 1
		p uint64 = 1
	)
	for ; n > 2; n-- {
		r, p = r+p, r
	}
	return r
}

func main() {
	var n uint64
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscan(scanner.Text(), &n)
		fmt.Println(fibo(n))
	}
}
