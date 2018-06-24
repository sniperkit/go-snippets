package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var m map[rune]rune

func init() {
	codel := `nug
rbc vjnmkf kd yxyqci na rbc zjkfoscdd ew rbc ujllmcp
tc rbkso rbyr ejp mysljylc kd kxveddknmc re jsicpdrysi
de kr kd eoya kw aej icfkici re zjkr`
	decol := `bjv
the public is amazed by the quickness of the juggler
we think that our language is impossible to understand
so it is okay if you decided to quit`

	m = map[rune]rune{' ': ' '}
	var misscod, missdec []rune
	for i := 'a'; i <= 'z'; i++ {
		ix := strings.IndexRune(codel, i)
		if ix != -1 {
			m[i] = rune(decol[ix])
		} else {
			misscod = append(misscod, i)
		}

		ix = strings.IndexRune(decol, i)
		if ix == -1 {
			missdec = append(missdec, i)
		}
	}
	for ix, i := range misscod {
		m[i] = missdec[ix]
	}
}

func mf(a rune) rune {
	return m[a]
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(strings.Map(mf, scanner.Text()))
	}
}
