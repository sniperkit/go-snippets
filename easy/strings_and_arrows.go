package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func stringsAndArrows(s string) (r uint) {
	for i := 0; i < len(s)-4; i++ {
		if s[i:i+5] == ">>-->" || s[i:i+5] == "<--<<" {
			r++
		}
	}
	return r
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(stringsAndArrows(scanner.Text()))
	}
}
