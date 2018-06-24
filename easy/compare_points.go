package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func comparePoints(x1, y1, x2, y2 int) string {
	switch {
	case x1 == x2:
		if y1 == y2 {
			return "here"
		} else if y1 < y2 {
			return "N"
		} else {
			return "S"
		}
	case y1 == y2:
		if x1 < x2 {
			return "E"
		} else {
			return "W"
		}
	case x1 < x2:
		if y1 < y2 {
			return "NE"
		} else {
			return "SE"
		}
	default:
		if y1 < y2 {
			return "NW"
		} else {
			return "SW"
		}
	}
}

func main() {
	var x1, y1, x2, y2 int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d %d %d %d", &x1, &y1, &x2, &y2)
		fmt.Println(comparePoints(x1, y1, x2, y2))
	}
}
