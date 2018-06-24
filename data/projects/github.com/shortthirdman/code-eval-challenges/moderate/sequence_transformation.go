package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func seqTrans(p, q string) bool {
	p = "^" + strings.Replace(strings.Replace(p, "0", "A+", -1), "1", "((A+)|(B+))", -1) + "$"
	seqRegex := regexp.MustCompile(p)
	return seqRegex.MatchString(q)
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Fields(scanner.Text())
		if seqTrans(s[0], s[1]) {
			fmt.Println("Yes")
		} else {
			fmt.Println("No")
		}
	}
}
