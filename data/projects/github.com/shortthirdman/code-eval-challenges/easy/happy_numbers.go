package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func happy(a uint) (r uint) {
	for a > 0 {
		b := a % 10
		r += b * b
		a /= 10
	}
	return r
}

func happyNumbers(a uint) bool {
	b := []uint{a}
	for i := 0; a > 1 && i < 127; i++ {
		a = happy(a)
		for _, j := range b {
			if j == a {
				return false
			}
		}
		b = append(b, a)
	}
	return a == 1
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
		if happyNumbers(a) {
			fmt.Println(1)
		} else {
			fmt.Println(0)
		}
	}
}
