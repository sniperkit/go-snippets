package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	// write file
	writeFile("./test.txt", "hello world")
	// read file
	content := readFile("./test.txt")
	// print file content
	fmt.Printf("file content:\n\n%v", content)
}

func readFile(fileName string) string {
	content, err := ioutil.ReadFile(fileName)
	checkError(err)

	return (string(content))
}

func writeFile(fileName string, content string) {
	file, err := os.Create(fileName)
	defer file.Close()

	len, err := io.WriteString(file, content)
	checkError(err)

	fmt.Printf("written %v chars", len)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
