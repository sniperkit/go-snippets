package golangSnippets

import (
	"bufio"
	"os"
	"strings"
)

//ReadLinesFromFilePtr takes pointer to a readable file as argument & return linesay of strings line by line (without '\n')
func ReadLinesFromFilePtr(fi *os.File) ([]string, error) {
	//Not an empty file
	var lines []string
	_, err := fi.Stat()
	if err != nil {
		return lines, err
	} else {
		scanner := bufio.NewScanner(fi)
		for scanner.Scan() {
			lines = append(lines, strings.TrimSpace(scanner.Text()))
		}
		if scanner.Err() != nil {
			return lines, scanner.Err()
		}
		return lines, nil
	}
}
