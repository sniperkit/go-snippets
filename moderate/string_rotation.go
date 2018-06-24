package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func stringRotation(p, q string) bool {
	var n int
	for i := strings.Index(q, string(p[0])); i >= 0; i += n {
		if strings.HasSuffix(q, string(p[0:len(p)-i])) && (i == 0 || strings.HasPrefix(q, string(p[len(p)-i:len(p)]))) {
			return true
		}
		n = 1 + strings.Index(q[i+1:len(q)], string(p[0]))
		if n == 0 {
			break
		}
	}
	return false
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ",")
		if stringRotation(s[0], s[1]) {
			fmt.Println("True")
		} else {
			fmt.Println("False")
		}
	}
}
