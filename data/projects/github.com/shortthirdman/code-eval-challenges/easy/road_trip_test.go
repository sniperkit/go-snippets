package main

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

func TestRoadTrip(t *testing.T) {
	for k, v := range map[string]string{
		"Rkbs,5453; Wdqiz,1245; Rwds,3890; Ujma,5589; Tbzmo,1303;":      "1245,58,2587,1563,136",
		"Vgdfz,70; Mgknxpi,3958; Nsptghk,2626; Wuzp,2559; Jcdwi,3761;":  "70,2489,67,1135,197",
		"Yvnzjwk,5363; Pkabj,5999; Xznvb,3584; Jfksvx,1240; Inwm,5720;": "1240,2344,1779,357,279",
		"Ramytdb,2683; Voclqmb,5236;":                                   "2683,2553"} {
		if r := roadTrip(k); r != v {
			t.Errorf("failed: roadTrip %s is %s, got %s",
				k, v, r)
		}
	}
}

func roadTrip(s string) string {
	var v, c int
	t := strings.Split(s, ";")
	l, m := make([]int, len(t)-1), make([]string, len(t)-1)
	for i := range l {
		u := strings.Split(t[i], ",")
		fmt.Sscan(u[1], &v)
		l[i] = v
	}
	sort.Ints(l)
	for i := range l {
		m[i], c = fmt.Sprint(l[i]-c), l[i]
	}
	return strings.Join(m, ",")
}
