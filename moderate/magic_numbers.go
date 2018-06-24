package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func isMagic(a int) bool {
	var (
		dig, r uint
		ns     []uint
	)
	for a > 0 {
		r = uint(a % 10)
		if r == 0 || dig&(1<<r) > 0 {
			return false
		}
		dig |= 1 << r
		ns = append(ns, r)
		a /= 10
	}
	dig, r = 0, 0
	for _ = range ns {
		r = (r + ns[(uint(len(ns))-1-r)]) % uint(len(ns))
		if dig&(1<<r) > 0 {
			return false
		}
		dig |= 1 << r
	}
	return r == 0
}

var magic []int

func init() {
	for i := 1; i <= 9876; i++ {
		if isMagic(i) {
			magic = append(magic, i)
		}
	}
}

func magicNumbers(a, b int) string {
	var r []string
	for i := 0; i < len(magic) && magic[i] <= b; i++ {
		if magic[i] >= a {
			r = append(r, fmt.Sprint(magic[i]))
		}
	}
	if len(r) == 0 {
		return "-1"
	}
	return strings.Join(r, " ")
}

func main() {
	var a, b int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d %d", &a, &b)
		fmt.Println(magicNumbers(a, b))
	}
}
