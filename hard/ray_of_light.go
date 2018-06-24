package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var (
	code map[rune]int
	tile map[int]rune
)

func init() {
	code = map[rune]int{' ': 0, '\\': 0, '/': 0, 'X': 0, '*': 1, '#': 2, 'o': 3}
	tile = map[int]rune{0: ' ', 1: '*', 2: '#', 3: 'o'}
}

type ray struct {
	x   int
	y   int
	dir int
	ttl int
}

func (a ray) nextTile(s []int) (int, int, int) {
	switch a.dir {
	case 0:
		if a.x == 9 || a.y == 0 {
			return -1, 0, 0
		}
		return s[a.x+1+(a.y-1)*10], a.x + 1, a.y - 1
	case 1:
		if a.x == 9 || a.y == 9 {
			return -1, 0, 0
		}
		return s[a.x+1+(a.y+1)*10], a.x + 1, a.y + 1
	case 2:
		if a.x == 0 || a.y == 9 {
			return -1, 0, 0
		}
		return s[a.x-1+(a.y+1)*10], a.x - 1, a.y + 1
	case 3:
		if a.x == 0 || a.y == 0 {
			return -1, 0, 0
		}
		return s[a.x-1+(a.y-1)*10], a.x - 1, a.y - 1
	}
	return -1, 0, 0
}

func (a ray) lrTile(s []int) (int, int, int, int, int, int) {
	switch a.dir {
	case 0:
		return s[a.x+(a.y-1)*10], a.x, a.y - 1, s[a.x+1+a.y*10], a.x + 1, a.y
	case 1:
		return s[a.x+1+a.y*10], a.x + 1, a.y, s[a.x+(a.y+1)*10], a.x, a.y + 1
	case 2:
		return s[a.x+(a.y+1)*10], a.x, a.y + 1, s[a.x-1+a.y*10], a.x - 1, a.y
	case 3:
		return s[a.x-1+a.y*10], a.x - 1, a.y, s[a.x+(a.y-1)*10], a.x, a.y - 1
	}
	return -1, 0, 0, 0, 0, 0
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	s := make([]int, 100)
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		var todo []ray
		been := make(map[int][]ray)
		for ix, i := range scanner.Text() {
			if i == '\\' || i == '/' || i == 'X' {
				x, y := ix%10, ix/10
				if i == '\\' || i == 'X' {
					todo = append(todo, ray{x, y, 1, 20})
					todo = append(todo, ray{x, y, 3, 20})
				}
				if i == '/' || i == 'X' {
					todo = append(todo, ray{x, y, 0, 20})
					todo = append(todo, ray{x, y, 2, 20})
				}
			}
			s[ix] = code[i]
		}

		for len(todo) > 0 {
			b := todo[0]
			todo = todo[1:]

			nt, x, y := b.nextTile(s)
			if b.ttl > 1 {
				if nt == 0 || nt == 1 {
					if nb, f := been[x+10*y]; f {
						for _, i := range nb {
							if i.dir == b.dir && i.ttl >= b.ttl {
								continue
							}
						}
					}
					if nt == 0 {
						todo = append(todo, ray{x, y, b.dir, b.ttl - 1})
					} else {
						todo = append(todo, ray{x, y, b.dir, b.ttl})
						todo = append(todo, ray{x, y, (b.dir + 1) % 4, b.ttl})
						todo = append(todo, ray{x, y, (b.dir + 3) % 4, b.ttl})
					}
				} else if nt == 2 {
					lt, lx, ly, rt, rx, ry := b.lrTile(s)
					if lt == 0 {
						todo = append(todo, ray{lx, ly, (b.dir + 3) % 4, b.ttl - 1})
					} else if lt == 1 {
						todo = append(todo, ray{lx, ly, (b.dir + 3) % 4, b.ttl})
					}
					if rt == 0 {
						todo = append(todo, ray{rx, ry, (b.dir + 1) % 4, b.ttl - 1})
					} else if lt == 1 {
						todo = append(todo, ray{rx, ry, (b.dir + 1) % 4, b.ttl})
					}
				}
			}
			if nb, f := been[b.x+10*b.y]; f {
				for ix, i := range nb {
					if i.dir == b.dir {
						been[b.x+10*b.y][ix] = b
						f = false
						break
					}
				}
				if f {
					been[b.x+10*b.y] = append(been[b.x+10*b.y], b)
				}
			} else {
				been[b.x+10*b.y] = []ray{b}
			}
		}

		for i := 0; i < 100; i++ {
			if s[i] == 0 {
				b, f := been[i]
				if f {
					p := 0
					for _, j := range b {
						if j.dir == 0 || j.dir == 2 {
							p |= 1
						} else {
							p |= 2
						}
					}
					switch p {
					case 1:
						fmt.Printf("/")
					case 2:
						fmt.Printf(`\`)
					case 3:
						fmt.Printf("X")
					}
				} else {
					fmt.Printf(" ")
				}
			} else {
				fmt.Printf("%c", tile[s[i]])
			}
		}
		fmt.Println()
	}
}
