package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/gophergala/gomp/preproc"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var buffer bytes.Buffer
	for scanner.Scan() {
		buffer.WriteString(scanner.Text() + "\n")
	}

	result, err := preproc.PreprocFile(buffer.String(), "stdin")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Gompp error while using preproc.PreprocFile:\n", err.Error())
		os.Exit(-1)
	}
	fmt.Print(result)
}
