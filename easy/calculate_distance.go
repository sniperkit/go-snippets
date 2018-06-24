package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

func sqrt(a int) int {
	return int(math.Sqrt(float64(a)))
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
		fmt.Sscanf(scanner.Text(), "(%d, %d) (%d, %d)", &x1, &y1, &x2, &y2)
		x, y := x1-x2, y1-y2
		fmt.Println(sqrt(x*x + y*y))
	}
}
