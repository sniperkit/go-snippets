package main // import "github.com/CenturyLinkLabs/sample-go-adapter"

import (
	"log"
	"github.com/CenturyLinkLabs/pmxadapter"
	"github.com/CenturyLinkLabs/sample-go-adapter/sample"
)

func main() {
	adapter := new(sample.SampleAdapter)
	server := pmxadapter.NewServer(adapter)

	log.Printf("Starting Sample Adapter")
	server.Start()
}
