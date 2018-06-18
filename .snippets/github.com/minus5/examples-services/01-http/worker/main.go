package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	port = ":9002"
)

func add(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	var numbers []int
	decoder.Decode(&numbers)
	fmt.Println("Got request", numbers)
	sum := 0
	for _, value := range numbers {
		sum += value
	}
	fmt.Println("Sending response", sum)
	fmt.Fprintf(rw, "%d", sum)
}

func main() {
	fmt.Printf("Started worker at http://localhost%v.\n", port)
	http.HandleFunc("/", add)
	http.ListenAndServe(port, nil)
}
