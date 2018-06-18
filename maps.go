package main

import "fmt"

func main() {
	m := make(map[string]int)

	m["key1"] = 1
	m["key2"] = 2

	for k, v := range m {
		fmt.Println(k, v)
	}
	fmt.Println(m)
	fmt.Println("val", m["key1"])

	fmt.Println(len(m))

	delete(m, "key2")

	fmt.Println(m)

	_, prs := m["key2"]
	fmt.Println(prs)
}
