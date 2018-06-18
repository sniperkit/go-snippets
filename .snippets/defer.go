package main

import "fmt"
import "os"

func createFile(p string) *os.File {
	fmt.Println("Creating file")
	f, err := os.Create(p)
	if err != nil {
		panic(err)
	}
	return f
}

func writeToFile(f *os.File) {
	fmt.Println("Writing to file")
	fmt.Fprintln(f, "data")
}

func closeFile(f *os.File) {
	fmt.Println("Closing file")
	f.Close()
}

func main() {
	f := createFile("/tmp/defer.txt")
	defer closeFile(f)
	writeToFile(f)
}
