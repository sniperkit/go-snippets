package main

import (
	"fmt"
	"math"
	"strings"
)

func main() {

	sum := 32654.23
	k := math.Floor(sum / 500)
	if sum > k {
		k += 1
	}
	fmt.Println(k * 500)
	fmt.Println("Hello, playground")

	s := "TyCdPveVfQ9enFPNUgnMyUFu"
	fmt.Println(strings.ToLower(s))
}

