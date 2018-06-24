package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var digits map[uint8]int

func init() {
	digits = map[uint8]int{'0': 252, '1': 96,
		'2': 218, '3': 242, '4': 102,
		'5': 182, '6': 190, '7': 224,
		'8': 254, '9': 246}
}

func valid(y int, d bool, x int) bool {
	if d && (x&1 == 0) {
		return false
	}
	return y&x == y
}

func brokenLcd(p, q string) bool {
	var v bool
	t := strings.Fields(p)
	for u := 0; u <= len(t)-len(q)+1; u++ {
		var j int
		v = true
		for i := 0; i < len(q); i++ {
			d := i < len(q)-1 && q[i+1] == '.'
			var x int
			fmt.Sscanf(t[u+j], "%b", &x)
			if !valid(digits[q[i]], d, x) {
				v = false
				break
			}
			if d {
				i++
			}
			j++
		}
		if v {
			return true
		}
	}
	return v
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ";")
		if brokenLcd(s[0], s[1]) {
			fmt.Println("1")
		} else {
			fmt.Println("0")
		}
	}
}
