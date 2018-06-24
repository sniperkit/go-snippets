package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

var c []int

func init() {
	c = make([]int, 26)
}

func clean(char rune) rune {
	if char >= 'a' && char <= 'z' {
		return char
	}
	return -1
}

func beauty(s string) (r int) {
	t := strings.Map(clean, strings.ToLower(s))
	for _, i := range t {
		c[i-'a']++
	}
	sort.Ints(c)
	for i := 26; i > 0 && c[i-1] > 0; i-- {
		r += i * c[i-1]
		c[i-1] = 0
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
		fmt.Println(beauty(scanner.Text()))
	}
}
