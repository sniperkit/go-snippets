package main

import (
	"strings"
	"testing"
)

type tuple struct {
	n int
	q string
}

func TestSlist(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{1, "aa"}:    "a",
		tuple{2, "ab"}:    "aa,ab,ba,bb",
		tuple{3, "pop"}:   "ooo,oop,opo,opp,poo,pop,ppo,ppp",
		tuple{4, "caff"}:  "aaaa,aaac,aaaf,aaca,aacc,aacf,aafa,aafc,aaff,acaa,acac,acaf,acca,accc,accf,acfa,acfc,acff,afaa,afac,afaf,afca,afcc,afcf,affa,affc,afff,caaa,caac,caaf,caca,cacc,cacf,cafa,cafc,caff,ccaa,ccac,ccaf,ccca,cccc,cccf,ccfa,ccfc,ccff,cfaa,cfac,cfaf,cfca,cfcc,cfcf,cffa,cffc,cfff,faaa,faac,faaf,faca,facc,facf,fafa,fafc,faff,fcaa,fcac,fcaf,fcca,fccc,fccf,fcfa,fcfc,fcff,ffaa,ffac,ffaf,ffca,ffcc,ffcf,fffa,fffc,ffff",
		tuple{2, "afdqo"}: "aa,ad,af,ao,aq,da,dd,df,do,dq,fa,fd,ff,fo,fq,oa,od,of,oo,oq,qa,qd,qf,qo,qq"} {
		if r := slist(sortu(k.q), "", k.n); r != v {
			t.Errorf("failed: slist %d %s is %s, got %s",
				k.n, k.q, v, r)
		}
	}
}

func sortu(s string) (r string) {
	for _, i := range s {
		var done bool
		for jx, j := range r {
			if j == i {
				done = true
				break
			} else if j > i {
				if jx > 0 {
					r = r[0:jx] + string(i) + r[jx:len(r)]
				} else {
					r = string(i) + r
				}
				done = true
				break
			}
		}
		if !done {
			r = r + string(i)
		}
	}
	return r
}

func slist(s, t string, n int) string {
	if n == 0 {
		return t
	}
	var r []string
	for _, i := range s {
		r = append(r, slist(s, t+string(i), n-1))
	}
	return strings.Join(r, ",")
}
