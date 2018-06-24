package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const k = ` !"#$%&'()*+,-./0123456789:<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`

func gronsfeld(p, q string) string {
	r := make([]byte, len(q))
	for i := 0; i < len(q); i++ {
		r[i] = k[(len(k)+strings.IndexRune(k, rune(q[i]))-int(p[i%len(p)]-'0'))%len(k)]
	}
	return string(r)
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ";")
		fmt.Println(gronsfeld(s[0], s[1]))
	}
}
