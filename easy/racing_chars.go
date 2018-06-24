package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func racing(n int, q string) (int, string) {
	p := n
	a, b := p-1, p+2
	if a < 0 {
		a = 0
	}
	if b > len(q) {
		b = len(q)
	}
	c := strings.Index(q[a:b], "C")
	if c == -1 {
		c = a + strings.Index(q[a:b], "_")
	} else {
		c += a
	}
	if c < p {
		q = q[:c] + "/" + q[c+1:]
	} else if c == p {
		q = q[:c] + "|" + q[c+1:]
	} else {
		q = q[:c] + `\` + q[c+1:]
	}
	return c, q
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)

	scanner.Scan()
	q := scanner.Text()
	c := strings.Index(q, "C")
	if c == -1 {
		c = strings.Index(q, "_")
	}
	fmt.Println(q[:c] + "|" + q[c+1:])
	for scanner.Scan() {
		if len(scanner.Text()) > 0 {
			c, q = racing(c, scanner.Text())
			fmt.Println(q)
		}
	}
}
