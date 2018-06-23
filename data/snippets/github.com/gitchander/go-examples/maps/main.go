package main

import (
	"fmt"
	"sort"
)

// Source: https://www.callicoder.com/golang-maps/

func main() {
	exampleNil()
	exampleMake()
	exampleInit()
	exampleGet()
	exampleGetAndCheck()
	exampleCheck()
	exampleDelete()
	exampleReference()
	exampleIterating()
	exampleSortedPrint()
	exampleEmptyValues()
	exampleArrayKeys()
	exampleStructKeys()
}

func exampleNil() {
	var m map[string]int
	fmt.Println(m)
	if m == nil {
		fmt.Println("m is nil")
	}
	// m["one"] = 1 // panic: assignment to entry in nil map
}

func exampleMake() {
	var m = make(map[string]int)
	fmt.Println(m)
	if m == nil {
		fmt.Println("map is nil")
	} else {
		fmt.Println("map is not nil")
	}
	m["one"] = 1
	m["two"] = 2
	m["three"] = 3
	fmt.Println(m)
}

func exampleInit() {
	var m = map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
	}
	fmt.Println(m)
}

func exampleGet() {
	var m = map[string]int{
		"January":  1,
		"February": 2,
		"March":    3,
	}

	monthId := m["January"]
	fmt.Println("January id =", monthId)
}

func exampleGetAndCheck() {
	var m = map[string]int{
		"January":  1,
		"February": 2,
		"March":    3,
	}

	key := "January"
	month, ok := m[key]
	if ok {
		fmt.Printf("m has month %s: %d\n", key, month)
	} else {
		fmt.Printf("m has not %s\n", key)
	}

	key = "October"
	month, ok = m[key]
	if ok {
		fmt.Printf("m has month %s: %d\n", key, month)
	} else {
		fmt.Printf("m has not %s\n", key)
	}
}

func exampleCheck() {
	var m = map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
	}

	key := "four"
	_, ok := m[key]
	if ok {
		fmt.Println("m has", key)
	} else {
		fmt.Println("m has not", key)
	}

	key = "six"
	_, ok = m[key]
	if ok {
		fmt.Println("m has", key)
	} else {
		fmt.Println("m has not", key)
	}
}

func exampleDelete() {
	var fileExtensions = map[string]string{
		"C++":    ".cpp",
		"Java":   ".java",
		"Kotlin": ".kt",
		"Golang": ".go",
		"Python": ".py",
	}
	fmt.Println(fileExtensions)
	delete(fileExtensions, "C++")
	delete(fileExtensions, "Javascript")
	fmt.Println(fileExtensions)
}

func exampleReference() {
	var m1 = map[string]int{
		"January":  1,
		"February": 2,
		"March":    3,
	}

	var m2 = m1

	fmt.Println("m1 =", m1)
	fmt.Println("m2 =", m2)

	fmt.Println()

	m2["April"] = 4

	fmt.Println("m1 =", m1)
	fmt.Println("m2 =", m2)
}

func exampleIterating() {
	var personAge = map[string]int{
		"Tom":   10,
		"Jerry": 5,
		"Spike": 12,
	}
	for name, age := range personAge {
		fmt.Println(name, age)
	}
}

func exampleSortedPrint() {
	var m = map[int]string{
		0: "Sunday",
		1: "Monday",
		2: "Tuesday",
		3: "Wednesday",
		4: "Thursday",
		5: "Friday",
		6: "Saturday",
	}

	fmt.Println("\nUnsorted days:")
	for day, name := range m {
		fmt.Printf("%d: %s\n", day, name)
	}

	days := make([]int, 0, len(m))
	for day := range m {
		days = append(days, day)
	}
	sort.Ints(days)

	fmt.Println("\nSorted days:")
	for _, day := range days {
		fmt.Printf("%d: %s\n", day, m[day])
	}
}

func exampleEmptyValues() {
	var m map[int]struct{} // declare map
	fmt.Println(m)
	m = make(map[int]struct{}) // make map
	m[100] = struct{}{}
	m[-34] = struct{}{}
	m[0] = struct{}{}
	fmt.Println(m)
}

func exampleArrayKeys() {
	const n = 3
	var m = make(map[[n]int]string)

	// keys
	var (
		alpha = [n]int{0, 1, 2}
		omega = [n]int{1, 2, 3}
		zetta = [n]int{2, 3, 4}
	)

	m[alpha] = "Alpha"
	m[omega] = "Omega"

	fmt.Println(m)

	fmt.Println("alpha:", m[alpha])
	fmt.Println("zetta:", m[zetta])
}

type Size struct {
	Width  int
	Height int
}

func exampleStructKeys() {
	var m = make(map[Size]string)
	m[Size{Width: 100, Height: 200}] = "100x200"
	m[Size{Width: 245, Height: 11}] = "245x11"
	m[Size{Width: 0, Height: 0}] = "Zero"
	fmt.Println(m)
}
