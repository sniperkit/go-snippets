package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var (
	ronum []uint
	rostr []string
)

func init() {
	ronum, rostr = []uint{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1},
		[]string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
}

func roman(a uint) (r string) {
	for i := 0; a > 0; {
		if a >= ronum[i] {
			r += rostr[i]
			a -= ronum[i]
		} else {
			i++
		}
	}
	return r
}

func main() {
	var a uint
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscan(scanner.Text(), &a)
		fmt.Println(roman(a))
	}
}
