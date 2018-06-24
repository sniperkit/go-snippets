package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func numberPairs(q string) string {
	var (
		n int
		r []string
	)
	t := strings.Split(q, ";")
	u := strings.Split(t[0], ",")
	fmt.Sscan(t[1], &n)
	v, vk := make(map[int]bool, len(u)), make([]int, len(u))
	for i := range u {
		fmt.Sscan(u[i], &vk[i])
		v[vk[i]] = true
	}
	for i := 0; i < len(vk) && 2*vk[i] < n; i++ {
		if v[n-vk[i]] {
			r = append(r, fmt.Sprint(vk[i], ",", n-vk[i]))
		}
	}
	if len(r) == 0 {
		return "NULL"
	}
	return strings.Join(r, ";")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(numberPairs(scanner.Text()))
	}
}
