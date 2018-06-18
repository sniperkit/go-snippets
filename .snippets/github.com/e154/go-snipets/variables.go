package main

import "fmt"

var i int
var c, python, java = true, false, "no!"
var (
	array = [...]string{
		"a",
		"b",
	}
	array2 = [2]string{
		"a",
		"b",
	}
	array3 = []string{
		"a",
		"b",
	}
)

func main() {
	k := 3
	fmt.Println(i, c, python, java, k)
}
