package gensym_test

import (
	"testing"
	. "github.com/gophergala/gomp/gensym"
)

func check(t *testing.T, expected, actual string) {
	if actual != expected {
		t.Errorf("Fail: %s != %s\n", expected, actual)
	}
}

func TestMkGen(t *testing.T) {
	f := MkGen("")
	if VocabOn {
		for i := 0; i < 108; i++ {
			f()
		}
	}
	check(t, "gompsym0", f())
	check(t, "gompsym1", f())
}

func TestTwoFuncs(t *testing.T) {
	const n = 500
	f1 := MkGen("")
	f2 := MkGen("")
	s1 := make([]string, n)
	s2 := make([]string, n)
	for i := 0; i < n; i++ {
		s1[i] = f1()
	}
	for i := 0; i < n; i++ {
		s2[i] = f2()
	}
	for i := 0; i < n; i++ {
		if s1[i] != s2[i] {
			t.Fatal("Two invocations of MkGen gave different results")
		}
	}
}

func TestTokens(t *testing.T) {
	var src = `package main

import (
	"fmt"
)

func main() {
	var gompsym0 int32 = 5
	gompsym1 := "asd"
	fmt.Printf("Hello, %d, %s and gompsym2\n", gompsym0, gompsym1)
}
`

	f := MkGen(src)
	if VocabOn {
		for i := 0; i < 108; i++ {
			f()
		}
	}
	check(t, "gompsym2", f())
}
