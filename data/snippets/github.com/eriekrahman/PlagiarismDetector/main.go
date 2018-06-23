package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/PlagiarismDetector/utils"
)

func main() {	
	if len(os.Args) != 3 {
		fmt.Println("Two arguments should be provided to indicate two texts to compare")
		fmt.Println("Syntax: go run main.go [file_path_1] [file_path_2]")
	} else {
		// Get the arguments
		firstFile := os.Args[1]
		secondFile := os.Args[2]
		
		// Load 2 files
		firstContent := utils.LoadFile(firstFile)
		secondContent := utils.LoadFile(secondFile)
		
		// Calculate similarity
		CalculateSimilarity(firstContent, secondContent)
	}
}

/*****************************************************************************************************
 * Function: CalculateSimilarity
 * Purpose: To calculate similarity and print the result
 * Parameter: firstText as the content of first file, secondText as the content of second file
 * Output: N/A
 *****************************************************************************************************/
func CalculateSimilarity(firstText string, secondText string){
	
	// Initialization
	var weight, averageWeight, totalWeight, maxWeight float32 = 0, 0, 0, 0
	var mostSimilarSentence1, mostSimilarSentence2 string
	var numOfComparison int = 0
	
	// TODO: To improve the result, before calculating the similarity, we can do word stemming
	
	// Remove useless punctuation
	firstText = utils.RemoveUnusedPunctuation(firstText)
	secondText = utils.RemoveUnusedPunctuation(secondText)
	
	// Collect list of sentences by splitting the text by period sign, to indicate 1 sentence
	listOfSentenceFromText1 := strings.Split(strings.TrimSpace(strings.ToLower(firstText)), ".")
	listOfSentenceFromText2 := strings.Split(strings.TrimSpace(strings.ToLower(secondText)), ".")
	
	// For each pair of sentences across the text, do weight calculation
	for _,sentence1 := range listOfSentenceFromText1 {
		for _,sentence2 := range listOfSentenceFromText2 {
			if len(strings.TrimSpace(string(sentence1))) > 0 && len(strings.TrimSpace(string(sentence2))) > 0 {
				
				// Collect list of words by splitting the sentence
				listOfWords1 := strings.Split(strings.TrimSpace(string(sentence1)), " ")
				listOfWords2 := strings.Split(strings.TrimSpace(string(sentence2)), " ")
				
				// Calculate the similarity
				weight = utils.SmithWaterman(listOfWords1, listOfWords2)
				if weight > maxWeight {
					maxWeight = weight
					mostSimilarSentence1 = string(sentence1)
					mostSimilarSentence2 = string(sentence2)
				}
				totalWeight += weight
				numOfComparison++
			}
		}
	}
	
	// Calculate the average weight
	averageWeight = totalWeight / float32(numOfComparison)
	
	// Print the result
	fmt.Println();
	fmt.Println("Text 1: ", firstText)
	fmt.Println();
	fmt.Println("Text 2: ", secondText)
	fmt.Println();
	fmt.Printf("%% Average Similarity = %.2f %%", averageWeight)
	fmt.Println();
	fmt.Printf("%% Max Similarity = %.2f %%", maxWeight)
	fmt.Println();
	fmt.Println("The sentences in File 1 and File 2 that are similar a lot:");
	fmt.Println(" - Text in File 1: ", mostSimilarSentence1);
	fmt.Println(" - Text in File 2: ", mostSimilarSentence2);
	fmt.Println();
}
