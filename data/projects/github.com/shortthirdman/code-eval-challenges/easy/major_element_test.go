package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestMajor(t *testing.T) {
	for k, v := range map[string]string{
		"92,19,19,76,19,21,19,85,19,19,19,94,19,19,22,67,83,19,19,54,59,1,19,19":              "19",
		"92,11,30,92,1,11,92,38,92,92,43,92,92,51,92,36,97,92,92,92,43,22,84,92,92":           "92",
		"4,79,89,98,48,42,39,79,55,70,21,39,98,16,96,2,10,24,14,47,0,50,95,20,95,48,50,12,42": "None",
		"1,2,3,4,4,4":           "None",
		"4,4,4,1,2,3":           "None",
		"1,2,3,100,100,100,100": "100"} {
		if r := major(k); r != v {
			t.Errorf("failed: major %s is %s, got %s",
				k, v, r)
		}
	}
}

func BenchmarkMajor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var s []string
		for j := i + 1; j > 0; j /= 101 {
			s = append(s, fmt.Sprint(j%101))
		}
		major(strings.Join(s, ","))
	}
}

func major(s string) string {
	var k, maxnum, maxcount int
	num := make([]int, 101)
	t := strings.Split(s, ",")
	for ix, i := range t {
		fmt.Sscan(i, &k)
		num[k]++
		if num[k] > maxcount {
			maxnum, maxcount = k, num[k]
		}
		if len(t)/2+1 > maxcount+len(t)-ix-1 {
			break
		} else if maxcount >= len(t)/2+1 {
			return fmt.Sprint(maxnum)
		}
	}
	return "None"
}
