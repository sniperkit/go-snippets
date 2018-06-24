package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func spiral(a, b int, q string) string {
	var (
		r            []string
		i, j, tn, tw int
	)
	t := strings.Fields(q)
	te, ts := b-1, a-1
	for {
		for j <= te {
			r = append(r, t[i*b+j])
			j++
		}
		j--
		i++
		tn++
		if i > ts {
			break
		}
		for i <= ts {
			r = append(r, t[i*b+j])
			i++
		}
		i--
		j--
		te--
		if j < tw {
			break
		}
		for j >= tw {
			r = append(r, t[i*b+j])
			j--
		}
		j++
		i--
		ts--
		if i < tn {
			break
		}
		for i >= tn {
			r = append(r, t[i*b+j])
			i--
		}
		i++
		j++
		tw++
		if j > te {
			break
		}
	}
	return strings.Join(r, " ")
}

func main() {
	var a, b int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ";")
		fmt.Sscan(s[0], &a)
		fmt.Sscan(s[1], &b)
		fmt.Println(spiral(a, b, s[2]))
	}
}
