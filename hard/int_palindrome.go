package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func isPalindrome(a int) bool {
	if a%10 == 0 {
		return a == 0
	}
	var rev int
	for ; a > rev; a /= 10 {
		rev = 10*rev + a%10
		if a == rev {
			return true
		}
	}
	return rev == a
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		var l, r int
		fmt.Sscanf(scanner.Text(), "%d %d", &l, &r)
		var n int
		for i := l; i <= r; i++ {
			prev := -1
			for j := i; j <= r; j++ {
				var p int
				if prev > -1 {
					p = prev
					if isPalindrome(j) {
						p++
					}
				} else {
					p = 0
					for k := i; k <= j; k++ {
						if isPalindrome(k) {
							p++
						}
					}
				}
				if p%2 == 0 {
					n++
				}
				prev = p
			}
		}
		fmt.Println(n)
	}
}
