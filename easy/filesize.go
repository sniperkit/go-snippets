package main

import (
	"fmt"
	"log"
	"os"
)

func filesize(q string) int64 {
	data, err := os.Open(q)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	stat, err := data.Stat()
	if err != nil {
		log.Fatal(err)
	}
	return stat.Size()
}

func main() {
	fmt.Println(filesize(os.Args[1]))
}
