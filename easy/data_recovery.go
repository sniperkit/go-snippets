package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func recovery(s, t string) string {
	var k int
	w := strings.Fields(s)
	r := make([]string, len(w))
	for ix, i := range strings.Fields(t) {
		fmt.Sscan(i, &k)
		r[k-1] = w[ix]
	}
	for ix := range r {
		if r[ix] == "" {
			r[ix] = w[len(w)-1]
			break
		}
	}
	return strings.Join(r, " ")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		t := strings.Split(scanner.Text(), ";")
		if len(t) == 2 {
			fmt.Println(recovery(t[0], t[1]))
		}
	}
}
