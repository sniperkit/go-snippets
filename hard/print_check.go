package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var s0, s1, s2, s3 []string

func init() {
	s0 = []string{"", "One", "Two", "Three", "Four", "Five",
		"Six", "Seven", "Eight", "Nine"}
	s1 = []string{"Ten", "Eleven", "Twelve", "Thirteen", "Fourteen",
		"Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen"}
	s2 = []string{"Twenty", "Thirty", "Forty", "Fifty",
		"Sixty", "Seventy", "Eighty", "Ninety"}
	s3 = []string{"Hundred", "Thousand", "Million"}
}

func wrd(a1, a2, a3 byte) (string, bool) {
	if a1+a2+a3 == 0 {
		return "", false
	}
	var r string
	if a1 > 0 {
		r = s0[a1] + s3[0]
	}
	if a2 > 1 {
		r += s2[a2-2] + s0[a3]
	} else if a2 > 0 {
		r += s1[a3]
	} else {
		r += s0[a3]
	}
	return r, true
}

func printCheck(n int) string {
	if n == 0 {
		return "ZeroDollars"
	}
	var r string
	c := make([]byte, 9)
	for i := 8; n > 0 && i >= 0; i-- {
		c[i], n = byte(n%10), n/10
	}
	if s, f := wrd(c[0], c[1], c[2]); f {
		r = s + s3[2]
	}
	if s, f := wrd(c[3], c[4], c[5]); f {
		r += s + s3[1]
	}
	s, _ := wrd(c[6], c[7], c[8])
	return r + s + "Dollars"
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
		fmt.Println(printCheck(n))
	}
}
