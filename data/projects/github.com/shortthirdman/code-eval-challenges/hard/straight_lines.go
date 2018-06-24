package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type point struct {
	x, y int
}

func inLine(a, b, c point) bool {
	dx1, dy1, dx2, dy2 := a.x-b.x, a.y-b.y, a.x-c.x, a.y-c.y
	return dx1*dy2 == dx2*dy1
}

func line(a, b int, t []point) bool {
	for i := 0; i < len(t); i++ {
		if i != a && i != b && inLine(t[i], t[a], t[b]) {
			return i >= b
		}
	}
	return false
}

func main() {
	var x, y int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " | ")
		var t []point
		for _, i := range s {
			fmt.Sscanf(i, "%d %d", &x, &y)
			t = append(t, point{x, y})
		}
		var r int
		for i := 0; i < len(t)-2; i++ {
			for j := i + 1; j < len(t)-1; j++ {
				if line(i, j, t) {
					r++
				}
			}
		}
		fmt.Println(r)
	}
}
