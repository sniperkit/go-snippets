package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type seg struct {
	s     string
	left  []*seg
	right []*seg
}

func l(a, b string, n int) bool {
	return a[:n-1] == b[len(b)-n+1:]
}

func r(a, b string, n int) bool {
	return a[len(a)-n+1:] == b[:n-1]
}

func allident(direction []*seg) bool {
	for i := 1; i < len(direction); i++ {
		if direction[i].s != direction[0].s {
			return false
		}
	}
	return true
}

func glue(q string) string {
	p := strings.Split(strings.Trim(q, "|"), "|")
	n := len(p[0])
	segs := make([]seg, len(p))
	for ix, i := range p {
		segs[ix] = seg{s: i}
	}

	for ix := 0; ix < len(p)-1; ix++ {
		for jx := ix + 1; jx < len(p); jx++ {
			if l(segs[ix].s, segs[jx].s, n) {
				segs[ix].left = append(segs[ix].left, &segs[jx])
				segs[jx].right = append(segs[jx].right, &segs[ix])
			}
			if r(segs[ix].s, segs[jx].s, n) {
				segs[ix].right = append(segs[ix].right, &segs[jx])
				segs[jx].left = append(segs[jx].left, &segs[ix])
			}
		}
	}

	for ix := 0; ix < len(p); ix++ {
		f := true
		for f {
			f = false
			if len(segs[ix].left) == 1 {
				f = true
				jx := segs[ix].left[0]
				if jx.s != "" {
					s := jx.s
					segs[ix].s = s[:len(s)-n+1] + segs[ix].s
					segs[ix].left, *segs[ix].left[0] = segs[ix].left[0].left, seg{"", nil, nil}
					for jx := 0; jx < len(segs[ix].left); jx++ {
						for kx := 0; kx < len(segs[ix].left[jx].right); kx++ {
							if segs[ix].left[jx].right[kx].s == "" {
								segs[ix].left[jx].right[kx] = &segs[ix]
							}
						}
					}
					for jx := len(segs[ix].left) - 1; jx >= 0; jx-- {
						if segs[ix].left[jx] == &segs[ix] {
							segs[ix].left = append(segs[ix].left[:jx], segs[ix].left[jx+1:]...)
						}
					}
					for jx := len(segs[ix].right) - 1; jx >= 0; jx-- {
						if segs[ix].right[jx] == &segs[ix] {
							segs[ix].right = append(segs[ix].right[:jx], segs[ix].right[jx+1:]...)
						}
					}
				} else {
					segs[ix].left = []*seg{}
				}
			}
			if len(segs[ix].right) == 1 {
				f = true
				jx := segs[ix].right[0]
				if jx.s != "" {
					segs[ix].s = segs[ix].s + jx.s[n-1:]
					segs[ix].right, *segs[ix].right[0] = segs[ix].right[0].right, seg{"", nil, nil}
					for jx := 0; jx < len(segs[ix].right); jx++ {
						for kx := 0; kx < len(segs[ix].right[jx].left); kx++ {
							if segs[ix].right[jx].left[kx].s == "" {
								segs[ix].right[jx].left[kx] = &segs[ix]
							}
						}
					}
					for jx := len(segs[ix].left) - 1; jx >= 0; jx-- {
						if segs[ix].left[jx] == &segs[ix] {
							segs[ix].left = append(segs[ix].left[:jx], segs[ix].left[jx+1:]...)
						}
					}
					for kx := len(segs[ix].right) - 1; kx >= 0; kx-- {
						if segs[ix].right[kx] == &segs[ix] {
							segs[ix].right = append(segs[ix].right[:kx], segs[ix].right[kx+1:]...)
						}
					}
				} else {
					segs[ix].right = []*seg{}
				}
			}
			if len(segs[ix].left) > 1 && allident(segs[ix].left) {
				f = true
				// take first, remove it from others
				segs[ix].left = segs[ix].left[:1]
				for jx, j := range segs[ix].left[0].right {
					if j != &segs[ix] {
						for kx, k := range j.left {
							if k == segs[ix].left[0] {
								segs[ix].left[0].right[jx].left, segs[ix].left[0].right[jx].left[kx] = segs[ix].left[0].right[jx].left[:len(segs[ix].left[0].right[jx].left)-1], segs[ix].left[0].right[jx].left[len(segs[ix].left[0].right[jx].left)-1]
								break
							}
						}
					}
				}
			}
			if len(segs[ix].right) > 1 && allident(segs[ix].right) {
				f = true
				// take first, remove it from others
				segs[ix].right = segs[ix].right[:1]
				for jx, j := range segs[ix].right[0].left {
					if j != &segs[ix] {
						for kx, k := range j.right {
							if k == segs[ix].right[0] {
								segs[ix].right[0].left[jx].right, segs[ix].right[0].left[jx].right[kx] = segs[ix].right[0].left[jx].right[:len(segs[ix].right[0].left[jx].right)-1], segs[ix].right[0].left[jx].right[len(segs[ix].right[0].left[jx].right)-1]
								break
							}
						}
					}
				}

			}
		}
	}

	p = []string{}
	for _, i := range segs {
		if i.s != "" {
			p = append(p, i.s)
		}
	}
	return strings.Join(p, "|")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(glue(strings.Trim(scanner.Text(), "|")))
	}
}
