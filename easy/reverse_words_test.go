package main

import (
	"strings"
	"testing"
)

func TestReverse(t *testing.T) {
	for k, v := range map[string]string{
		"Hello World":              "World Hello",
		"Hello CodeEval":           "CodeEval Hello",
		"can you talk like Yoda ?": "? Yoda like talk you can"} {
		if r := reverse(k); r != v {
			t.Errorf("failed: reverse %s is %s, got %s",
				k, v, r)
		}
	}
}

func TestReverseS(t *testing.T) {
	for k, v := range map[string]string{
		"Hello World":              "World Hello",
		"Hello CodeEval":           "CodeEval Hello",
		"can you talk like Yoda ?": "? Yoda like talk you can"} {
		if r := reverseS(k); r != v {
			t.Errorf("failed: reverseS %s is %s, got %s",
				k, v, r)
		}
	}
}

func BenchmarkReverse(b *testing.B) {
	h := []byte{32, 97, 98, 99}
	for i := 0; i < b.N; i++ {
		var s []byte
		for j := i; j > 0; j >>= 2 {
			s = append(s, h[i%4])
		}
		reverse(string(s))
	}
}

func BenchmarkReverseS(b *testing.B) {
	h := []byte{32, 97, 98, 99}
	for i := 0; i < b.N; i++ {
		var s []byte
		for j := i; j > 0; j >>= 2 {
			s = append(s, h[i%4])
		}
		reverseS(string(s))
	}
}

func reverse(s string) string {
	var (
		c []byte
		t []string
	)
	for _, i := range s {
		if i == ' ' {
			t, c = append([]string{string(c)}, t...), []byte{}
		} else {
			c = append(c, byte(i))
		}
	}
	t = append([]string{string(c)}, t...)
	return strings.Join(t, " ")
}

func reverseS(s string) string {
	var (
		c string
		t []string
	)
	for _, i := range s {
		if i == ' ' {
			t, c = append([]string{c}, t...), ""
		} else {
			c += string(i)
		}
	}
	t = append([]string{c}, t...)
	return strings.Join(t, " ")
}
