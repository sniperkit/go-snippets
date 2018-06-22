package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	fmt.Println("Hello, Linen Project")
	fin, err := os.Open("res/dict.txt")
	if err == nil {
		fmt.Println("File exists!")
		r := bufio.NewReader(fin)
		for line, _, err := r.ReadLine(); err != io.EOF; line, _, err = r.ReadLine() {
			fmt.Printf("Lines: %v (error %v)\n", string(line), err)
		}
	} else {
		fmt.Println("Error: the file to be read is not found")
	}
}

