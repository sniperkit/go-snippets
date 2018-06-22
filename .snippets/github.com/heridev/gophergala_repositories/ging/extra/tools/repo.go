package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

var count = flag.Int("count", 10, "Number of repositories")

type Repository struct {
	Path     string `json:"path"`
	Synopsis string `json:"synopsis,omitempty"`
}

func main() {
	results := struct {
		Repositories []Repository `json:"results"`
	}{[]Repository{}}
	err := json.NewDecoder(os.Stdin).Decode(&results)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
		os.Exit(-1)
	}
	var rs = results.Repositories[0:*count]
	err = json.NewEncoder(os.Stdout).Encode(rs)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
		os.Exit(-1)
	}
}
