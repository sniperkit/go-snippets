package main

import "fmt"

func main() {
	str1 := "Write a string " +
		"over multiple lines " +
		"in the source code ..."

	fmt.Println(str1)

	str2 := `Write and print
  this string with linebreaks ...
  `

	fmt.Println(str2)
}
