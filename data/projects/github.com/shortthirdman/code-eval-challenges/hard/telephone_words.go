package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func telephone(s, t string) []string {
	if len(s) == 0 {
		return []string{t}
	}
	var r []string
	switch s[0] {
	case '0':
		r = telephone(s[1:len(s)], t+"0")
	case '1':
		r = telephone(s[1:len(s)], t+"1")
	case '2':
		r = append(telephone(s[1:len(s)], t+"a"),
			append(telephone(s[1:len(s)], t+"b"),
				telephone(s[1:len(s)], t+"c")...)...)
	case '3':
		r = append(telephone(s[1:len(s)], t+"d"),
			append(telephone(s[1:len(s)], t+"e"),
				telephone(s[1:len(s)], t+"f")...)...)
	case '4':
		r = append(telephone(s[1:len(s)], t+"g"),
			append(telephone(s[1:len(s)], t+"h"),
				telephone(s[1:len(s)], t+"i")...)...)
	case '5':
		r = append(telephone(s[1:len(s)], t+"j"),
			append(telephone(s[1:len(s)], t+"k"),
				telephone(s[1:len(s)], t+"l")...)...)
	case '6':
		r = append(telephone(s[1:len(s)], t+"m"),
			append(telephone(s[1:len(s)], t+"n"),
				telephone(s[1:len(s)], t+"o")...)...)
	case '7':
		r = append(telephone(s[1:len(s)], t+"p"),
			append(telephone(s[1:len(s)], t+"q"),
				append(telephone(s[1:len(s)], t+"r"),
					telephone(s[1:len(s)], t+"s")...)...)...)
	case '8':
		r = append(telephone(s[1:len(s)], t+"t"),
			append(telephone(s[1:len(s)], t+"u"),
				telephone(s[1:len(s)], t+"v")...)...)
	case '9':
		r = append(telephone(s[1:len(s)], t+"w"),
			append(telephone(s[1:len(s)], t+"x"),
				append(telephone(s[1:len(s)], t+"y"),
					telephone(s[1:len(s)], t+"z")...)...)...)
	}
	return r
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(strings.Join(telephone(scanner.Text(), ""), ","))
	}
}
