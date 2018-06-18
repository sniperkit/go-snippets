package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {

	file, err := ioutil.TempFile("", "example_tmp_file")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	fmt.Println("create temp file:", file.Name())

	body := []byte("temp file body")
	if _, err = file.Write(body); err != nil {
		log.Fatal(err)
	}
}
