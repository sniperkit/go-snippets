package main

import (
	"fmt"
	"hackerrank"
	"os"
)

const challengeNameIndex = 0

/**
 * Main function of the hacker rank challenge executor.
 */
func main() {

	//  let's use an slice of the args array starting from the second index.
	argsWithoutProg := os.Args[1:]
	var challengeName string

	if len(argsWithoutProg) == 0 {
		fmt.Println("You must provide a challenge name as the argument.")
		os.Exit(1)
	}

	// The first argument of the array contains the challenge name
	challengeName = argsWithoutProg[challengeNameIndex]

	// go does not provide a good reflection mechanis to gets a package's function by its name :(
	switch challengeName {
	case "IceCreamParlor":
		hackerrank.IceCreamParlor()
	case "SherlockArray":
		hackerrank.SherlockArray()
	case "MaximizeSum":
		hackerrank.MaximizeSum()
	case "InsertionSortPart2":
		hackerrank.InsertionSortPart2()
	case "ResultCompare":
		hackerrank.ResultCompare()
	case "QuickSortPart1":
		hackerrank.QuickSortPart1()
	case "QuickSortPart2":
		hackerrank.QuickSortPart2()
	case "InsertionSortAdvancedAnalysis":
		hackerrank.InsertionSortAdvancedAnalysis()

	default:
		fmt.Println("You must provide an existing challenge name.")
	}
}
