package main

import (
	"fmt"
	"math/rand"
	"strings"
)

const (
	messg = "012222 1114142503 0313012513 03141418192102 0113 2419182119021713 06131715070119"
	key   = "BHISOECRTMGWYVALUZDNFJKPQX"
)

var (
	alpha        map[rune]rune
	seen         map[rune]bool
	mess, fwords []string
	p            []int
)

func init() {
	fwords = []string{"ALL", "AT", "ATHINAI", "BE", "BEEN", "BITTORRENT", "BLIZZARD", "CAN",
		"CALL", "COME", "COULD", "DID", "DO", "DOWN", "EACH", "FIRST", "FOR", "GET", "GO",
		"HAVE", "HAD", "LINUX", "LOL", "LOLCAT", "MIDNIGHT", "MORE", "KTHXBAI", "NUMBER",
		"PEERS", "PROTOCOL", "ROTFL", "SCNR", "SEEDING", "START", "YOUTUBE"}
}

func wfit(wx int, w, g string) bool {
	for i := 0; i < len(w); i++ {
		if a, f := alpha[rune(w[i])]; f {
			if a != rune(g[i]) {
				return false
			}
		} else if _, f := seen[rune(g[i])]; f {
			return false
		} else {
			alpha[rune(w[i])] = rune(g[i])
			seen[rune(g[i])] = true
		}
	}
	if wx == len(mess)-1 {
		return true
	}
	var f bool
	for i := 0; i < len(mess); i++ {
		if i <= wx {
			continue
		}
		for _, j := range fwords {
			if wstruct(mess[p[i]]) == wstruct(j) && wfit(i, mess[p[i]], j) {
				f = true
				break
			}
		}
	}
	return f
}

func wstruct(a string) string {
	var (
		curr      int
		strd, ret string
	)
	for _, i := range a {
		if ix := strings.IndexRune(strd, rune(i)); ix >= 0 {
			ret += string('A' + ix)
		} else {
			strd += string(i)
			ret += string('A' + curr)
			curr++
		}
	}
	return ret
}

func main() {
	mess = strings.Fields(messg)
	for ix, i := range mess {
		var wrd []string
		for j := 0; j < len(i); j += 2 {
			var w int
			fmt.Sscan(i[j:j+2], &w)
			wrd = append(wrd, string(w+64))
		}
		mess[ix] = strings.Join(wrd, "")
	}

	for {
		var f bool
		p = rand.Perm(len(mess))
		alpha = make(map[rune]rune)
		seen = make(map[rune]bool)
		for i := 0; i < len(mess); i++ {
			f = false
			for _, j := range fwords {
				if wstruct(mess[p[i]]) == wstruct(j) && wfit(i, mess[p[i]], j) {
					f = true
					break
				}
			}
			if !f {
				break
			}
		}
		if f {
			break
		}
	}

	for ix, i := range mess {
		var wrd []string
		for j := 0; j < len(i); j++ {
			if a, f := alpha[rune(i[j])]; f {
				wrd = append(wrd, string(a))
			} else {
				wrd = append(wrd, "?")
			}
		}
		mess[ix] = strings.Join(wrd, "")
	}
	fmt.Println(strings.Join(mess, " "))
}
