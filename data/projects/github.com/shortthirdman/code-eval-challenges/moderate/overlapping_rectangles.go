package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type point struct {
	x, y int
}

func within(a, b, c point) bool {
	return c.x >= a.x && c.x <= b.x && c.y >= b.y && c.y <= a.y
}

func overlapping(a, b, c, d point) bool {
	for _, i := range []point{a, point{a.x, b.y}, point{b.x, a.y}, b} {
		if within(c, d, i) {
			return true
		}
	}
	for _, i := range []point{c, point{c.x, d.y}, point{d.x, c.y}, d} {
		if within(a, b, i) {
			return true
		}
	}
	return false
}

func main() {
	var a, b, c, d point
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d,%d,%d,%d,%d,%d,%d,%d",
			&a.x, &a.y, &b.x, &b.y, &c.x, &c.y, &d.x, &d.y)
		if overlapping(a, b, c, d) {
			fmt.Println("True")
		} else {
			fmt.Println("False")
		}
	}
}
