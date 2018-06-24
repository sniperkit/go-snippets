package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func slide(a []int) []int {
	lastNum, lastNumID, lastZero := -1, -1, -1
	for i := range a {
		if a[i] == 0 {
			if lastZero == -1 {
				lastZero = i
			}
		} else if lastNum == a[i] {
			a[lastNumID] = 2 * lastNum
			lastNum = -1
			a[i] = 0
			if lastZero == -1 {
				lastZero = i
			}
		} else {
			lastNum = a[i]
			if lastZero == -1 {
				lastNumID = i
			} else {
				lastNumID = lastZero
				a[lastZero] = a[i]
				a[i] = 0
				lastZero++
			}
		}
	}
	return a
}

func twentyFortyEight(n int, p, q string) string {
	t := strings.Split(q, "|")
	m := make([][]int, n)
	for i := 0; i < n; i++ {
		u := strings.Fields(t[i])
		m[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Sscan(u[j], &m[i][j])
		}
	}
	l := make([]int, n)
	for i := 0; i < n; i++ {
		switch p {
		case "LEFT":
			copy(l, m[i])
			l = slide(l)
			copy(m[i], l)
		case "RIGHT":
			for j := 0; j < n; j++ {
				l[j] = m[i][n-1-j]
			}
			l = slide(l)
			for j := 0; j < n; j++ {
				m[i][n-1-j] = l[j]
			}
		case "UP":
			for j := 0; j < n; j++ {
				l[j] = m[j][i]
			}
			l = slide(l)
			for j := 0; j < n; j++ {
				m[j][i] = l[j]
			}
		case "DOWN":
			for j := 0; j < n; j++ {
				l[j] = m[n-1-j][i]
			}
			l = slide(l)
			for j := 0; j < n; j++ {
				m[n-1-j][i] = l[j]
			}
		}
	}
	r := make([]string, n)
	for i := 0; i < n; i++ {
		u := make([]string, n)
		for j := 0; j < n; j++ {
			u[j] = fmt.Sprint(m[i][j])
		}
		r[i] = strings.Join(u, " ")
	}
	return strings.Join(r, "|")
}

func main() {
	var n int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ";")
		fmt.Sscan(s[1], &n)
		fmt.Println(twentyFortyEight(n, s[0], s[2]))
	}
}
