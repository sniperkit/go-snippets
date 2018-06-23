package slices

import (
	"fmt"
	"testing"
)

func TestConvert(t *testing.T) {
	var intConverter func(...int) []interface{}
	Convert(&intConverter)
	intConverter(1, 2, 3, 4)

	type T1 struct{ Name string }
	var t1Converter func(...T1) []interface{}
	Convert(&t1Converter)
	t1Converter([]T1{{"A"}, {"B"}}...)

	var stringConverter func([]string, []string) ([]interface{}, []interface{})
	Convert(&stringConverter)
	stringConverter([]string{"a", "b"}, []string{"x", "y"})

	type T2 int
	var t2Converter func(T2, ...T2) (T2, []interface{})
	Convert(&t2Converter)
	t2Converter(T2(1), T2(2), T2(3))
}

func ExampleConvert() {
	var intConverter func(...int) []interface{}
	Convert(&intConverter)
	fmt.Printf("%v\n%v\n%v\n", intConverter(1, 2, 3)...)
	// Output: 1
	// 2
	// 3
}
