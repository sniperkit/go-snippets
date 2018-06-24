package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var parent map[int]int

func init() {
	parent = map[int]int{3: 8, 8: 30, 10: 20, 20: 8, 29: 20, 30: 0, 52: 30}
}

func lca(a, b int) int {
	for ; a != 0 && a != b; a = parent[a] {
		for i := b; i != 0; i = parent[i] {
			if a == i {
				return a
			}
		}
	}
	return a
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
		fmt.Println(lca(a, b))
	}
}
