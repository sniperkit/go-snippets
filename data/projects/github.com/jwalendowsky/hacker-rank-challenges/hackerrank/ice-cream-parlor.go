package hackerrank

import (
	"fmt"
	"stdin"
)

/**
 * Sunny and Johnny together have  dollars they want to spend on ice cream. The parlor offers  flavors, and they
 * want to choose two flavors so that they end up spending the whole amount.
 *
 * You are given the cost of these flavors. The cost of the  flavor is denoted by . You have to display
 * the indices of the two flavors whose sum is .
 */

// IceCreamParlor executes the Ice Cream Parlor algorithm.
func IceCreamParlor() {
	testsNumber := stdin.ReadInt() // The number of tests to be read

	for i := 0; i < testsNumber; i++ {
		amount := stdin.ReadInt()        // Amount to be spent by Sunny and Johnny
		flavorsLength := stdin.ReadInt() // Numbers of flavors
		findFlavorsOn(amount, flavorsLength)
	}
}

// This is an On solution to the problem that uses a map to keep track of the amount that is missing for a given value
func findFlavorsOn(amount int, flavorsLength int) {

	// Let's create a map with the missing amount to the final amount and the index of the element
	index := make(map[int]int)
	alreadyFound := false

	// Let's run an On traversal to read the array elements
	for i := 1; i <= flavorsLength; i++ {
		currentValue := stdin.ReadInt()

		// Even when we already found the result, let's read the remain integers of the line just to
		// keep the stream ok.
		if !alreadyFound {
			missingAmount := amount - currentValue
			firstIndex := index[currentValue] // Let's check if there's any item whose missing amount is the given index

			if firstIndex > 0 {
				fmt.Printf("%d %d\n", firstIndex, i)
			} else if missingAmount > 0 {
				index[missingAmount] = i
			}
		}
	}
}

// This is a On2 version of the algorithm
func findFlavorsOn2(amount int, flavorsLength int) {
	flavorsCosts := readIntegerArray(flavorsLength, stdin.ReadInt)

	// Let's run over the flavorsCosts array to find the values whose sum is equal to the given amount.
	for i := 0; i < len(flavorsCosts); i++ {
		for j := i + 1; j < len(flavorsCosts); j++ {
			if flavorsCosts[i]+flavorsCosts[j] == amount {
				fmt.Printf("%d %d\n", i+1, j+1)
				return
			}
		}
	}
}
