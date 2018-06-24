package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestAbsurd(t *testing.T) {
	if r := absurd(5, strings.Split("0,1,2,3,0", ",")); r != 0 {
		t.Errorf("failed: absurd 5 0,1,2,3,0 is 0, got %d",
			r)
	}
	if r := absurd(20, strings.Split("0,1,10,3,2,4,5,7,6,8,11,9,15,12,13,4,16,18,17,14", ",")); r != 4 {
		t.Errorf("failed: absurd 20 0,1,10,3,2,4,5,7,6,8,11,9,15,12,13,4,16,18,17,14 is 0, got %d",
			r)
	}
	if r := absurd(6, strings.Split("4,0,1,2,3,0", ",")); r != 0 {
		t.Errorf("failed: absurd 6 4,0,1,2,3,0 is 0, got %d",
			r)
	}
	if r := absurd(21, strings.Split("19,0,1,10,3,2,4,5,7,6,8,11,9,15,12,13,4,16,18,17,14", ",")); r != 4 {
		t.Errorf("failed: absurd 21 19,0,1,10,3,2,4,5,7,6,8,11,9,15,12,13,4,16,18,17,14 is 0, got %d",
			r)
	}
}

func TestAbsurdMap(t *testing.T) {
	if r := absurdMap(5, strings.Split("0,1,2,3,0", ",")); r != 0 {
		t.Errorf("failed: absurdMap 5 0,1,2,3,0 is 0, got %d",
			r)
	}
	if r := absurdMap(20, strings.Split("0,1,10,3,2,4,5,7,6,8,11,9,15,12,13,4,16,18,17,14", ",")); r != 4 {
		t.Errorf("failed: absurdMap 20 0,1,10,3,2,4,5,7,6,8,11,9,15,12,13,4,16,18,17,14 is 0, got %d",
			r)
	}
	if r := absurdMap(6, strings.Split("4,0,1,2,3,0", ",")); r != 0 {
		t.Errorf("failed: absurdMap 6 4,0,1,2,3,0 is 0, got %d",
			r)
	}
	if r := absurdMap(21, strings.Split("19,0,1,10,3,2,4,5,7,6,8,11,9,15,12,13,4,16,18,17,14", ",")); r != 4 {
		t.Errorf("failed: absurdMap 21 19,0,1,10,3,2,4,5,7,6,8,11,9,15,12,13,4,16,18,17,14 is 0, got %d",
			r)
	}
}

func BenchmarkAbsurd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := make([]string, i%1024+1)
		for j := 0; j < i%1024+1; j++ {
			f[j] = fmt.Sprint(j)
		}
		absurd(i%1024+2, f)
	}
}

func BenchmarkAbsurdMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f := make([]string, i%1024+1)
		for j := 0; j < i%1024+1; j++ {
			f[j] = fmt.Sprint(j)
		}
		absurdMap(i%1024+2, f)
	}
}

func absurd(n int, t []string) int {
	var v int
	m := make([]bool, n-1)
	for _, i := range t {
		fmt.Sscanf(i, "%d", &v)
		if m[v] {
			return v
		}
		m[v] = true
	}
	return -1
}

func absurdMap(n int, t []string) int {
	var v int
	m := make(map[int]bool, n-1)
	for _, i := range t {
		fmt.Sscanf(i, "%d", &v)
		if m[v] {
			return v
		}
		m[v] = true
	}
	return -1
}
