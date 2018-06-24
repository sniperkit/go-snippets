package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
)

func climbingStairs(a uint) string {
	if a < 4 {
		return fmt.Sprint(a)
	}
	r := big.NewInt(1)
	q := big.NewInt(0)
	for p := big.NewInt(1); a > 1; a-- {
		q.Set(r)
		r.Add(r, p)
		p.Set(q)
	}
	return fmt.Sprint(r)
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		var a uint // line may be empty, default to 0
		fmt.Sscan(scanner.Text(), &a)
		if a > 0 {
			fmt.Println(climbingStairs(a))
		}
	}
}
