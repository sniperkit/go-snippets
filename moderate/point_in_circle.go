package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func inCircle(a, b, c float64) bool {
	return a*a+b*b <= c*c
}

func main() {
	f := make([]float64, 5)
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "Center: (%f, %f); Radius: %f; Point: (%f, %f)",
			&f[0], &f[1], &f[2], &f[3], &f[4])
		fmt.Println(inCircle(f[0]-f[3], f[1]-f[4], f[2]))
	}
}
