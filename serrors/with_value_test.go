package serrors_test

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-zero-boilerplate/facades/serrors"
)

func ExampleWithValue1() {
	cause := errors.New("whoops")
	err := serrors.WithValue(cause, "key1", "val1")
	fmt.Println(err)

	// Output: [key1=val1] whoops
}

func ExampleWithValue2() {
	cause := errors.New("whoops")
	err := serrors.WithValue(serrors.WithValue(cause, "key1", "val1"), "key2", "val2")
	fmt.Println(err)

	// Output: [key2=val2] [key1=val1] whoops
}

func ExampleWithValue3() {
	cause := errors.New("whoops")
	err := serrors.WithValue(serrors.WithValue(serrors.WithValue(cause, "key1", "val1"), "key2", "val2"), "key1", "val1b")

	fmt.Println(err)

	// Output: [key1=val1b] [key2=val2] [key1=val1] whoops
}

func ExampleUniqueValues() {
	cause := errors.New("whoops")
	err := serrors.WithValue(serrors.WithValue(serrors.WithValue(cause, "key1", "val1"), "key2", "val2"), "key1", "val1b")

	orderedKeys, valueMap := serrors.UniqueValues(err)

	valStrs := []string{}
	for _, key := range orderedKeys {
		valStrs = append(valStrs, fmt.Sprintf("[%s=%+v] ", key, valueMap[key]))
	}

	fmt.Println(strings.Join(valStrs, ""))

	// Output: [key1=val1] [key2=val2]
}
