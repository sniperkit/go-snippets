package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func stackpush(stack []int, n int) []int {
	return append(stack, n)

}

func stackpop(stack []int) ([]int, int) {
	var n int
	if len(stack) > 0 {
		n, stack = stack[len(stack)-1], stack[:len(stack)-1]
	}
	return stack, n
}

func stack(q string) string {
	var (
		stack []int
		n     int
		r     []string
	)
	t := strings.Fields(q)
	for _, i := range t {
		fmt.Sscan(i, &n)
		stack = stackpush(stack, n)
	}
	for len(stack) > 0 {
		stack, n = stackpop(stack)
		r = append(r, fmt.Sprint(n))
		stack, _ = stackpop(stack)
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
		fmt.Println(stack(scanner.Text()))
	}
}
