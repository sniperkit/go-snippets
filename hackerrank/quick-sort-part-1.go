package hackerrank

// QuickSortPart1 resolves the Quick Sort 1 hacker rank challenge.
func QuickSortPart1() {
	numberOfItens := readInteger()
	items := readIntegerArray(numberOfItens, readInteger)

	partitionInts(items, 0, len(items)-1)

	printArray(items)
}

// partitionInts creates a QuickSort partition of an array slice returning and index
// (min <= finalPositionIndex <= max) with the following definitions:
// - This index is the final position of the element on the array.
// - All the elements before the index are smaller/equal to the value of the element.
// - All the elements after the index are grater than the value of the element.
func partitionInts(items []int, min int, max int) (finalPosition int) {
	finalPosition = min  // Let's start assuming that the final position is the minimum.
	rightPosition := max // Let's keep track of the right position of the array

	// Let's walk from both sides of the array until we reach the middle of it.
	for finalPosition < rightPosition {

		// If the value of the element on the final position is greater than the value of the element on
		// finalPosition + 1, let's swap than and increment the final position.
		if items[finalPosition] >= items[finalPosition+1] {
			items[finalPosition], items[finalPosition+1] = items[finalPosition+1], items[finalPosition]
			finalPosition++
		} else if items[finalPosition] >= items[rightPosition] {

			// Otherwise, if the element on the final position is greater than the element on the right position,
			// let's swipe the element on the right position with the element of the finalPosition + 1 (since
			// we know such element is bigger than the element on the final position. After that, let's swap
			// the value of the final position with the element on final position + 1, updating the final position.
			// items[finalPosition+1], items[rightPosition] = items[rightPosition], items[finalPosition+1]
			items[finalPosition+1], items[rightPosition] = items[rightPosition], items[finalPosition+1]
			items[finalPosition], items[finalPosition+1] = items[finalPosition+1], items[finalPosition]
			finalPosition++
		} else {

			// Otherwise, let's decrement the right position
			rightPosition--
		}
	}

	return
}

/*
// partitionInts creates a QuickSort partition of an array slice returning and index
// (min <= finalPositionIndex <= max) with the following definitions:
// - This index is the final position of the element on the array.
// - All the elements before the index are smaller/equal to the value of the element.
// - All the elements after the index are grater than the value of the element.
func partitionInts(items []int, min int, max int) (finalPosition int) {
	finalPosition = min
	i := min + 1
	j := max

	for true {
		for items[i] < items[min] {
			i++
			if i == max {
				break
			}
		}

		for items[min] < items[j] {
			j--
			if j == min {
				break
			}
		}

		if i >= j {
			break
		}
		items[i], items[j] = items[j], items[i]
	}

	items[min], items[j] = items[j], items[min]

	finalPosition = j
	return
}s
*/
