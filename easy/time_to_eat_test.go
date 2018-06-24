package main

import (
	"sort"
	"strings"
	"testing"
)

func TestTimeToEat(t *testing.T) {
	for k, v := range map[string]string{
		"02:26:31 14:44:45 09:53:27": "14:44:45 09:53:27 02:26:31",
		"05:33:44 21:25:41":          "21:25:41 05:33:44"} {
		if r := timeToEat(k); r != v {
			t.Errorf("failed: timeToEat %s is %s, got %s",
				k, v, r)
		}
	}
}

func timeToEat(q string) string {
	s := strings.Fields(q)
	sort.Sort(sort.Reverse(sort.StringSlice(s)))
	return strings.Join(s, " ")
}
