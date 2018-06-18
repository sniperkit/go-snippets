package main

import "fmt"

func main() {
	openFile()
	defer closeFile() // defer executes before the func main exit
	writeFile()
}

func openFile() {
	fmt.Println("Open File")
}

func writeFile() {
	fmt.Println("Write File ...")
}

func closeFile() {
	fmt.Println("Close File")
}
