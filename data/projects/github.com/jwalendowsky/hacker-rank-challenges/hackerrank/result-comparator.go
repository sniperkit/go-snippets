package hackerrank

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"stdout"
	"strings"
)

const fileNamePosition = 2 // File name position in the runtime args array

// ResultCompare compares the result from stdout and a file whose name was given
func ResultCompare() {
	if len(os.Args) < 3 {
		fmt.Println("You must provide a file name to compare the result with.")
		os.Exit(1)
	}

	fileName := os.Args[fileNamePosition]

	file, error := os.Open(fileName)

	if error != nil {
		fmt.Printf("Error opening file %s.\n", fileName)
		os.Exit(1)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	line, readError := reader.ReadString('\n')

	fmt.Printf("Actual - Expected - Ok?\n")

	for readError != io.EOF {
		lineFromStdout := stdout.ReadLine()
		line = strings.TrimSuffix(line, "\n")

		fmt.Printf("%s - %s - %t\n", lineFromStdout, line, lineFromStdout == line)

		if lineFromStdout != line {
			fmt.Printf("ERROR: Expected %s, but received %s.", line, lineFromStdout)
			break
		}

		line, readError = reader.ReadString('\n')
	}
}
