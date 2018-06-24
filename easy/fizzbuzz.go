package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type node struct {
	value string
	next  *node
}

func (a node) isEmpty() bool {
	return len(a.value) == 0
}

func fizzbuzz(f, b, n uint) string {
	var r []string
	tail := node{value: "FB"}
	list := &tail
	for i, j := f-1, b-1; i > 0 || j > 0; i, j = (i+f-1)%f, (j+b-1)%b {
		var v string
		if i == 0 {
			v = "F"
		} else if j == 0 {
			v = "B"
		}
		list = &node{v, list}
	}
	tail.next = list

	for i := uint(1); i <= n; i++ {
		if list.isEmpty() {
			r = append(r, fmt.Sprint(i))
		} else {
			r = append(r, list.value)
		}
		list = list.next
	}
	return strings.Join(r, " ")
}

func main() {
	var f, b, n uint
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d %d %d", &f, &b, &n)
		fmt.Println(fizzbuzz(f, b, n))
	}
}
