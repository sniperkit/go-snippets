package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func simpleSort(q string) string {
	t := strings.Fields(q)
	u := make([]float64, len(t))
	for ix, i := range t {
		fmt.Sscan(i, &u[ix])
	}
	sort.Float64s(u)
	for ix, i := range u {
		t[ix] = fmt.Sprintf("%.3f", i)
	}
	return strings.Join(t, " ")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(simpleSort(scanner.Text()))
	}
}
