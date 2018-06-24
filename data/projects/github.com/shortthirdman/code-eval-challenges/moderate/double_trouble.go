package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func doubleTrouble(q string) uint {
	var r uint = 1
	n := len(q) / 2
	for i := 0; i < n; i++ {
		if ((q[i] == 'A') && (q[i+n] == 'B')) ||
			((q[i] == 'B') && (q[i+n] == 'A')) {
			r = 0
			break
		} else if (q[i] == '*') && (q[i+n] == '*') {
			r *= 2
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
		fmt.Println(doubleTrouble(scanner.Text()))
	}
}
