package main

import (
	"fmt"
	"strings"
)

func do(s string) string {
	var result string

	s = strings.ToLower(s)
	for i := 0; i < len(s); i++ {
		if s[i] == 'y' || s[i] == 'a' || s[i] == 'e' || s[i] == 'i' || s[i] == 'o' || s[i] == 'u' {
			continue
		}
		result += fmt.Sprintf(".%c", s[i])
	}
	return result
}

// func Test(t *testing.T) {
// 	a := assert.New(t)

// 	a.Equal(".t.r", do("tour"))
// 	a.Equal(".c.d.f.r.c.s", do("Codeforces"))
// 	a.Equal(".w.p.w.l", do("wpwl"))
// 	a.Equal(".k.t.j.q.h.p.q.s.v.h.w", do("ktajqhpqsvhw"))
// }

func main() {
	var s string
	fmt.Scan(&s)
	fmt.Println(do(s))
}
