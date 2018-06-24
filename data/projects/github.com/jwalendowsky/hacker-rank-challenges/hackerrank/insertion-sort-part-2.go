package hackerrank

// InsertionSortPart2 runs the challenge to use insert sort for a couple of numbers.
func InsertionSortPart2() {
	numberOfItens := readInteger()
	items := readIntegerArray(numberOfItens, readInteger)

	insertionSort(items)
}

func insertionSort(items []int) {
	size := len(items)

	for i := 1; i < size; i++ {
		for j := i; j > 0 && items[j] < items[j-1]; j-- {
			temp := items[j]
			items[j] = items[j-1]
			items[j-1] = temp
		}
		printArray(items)
	}
}
