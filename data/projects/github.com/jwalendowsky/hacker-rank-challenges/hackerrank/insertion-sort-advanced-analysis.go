package hackerrank

import (
	"fmt"
	"gosort"
)

// InsertionSortAdvancedAnalysis Runs the InsertionSortAdvancedAnalysis challenge
func InsertionSortAdvancedAnalysis() {
	numberOfTests := readInteger()
	for k := 0; k < numberOfTests; k++ {
		numberOfItems := readInteger()
		items := readIntegerArray(numberOfItems, readInteger)
		counter := 0

		gosort.InsertionSort(gosort.IntsSort(items), func(i, j int) {
			counter++
		})
		fmt.Println(counter)
	}
}
