package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var (
	phrases = []string{", yeah!", ", this is crazy, I tell ya.",
		", can U believe this?", ", eh?", ", aw yea.",
		", yo.", "? No way!", ". Awesome!"}
)

func slang(q string, c int, l bool) (string, int, bool) {
	if q == "." || q == "!" || q == "?" {
		if l {
			return phrases[c], (c + 1) % len(phrases), false
		}
		return q, c, true
	}
	return q, c, l
}

func main() {
	var (
		c int
		l bool
		t string = "\n"
	)
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		t, c, l = slang(scanner.Text(), c, l)
		fmt.Print(t)
	}
	if t != "\n" {
		fmt.Println()
	}
}
