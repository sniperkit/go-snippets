package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

const (
	port = ":9001"
)

func read(rw http.ResponseWriter, req *http.Request) {
	fmt.Println("Got request")
	var values [3]uint32
	for i := range values {
		values[i] = rand.Uint32()
	}
	jsonValue, _ := json.Marshal(values)
	fmt.Println("Sending values", string(jsonValue))
	fmt.Fprintf(rw, string(jsonValue))
}

func main() {
	fmt.Printf("Started sensor at http://localhost%v.\n", port)
	http.HandleFunc("/", read)
	http.ListenAndServe(port, nil)
}
