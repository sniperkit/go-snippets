package hackerrank

import (
	"fmt"
	"stdin"
)

/**
 * Watson gives Sherlock an array  of length . Then he asks him to determine if there exists an element in the
 * array such that the sum of the elements on its left is equal to the sum of the elements on its right.
 * If there are no elements to the left/right, then the sum is considered to be zero.
 * Formally, find an , such that, 12i-1 i+1i+2N.
 *
 * Input Format
 * The first line contains , the number of test cases. For each test case, the first line contains , the number
 * of elements in the array . The second line for each test case contains  space-separated integers, denoting the array .
 *
 * Output Format
 * For each test case print YES if there exists an element in the array, such that the sum of the elements on its left
 * is equal to the sum of the elements on its right; otherwise print NO.
 */

// SherlockArray executes the Sherlock Array algorithm
func SherlockArray() {
	testsNumber := stdin.ReadInt() // The number of tests to be read

	for i := 0; i < testsNumber; i++ {
		itensNumber := stdin.ReadInt()        // The number of itens in the array
		itensSums := make([]int, itensNumber) // the sums of itens from left to right
		itensSum := 0                         // The current itens sum from left to right
		rightItensSum := 0                    // the itens sum from right to left
		found := "NO"

		// let's read the array and store the itens sum from left to right in each index.
		for j := 0; j < itensNumber; j++ {
			itensSum += stdin.ReadInt()
			itensSums[j] = itensSum
		}

		// let's iterate the array from right to left and check if there is at least one result ok.
		for j := len(itensSums) - 1; j >= 0; j-- {
			sum := 0

			// if we have itens from the left, we'll compare using the sum of itens till j - 1 and compare to the
			// current rightItemsSum
			if j > 0 {
				sum = itensSums[j-1]
			}

			if rightItensSum == sum {
				found = "YES"
				break
			}
			rightItensSum += itensSums[j] - sum // Let's store the sum to verify
		}

		fmt.Println(found)
	}
}
