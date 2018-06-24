package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func cls(a [][]rune) {
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			a[i][j] = ' '
		}
	}
}

func main() {
	var (
		cx     int
		cy     int
		c      rune
		insert bool
	)
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	term := make([][]rune, 10)
	for i := 0; i < 10; i++ {
		term[i] = make([]rune, 10)
	}
	cls(term)

	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		c = rune(scanner.Text()[0])
		if c == '\n' {
			continue
		} else if c == '^' {
			scanner.Scan()
			c = rune(scanner.Text()[0])
			if c != '^' {
				switch c {
				case 'c':
					cls(term)
				case 'h':
					cx, cy = 0, 0
				case 'b':
					cx = 0
				case 'd':
					if cy < 9 {
						cy++
					}
				case 'u':
					if cy > 0 {
						cy--
					}
				case 'l':
					if cx > 0 {
						cx--
					}
				case 'r':
					if cx < 9 {
						cx++
					}
				case 'e':
					for i := cx; i < 10; i++ {
						term[cy][i] = ' '
					}
				case 'i':
					insert = true
				case 'o':
					insert = false
				default:
					cy = int(c - '0')
					scanner.Scan()
					cx = int(rune(scanner.Text()[0]) - '0')
				}
				continue
			}
		}
		if insert {
			for i := cx + 1; i < 10; i++ {
				term[cy][i] = term[cy][i-1]
			}
		}
		term[cy][cx] = c
		if cx < 9 {
			cx++
		}
	}
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			fmt.Printf("%c", term[i][j])
		}
		fmt.Println()
	}
}
