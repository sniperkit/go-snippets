package hackerrank

// QuickSortPart2 resolves the Quick Sort 1 hacker rank challenge.
func QuickSortPart2() {
	numberOfItens := readInteger()
	items := readIntegerArray(numberOfItens, readInteger)

	quickSortInts(items, 0, len(items)-1)
}

func quickSortInts(items []int, min int, max int) {

	// If the min is equal or greater than the max, our sort endend.
	if min >= max {
		return
	}

	// Let's partition the array in a way that the return is the position of an element that will not
	// change and all elements at his left are smaller than it and all at his right are bigger.
	finalPositionIndex := partitionIntsPrinting(items, min, max)
	printArrayWithBoundaries(items, min, max)

	// Now let's use divide and conquer to sort all elements at the left and at the right of the final
	// position index
	quickSortInts(items, min, finalPositionIndex-1)
	quickSortInts(items, finalPositionIndex+1, max)
	// printArrayWithBoundaries(items, min, max)
}

// partitionInts creates a QuickSort partition of an array slice returning and index
// (min <= finalPositionIndex <= max) with the following definitions:
// - This index is the final position of the element on the array.
// - All the elements before the index are smaller/equal to the value of the element.
// - All the elements after the index are grater than the value of the element.
func partitionIntsPrinting(items []int, min int, max int) (finalPosition int) {
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
