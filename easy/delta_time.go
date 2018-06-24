package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func delta(h1, m1, s1, h2, m2, s2 int) string {
	t := abs((h1-h2)*3600 + (m1-m2)*60 + s1 - s2)
	return fmt.Sprintf("%02d:%02d:%02d", t/3600, (t/60)%60, t%60)
}

func main() {
	var h1, m1, s1, h2, m2, s2 int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d:%d:%d %d:%d:%d", &h1, &m1, &s1, &h2, &m2, &s2)
		fmt.Println(delta(h1, m1, s1, h2, m2, s2))
	}
}
