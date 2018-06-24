package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type state struct {
	knownOut, knownIn     uint16
	unknownOut, unknownIn int8
}

func kernighan(a uint16) (r int8) {
	for ; a > 0; a &= a - 1 {
		r++
	}
	return r
}

func all(a uint16) (r []uint) {
	for i := uint(0); a > 0; a, i = a>>1, i+1 {
		if a&1 > 0 {
			r = append(r, i)
		}
	}
	return r
}

func removeIn(ss []state, x uint) (r []state) { // ex
	for ix, i := range ss {
		if i.knownIn&(1<<x) == 0 {
			if i.knownOut&(1<<x) > 0 {
				ss[ix].knownOut -= (1 << x)
			} else if i.unknownOut > 0 {
				r = append(r, state{i.knownOut, i.knownIn + 1<<x, i.unknownOut - 1, i.unknownIn})
			}
			r = append(r, state{ss[ix].knownOut, i.knownIn + 1<<x, i.unknownOut, i.unknownIn})
		}
	}
	return r
}
func removeOut(ss []state, x uint) (r []state) { // lx
	for ix, i := range ss {
		if i.knownOut&(1<<x) == 0 {
			if i.knownIn&(1<<x) > 0 {
				ss[ix].knownIn -= 1 << x
			} else if i.unknownIn > 0 {
				r = append(r, state{i.knownOut + 1<<x, i.knownIn, i.unknownOut, i.unknownIn - 1})
			}
			r = append(r, state{i.knownOut + 1<<x, ss[ix].knownIn, i.unknownOut, i.unknownIn})
		}
	}
	return r
}
func createIn(ss []state) (r []state) { // e0
	for _, i := range ss {
		kno := all(i.knownOut)
		for _, j := range kno {
			r = append(r, state{i.knownOut - 1<<j, i.knownIn + 1<<j, i.unknownOut, i.unknownIn})
		}
		if i.unknownOut > 0 {
			r = append(r, state{i.knownOut, i.knownIn, i.unknownOut - 1, i.unknownIn + 1})
		}
		r = append(r, state{i.knownOut, i.knownIn, i.unknownOut, i.unknownIn + 1})
	}
	return r
}
func createOut(ss []state) (r []state) { // l0
	for _, i := range ss {
		kno := all(i.knownIn)
		for _, j := range kno {
			r = append(r, state{i.knownOut + 1<<j, i.knownIn - 1<<j, i.unknownOut, i.unknownIn})
		}
		if i.unknownIn > 0 {
			r = append(r, state{i.knownOut, i.knownIn, i.unknownOut + 1, i.unknownIn - 1})
		}
		r = append(r, state{i.knownOut, i.knownIn, i.unknownOut + 1, i.unknownIn})
	}
	return r
}

func main() {
	var ix uint
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		u := make(map[uint]int8)
		ux := int8(1)
		s := strings.Split(scanner.Text(), "; ")
		t := strings.Split(s[1], "|")
		states := []state{state{0, 0, 0, 0}}
		for _, i := range t {
			e := i[0] == 'E'
			fmt.Sscanf(i[2:], "%d", &ix)
			var cx int8
			if ix > 0 {
				if x, f := u[ix]; f {
					cx = x
				} else {
					u[ix] = ux
					cx = ux
					ux++
				}
			}
			if e {
				if cx > 0 {
					states = removeIn(states, uint(cx))
				} else {
					states = createIn(states)
				}
			} else {
				if cx > 0 {
					states = removeOut(states, uint(cx))
				} else {
					states = createOut(states)
				}
			}
			for j := len(states) - 1; j > 0; j-- {
				for k := j - 1; k >= 0; k-- {
					if states[j].knownIn == states[k].knownIn && states[j].knownOut == states[k].knownOut && states[j].unknownIn == states[k].unknownIn {
						states[j], states = states[len(states)-1], states[:len(states)-1]
						break
					}
				}
			}
			if len(states) == 0 {
				break
			}
		}
		if len(states) == 0 {
			fmt.Println("CRIME TIME")
		} else {
			minIn := int8(15)
			for _, i := range states {
				in := kernighan(i.knownIn) + i.unknownIn
				if in < minIn {
					minIn = in
				}
			}
			fmt.Println(minIn)
		}
	}
}
