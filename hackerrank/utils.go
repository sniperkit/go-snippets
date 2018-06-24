package hackerrank

import (
	"fmt"
	"strconv"
)

// Reads an integer array using the mapper function to define the values of the array itens.
func readIntegerArray(arrayLength int, mapper func() int) (integerArray []int) {
	integerArray = make([]int, arrayLength) //  Let's allocate an "array" (slice) of arrayLength itens.

	for i := 0; i < arrayLength; i++ {
		integerArray[i] = mapper() // For each index of the array we call the mapper function
	}

	return
}

// Reads an integer from the stdin
func readInteger() (value int) {
	fmt.Scanf("%d", &value) // We can use fmt.Scan to read a value from stdin.
	return
}

// Gets the maximum value of two number
func max(firstValue int, secondValue int) int {
	// NOTE: golang does not have a built-in int Math.Max function
	if firstValue > secondValue {
		return firstValue
	}
	return secondValue
}

func printArray(items []int) {

	// Since Go doesn't have a join operation to arrays of types other the strings, we'll print by hand.

	printArrayWithBoundaries(items, 0, len(items)-1)
}

func printArrayWithBoundaries(items []int, min int, max int) {
	// Since Go doesn't have a join operation to arrays of types other the strings, we'll print by hand.

	result := ""
	counter := 0

	for i := min; i <= max; i++ {
		if counter > 0 {
			result += " "
		}
		counter++
		result += strconv.Itoa(items[i])
	}

	fmt.Println(result)
}
