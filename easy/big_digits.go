package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	digits = `
-**----*--***--***---*---****--**--****--**---**--
*--*--**-----*----*-*--*-*----*-------*-*--*-*--*-
*--*---*---**---**--****-***--***----*---**---***-
*--*---*--*-------*----*----*-*--*--*---*--*----*-
-**---***-****-***-----*-***---**---*----**---**--
--------------------------------------------------
`
	w, h = 5, 6
)

func bigDigits(q string) string {
	r := make([]string, h)
	for _, i := range q {
		if i >= '0' && i <= '9' {
			for j := 0; j < h; j++ {
				pos := 1 + j*w*10 + j + int(i-'0')*w
				r[j] += digits[pos : pos+w]
			}
		}
	}
	return strings.Join(r, "\n")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(bigDigits(scanner.Text()))
	}
}
