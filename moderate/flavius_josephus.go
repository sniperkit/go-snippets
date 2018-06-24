package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type node struct {
	value int
	next  *node
}

func flavius(n, m int) string {
	t := make([]string, n)
	tail := node{value: n - 1}
	list := &tail
	for i := n - 2; i >= 0; i-- {
		list = &node{i, list}
	}
	tail.next, list = list, &tail

	for i := 0; i < n; i++ {
		for j := 0; j < m-1; j++ {
			list = list.next
		}
		t[i] = fmt.Sprint(list.next.value)
		list.next = list.next.next
	}
	return strings.Join(t, " ")
}

func main() {
	var n, m int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d,%d", &n, &m)
		fmt.Println(flavius(n, m))
	}
}
