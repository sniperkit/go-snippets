package main

import (
	"fmt"
)

type testType []string

func test(t interface{}) {

	switch v := t.(type) {
	case testType:
		fmt.Println("testType", v)
	}
}

func main() {

	t := testType{}
	t = append(t, "q", "w")
	test(t)
}

