package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ";")
		var n, m int
		fmt.Sscanf(s[0], "%d,%d", &m, &n)
		prev, curr := make([]int, n), make([]int, n)
		for i := 0; i < n; i++ {
			if s[1][i] == '*' {
				prev[i] = -1
			} else if i > 0 && prev[i-1] == -1 {
				prev[i] = 1
			}
		}
		for i := 1; i < m; i++ {
			if s[1][i*n] == '*' {
				curr[0] = -1
			} else {
				if prev[0] == -1 {
					curr[0] = 1
				} else {
					curr[0] = 0
				}
				if n > 1 && prev[1] == -1 {
					curr[0]++
				}
			}
			for j := 1; j < n; j++ {
				if s[1][i*n+j] == '*' {
					curr[j] = -1
				} else {
					if curr[j-1] == -1 {
						curr[j] = 1
					} else {
						curr[j] = 0
					}
					if prev[j-1] == -1 {
						curr[j]++
					}
					if prev[j] == -1 {
						curr[j]++
					}
					if n > j+1 && prev[j+1] == -1 {
						curr[j]++
					}
				}
				if prev[j-1] == -1 {
					fmt.Print("*")
				} else {
					if curr[j-1] == -1 {
						prev[j-1]++
					}
					if curr[j] == -1 {
						prev[j-1]++
					}
					if prev[j] == -1 {
						prev[j-1]++
					}
					if j > 1 && curr[j-2] == -1 {
						prev[j-1]++
					}
					fmt.Print(prev[j-1])
				}
			}
			if prev[n-1] == -1 {
				fmt.Print("*")
			} else {
				if curr[n-1] == -1 {
					prev[n-1]++
				}
				if n > 1 && curr[n-2] == -1 {
					prev[n-1]++
				}
				fmt.Print(prev[n-1])
			}
			copy(prev, curr)
		}
		for j := 0; j < n; j++ {
			if prev[j] == -1 {
				fmt.Print("*")
			} else {
				if j < n-1 && prev[j+1] == -1 {
					prev[j]++
				}
				fmt.Print(prev[j])
			}
		}
		fmt.Println()
	}
}
