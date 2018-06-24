package main

import "fmt"

func robot(f []bool, x, y byte) int {
	if x == 3 && y == 3 {
		return 1
	}
	var ret int
	if x > 0 && !f[x-1+4*y] {
		g := make([]bool, 16)
		copy(g, f)
		g[x-1+4*y] = true
		ret += robot(g, x-1, y)
	}
	if y > 0 && !f[x+4*(y-1)] {
		g := make([]bool, 16)
		copy(g, f)
		g[x+4*(y-1)] = true
		ret += robot(g, x, y-1)
	}
	if x < 3 && !f[x+1+4*y] {
		g := make([]bool, 16)
		copy(g, f)
		g[x+1+4*y] = true
		ret += robot(g, x+1, y)
	}
	if y < 3 && !f[x+4*(y+1)] {
		g := make([]bool, 16)
		copy(g, f)
		g[x+4*(y+1)] = true
		ret += robot(g, x, y+1)
	}
	return ret
}

func main() {
	f := make([]bool, 16)
	f[0] = true
	fmt.Println(robot(f, 0, 0))
}
