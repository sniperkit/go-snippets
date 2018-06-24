package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestStack(t *testing.T) {
	for k, v := range map[string]string{
		"1 2 3 4":   "4 2",
		"10 -2 3 4": "4 -2",
		"1 2 3 4 5": "5 3 1"} {
		if r := stack(k); r != v {
			t.Errorf("failed: stack %s is %s, got %s",
				k, v, r)
		}
	}
}

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
