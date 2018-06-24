package hackerrank

import "fmt"

type prefixSum struct {
	value         int
	originalIndex int
}

// MaximizeSum is the function that solves the maximize sum algorithm.
func MaximizeSum() {
	numberOfTests := readInteger()

	for i := 0; i < numberOfTests; i++ {
		numberOfItems := readInteger()
		modulo := readInteger()

		// let's find the maximum modulo sum for numberOfItems itens
		fmt.Println(findMaximumModuloSum(numberOfItems, modulo))
	}
}

func findMaximumModuloSum(numberOfItems int, modulo int) (maxSum int) {
	currentPrefixSum := 0

	// Step 1: O(n)
	// Let's create an array of preffix sum items to track their preffix sum and their original indexes.

	prefixSums := make([]prefixSum, numberOfItems)

	for i := 0; i < numberOfItems; i++ {
		currentPrefixSum = (readInteger()%modulo + currentPrefixSum) % modulo
		prefixSums[i] = prefixSum{currentPrefixSum, i}
	}

	// Step 2: O(n2)
	// Let's sort the items and, for each exchange, check if the preffix sum is the max.
	sortItems(prefixSums, func(i int, j int, iPrefixSum prefixSum, jPrefixSum prefixSum) {
		if iPrefixSum.originalIndex < jPrefixSum.originalIndex {
			maxSum = max(maxSum, (jPrefixSum.value-iPrefixSum.value+modulo)%modulo)
		}
	})

	maxSum = max(maxSum, prefixSums[numberOfItems-1].value)

	return
}

func sortItems(items []prefixSum, listener func(i int, j int, iValue prefixSum, jValue prefixSum)) {
	// shellSort(items, listener)
	mergeSort(items, listener)
}

func shellSort(items []prefixSum, listener func(i int, j int, iValue prefixSum, jValue prefixSum)) {
	size := len(items)
	h := 1

	for h < size/3 {
		h = 3*h + 1
	}

	for h >= 1 {
		for i := h; i < size; i++ {
			for j := i; j >= h && less(items[j].value, items[j-h].value); j -= h {
				listener(j-h, j, items[j-h], items[j])
				exchange(items, j, j-h)
			}
		}
		h = h / 3
	}
}

func mergeSort(items []prefixSum, listener func(i int, j int, iValue prefixSum, jValue prefixSum)) {
	aux := make([]prefixSum, len(items))
	doMergeSort(items, aux, 0, len(items)-1, listener)
}

func doMergeSort(items []prefixSum, aux []prefixSum, low int, high int,
	listener func(i int, j int, iValue prefixSum, jValue prefixSum)) {

	if high <= low {
		return
	}

	mid := low + (high-low)/2

	doMergeSort(items, aux, low, mid, listener)
	doMergeSort(items, aux, mid+1, high, listener)
	merge(items, aux, low, mid, high, listener)
}

func merge(items []prefixSum, aux []prefixSum, low int, mid int, high int,
	listener func(i int, j int, iValue prefixSum, jValue prefixSum)) {

	for k := low; k <= high; k++ {
		aux[k] = items[k]
	}

	i := low
	j := mid + 1

	for k := low; k <= high; k++ {
		if i > mid {
			items[k] = aux[j]
			j++
		} else if j > high {
			items[k] = aux[i]
			i++
		} else if less(aux[j].value, aux[i].value) {
			items[k] = aux[j]
			listener(i, j, aux[i], aux[j])
			j++
		} else {
			items[k] = aux[i]
			i++
		}
	}
}

func less(firstValue int, secondValue int) bool {
	return firstValue < secondValue
}

func exchange(items []prefixSum, i int, j int) {
	temp := items[i]
	items[i] = items[j]
	items[j] = temp
}

// O(n^2) approach:

func findMaximumModuloSumOn2(numberOfItems int, modulo int) (maxSum int) {
	currentPrefixSum := 0

	// Step 1: O(n)
	// Let's create an array of preffix sum items to track their preffix sum and their original indexes.
	preffixSums := make([]int, numberOfItems)

	for i := 0; i < numberOfItems; i++ {
		currentPrefixSum = (readInteger()%modulo + currentPrefixSum) % modulo
		preffixSums[i] = currentPrefixSum
	}

	// Step 2: O(n^2)
	// Let's use brute force to understand the algorithm.
	// To find the preffix sum of a subarray, we use:
	// sumModular[i,j]=(prefixSum[j] - prefixSum[i−1] + M)%M

	maxSum = 0

	for i := 0; i < numberOfItems; i++ {
		for j := i - 1; j >= 0; j-- {
			maxSum = max(maxSum, (preffixSums[i]-preffixSums[j]+modulo)%modulo)
		}
		maxSum = max(maxSum, preffixSums[i]) // Don't forget sum from beginning.
	}

	return
}
