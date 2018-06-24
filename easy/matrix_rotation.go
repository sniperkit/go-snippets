package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func sq(a int) (r int) {
	for r = 1; r*r < a; r++ {
	}
	return r
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		m := strings.Fields(scanner.Text())
		l := sq(len(m))
		for i := 0; i < l; i++ {
			for j := l - 1; j >= 0; j-- {
				fmt.Print(m[j*l+i])
				if i < l-1 || j > 0 {
					fmt.Print(" ")
				} else {
					fmt.Print("\n")
				}
			}
		}
	}
}
