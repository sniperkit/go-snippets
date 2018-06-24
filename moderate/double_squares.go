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
	sq := make(map[int]bool)
	for i := 0; i < 46341; i++ {
		sq[i*i] = true
	}

	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	scanner.Scan()
	var n int
	fmt.Sscan(scanner.Text(), &n)
	for i := 0; i < n; i++ {
		scanner.Scan()
		var x int
		fmt.Sscan(scanner.Text(), &x)
		num, l, bot, top := 0, make(map[int]bool), sqrt(x/2), sqrt(x)+1
		for j := bot; j < top; j++ {
			t := x - j*j
			if sq[t] && !l[t] {
				l[j*j], num = true, num+1
			}
		}
		fmt.Println(num)
	}
}
