package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestCash(t *testing.T) {
	for k, v := range map[string]string{
		"15.94;16.00": "NICKEL,PENNY",
		"17;16":       "ERROR",
		"35;35":       "ZERO",
		"12.12;13.25": "ONE,DIME,PENNY,PENNY,PENNY",
		"45;50":       "FIVE",
		"311.09;500":  "ONE HUNDRED,FIFTY,TWENTY,TEN,FIVE,TWO,ONE,HALF DOLLAR,QUARTER,DIME,NICKEL,PENNY"} {
		if r := cash(k); r != v {
			t.Errorf("failed: cash %s is %s, got %s",
				k, v, r)
		}
	}
}

var (
	units []string
	value []int
)

func init() {
	units = []string{"PENNY", "NICKEL", "DIME", "QUARTER", "HALF DOLLAR",
		"ONE", "TWO", "FIVE", "TEN", "TWENTY", "FIFTY", "ONE HUNDRED"}
	value = []int{1, 5, 10, 25, 50, 100,
		200, 500, 1000, 2000, 5000, 10000}
}

func cash(q string) string {
	s := strings.Split(q, ";")
	var p, c, v int
	if strings.Contains(s[0], ".") {
		fmt.Sscanf(s[0], "%d.%d", &p, &v)
	} else {
		fmt.Sscan(s[0], &p)
	}
	p, v = p*100+v, 0
	if strings.Contains(s[1], ".") {
		fmt.Sscanf(s[1], "%d.%d", &c, &v)
	} else {
		fmt.Sscan(s[1], &c)
	}
	c = c*100 + v
	if p > c {
		return "ERROR"
	} else if p == c {
		return "ZERO"
	}
	var r []string
	change := make([]int, 12)
	for i := 11; i >= 0; i-- {
		for c-p >= value[i] {
			change[i]++
			c -= value[i]
		}
	}
	for i := 11; i >= 0; i-- {
		for change[i] > 0 {
			r = append(r, units[i])
			change[i]--
		}
	}
	return strings.Join(r, ",")
}
