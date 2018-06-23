package main

import "fmt"

// convert types take an int and return a string value.
type convert func(int) string

// value implements convert, returning x as string.
func value(x int) string {
	return fmt.Sprintf("%v", x)
}

// quote123 passes 123 to convert func and returns quoted string.
func quote123(fn convert) string {
	return fmt.Sprintf("%q", fn(123))
}

func main() {
	var result string

	result = value(123)
	fmt.Println(result)
	// Output: 123

	result = quote123(value)
	fmt.Println(result)
	// Output: "123"

	result = quote123(func(x int) string { return fmt.Sprintf("%b", x) })
	fmt.Println(result)
	// Output: "1111011"

	foo := func(x int) string { return "foo" }
	result = quote123(foo)
	fmt.Println(result)
	// Output: "foo"

	_ = convert(foo) // confirm foo satisfies convert at runtime

	// fails due to argument type
	// _ = convert(func(x float64) string { return "" })
}

