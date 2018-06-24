package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const dn, dz = "negative", "zero"

var (
	d0, d1, d2 map[string]int
	d3         []string
)

func init() {
	d0 = map[string]int{"zero": 0, "one": 1, "two": 2,
		"three": 3, "four": 4, "five": 5, "six": 6,
		"seven": 7, "eight": 8, "nine": 9}
	d1 = map[string]int{"ten": 10, "eleven": 11, "twelve": 12,
		"thirteen": 13, "fourteen": 14, "fifteen": 15, "sixteen": 16,
		"seventeen": 17, "eighteen": 18, "nineteen": 19}
	d2 = map[string]int{"twenty": 20, "thirty": 30, "forty": 40,
		"fifty": 50, "sixty": 60, "seventy": 70, "eighty": 80,
		"ninety": 90}
	d3 = []string{"hundred", "thousand", "million"}
}

func textToNumber(s string) int {
	t := strings.Fields(s)
	var n, h, q, r int
	if t[0] == dn {
		n = -1
		t = t[1:]
	} else {
		n = 1
	}
	for len(t) > 0 {
		if t[0] == dz {
			q = 0
		} else if v, f := d0[t[0]]; f {
			h = v
		} else if v, f := d1[t[0]]; f {
			q += v
		} else if v, f := d2[t[0]]; f {
			q += v
		} else if t[0] == d3[0] {
			q, h = h*100, 0
		} else if t[0] == d3[1] {
			r, h, q = r+(h+q)*1000, 0, 0
		} else if t[0] == d3[2] {
			r, h, q = (h+q)*1000000, 0, 0
		}
		t = t[1:]
	}
	return n * (r + q + h)
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(textToNumber(scanner.Text()))
	}
}
