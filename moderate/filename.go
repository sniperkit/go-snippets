package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func filename(q string) string {
	var r []string
	f := strings.Fields(q)
	p := regexp.QuoteMeta(f[0])
	p = strings.Replace(p, `\?`, ".", -1)
	p = strings.Replace(p, `\*`, ".*", -1)
	p = strings.Replace(p, `\[`, "[", -1)
	p = strings.Replace(p, `\]`, "]", -1)
	p = "^" + p + "$"
	pattern := regexp.MustCompile(p)
	for i := 1; i < len(f); i++ {
		if pattern.MatchString(f[i]) {
			r = append(r, f[i])
		}
	}
	if len(r) > 0 {
		return strings.Join(r, " ")
	} else {
		return "-"
	}
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(filename(scanner.Text()))
	}
}
