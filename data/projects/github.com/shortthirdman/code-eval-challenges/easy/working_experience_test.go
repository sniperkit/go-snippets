package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestMonth(t *testing.T) {
	for k, v := range map[string]int{
		"Feb 2004": 169,
		"Jul 2008": 222,
		"Aug 2011": 259,
		"Sep 2013": 284} {
		if r := month(k); r != v {
			t.Errorf("failed: month %s is %d, got %d",
				k, v, r)
		}
	}
}

func TestWorking(t *testing.T) {
	for k, v := range map[string]int{
		"Feb 2004-Dec 2009; Sep 2004-Jul 2008":                                                          5,
		"Aug 2013-Mar 2014; Apr 2013-Aug 2013; Jun 2014-Aug 2015; Apr 2003-Nov 2004; Apr 2014-Jan 2015": 4,
		"Mar 2003-Jul 2003; Nov 2003-Jan 2004; Apr 1999-Nov 1999":                                       1,
		"Apr 1992-Dec 1993; Feb 1996-Sep 1997; Jan 2002-Jun 2002; Sep 2003-Apr 2004; Feb 2010-Nov 2011": 6,
		"Feb 2004-May 2004; Jun 2004-Jul 2004":                                                          0} {
		if r := working(k); r != v {
			t.Errorf("failed: working %s is %d, got %d",
				k, v, r)
		}
	}
}

func BenchmarkMonth(b *testing.B) {
	mnth := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul",
		"Aug", "Sep", "Oct", "Nov", "Dec"}
	for i := 0; i < b.N; i++ {
		month(mnth[i%12] + " " + fmt.Sprint((i/12)%27+1990))
	}
}

var moy map[string]int

func init() {
	moy = map[string]int{"Jan": 0, "Feb": 1, "Mar": 2, "Apr": 3,
		"May": 4, "Jun": 5, "Jul": 6, "Aug": 7,
		"Sep": 8, "Oct": 9, "Nov": 10, "Dec": 11}
}

func month(s string) int {
	var k int
	fmt.Sscan(s[4:8], &k)
	return 12*(k-1990) + moy[s[0:3]]
}

func working(q string) int {
	work := make(map[int]bool)
	t := strings.Split(q, "; ")
	for _, i := range t {
		t0, t1 := month(i[0:8]), month(i[9:17])
		for j := t0; j <= t1; j++ {
			work[j] = true
		}
	}
	return len(work) / 12
}
