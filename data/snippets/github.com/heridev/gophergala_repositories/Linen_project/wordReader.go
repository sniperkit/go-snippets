// Simple go package to read a word at a time from a file
// Used to following example to learn from: http://golang.org/pkg/bufio/#example_Scanner_words
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	// A file input source.
	fin, err := os.Open("res/dict.txt")
	bagOfWords := make(map[string]int)

	if err == nil {
		r := bufio.NewReader(fin)
		for strFromReader, err := r.ReadString('\n'); err != io.EOF; strFromReader, err = r.ReadString('\n') {
			fmt.Printf("Line: %v (error %v)\n", strFromReader, err)
			// create another scanner for that string
			strScanner := bufio.NewScanner(strings.NewReader(strFromReader))

			// Set the split function for the scanning operation.
			strScanner.Split(bufio.ScanWords)
			// Count the words.
			count := 0
			for strScanner.Scan() {
				word := strScanner.Text()
				if count, ok := bagOfWords[word]; ok{
					bagOfWords[word] = count + 1
				} else {
					bagOfWords[word] = 1
				}
				fmt.Printf("%s, ", word)
				count++
			}
			if err := strScanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading input:", err)
			}
			fmt.Printf("%d\n", count)

		}
		fmt.Printf("\nHere is the result of the word count...\n")
		for key, value := range bagOfWords {
			fmt.Println("Key:", key, "Value:", value)
		}
	} else {
		fmt.Println("Error: the file to be read is not found")
	}
}

