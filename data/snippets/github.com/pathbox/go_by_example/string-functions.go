package main

import (
	"fmt"
	ss "strings"
)

var p = fmt.Println

func main() {
	p("Contains: ", ss.Contains("test", "es"))
	p("Count: ", ss.Count("test", "t"))
	p("HasPrefix: ", ss.HasPrefix("test", "te"))
	p("Index: ", ss.HasSuffix("test", "st"))
	p("Join: ", ss.Join([]string{"a", "b"}, "-"))
	p("Repeat:    ", ss.Repeat("a", 5))
	p("Replace:   ", ss.Replace("foo", "o", "0", -1))
	p("Replace:   ", ss.Replace("foo", "o", "0", 1))
	p("Split:     ", ss.Split("a-b-c-d-e", "-"))
	p("ToLower:   ", ss.ToLower("TEST"))
	p("ToUpper:   ", ss.ToUpper("test"))
	p()
	p("Len: ", len("hello"))
	p("Char:", "hello"[1])
}
