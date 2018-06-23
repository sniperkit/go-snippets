package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", test)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed to start server!")
	}
}

func test(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("test")
	io.WriteString(rw, "test")
}
