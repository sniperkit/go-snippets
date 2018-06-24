package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func lakes(q string) int {
	l := make(map[int]bool)
	s := strings.Split(q, "|")
	last := make([]int, len(s[0])/2+1)
	for ix, i := range s {
		var left int
		curr := make([]int, len(s[0])/2)
		t := strings.Fields(strings.TrimSpace(i))
		for jx, j := range t {
			if j == "o" {
				if ix > 0 {
					if jx > 0 && last[jx-1] > 0 {
						curr[jx] = last[jx-1]
					} else if last[jx] > 0 {
						curr[jx] = last[jx]
					}
					if last[jx+1] > 0 {
						if curr[jx] == 0 {
							curr[jx] = last[jx+1]
						} else if curr[jx] != last[jx+1] {
							tmp := last[jx+1]
							for k := jx + 1; k < len(last); k++ {
								if last[k] == tmp {
									last[k] = curr[jx]
								}
							}
							for k := 0; k < jx; k++ {
								if curr[k] == tmp {
									curr[k] = curr[jx]
								}
							}
							delete(l, tmp)
						}
					}
				}
				if left > 0 {
					if curr[jx] == 0 {
						curr[jx] = left
					} else if curr[jx] != left {
						tmp := left
						for k := jx; k < len(last); k++ {
							if last[k] == tmp {
								last[k] = curr[jx]
							}
						}
						for k := 0; k < jx; k++ {
							if curr[k] == tmp {
								curr[k] = curr[jx]
							}
						}
						delete(l, tmp)
					}
				}
				if curr[jx] == 0 {
					curr[jx] = ix<<5 + jx + 1
					l[curr[jx]] = true
				}
			}
			left = curr[jx]
		}
		copy(last, curr)
	}
	return len(l)
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(lakes(scanner.Text()))
	}
}
