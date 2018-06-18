package main

import (
	imagga "github.com/feedlabs/elasticfeed-plugin/indexer/photo-imagga"
	"github.com/feedlabs/elasticfeed/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPipeline(new(imagga.Indexer))
	server.Serve()
}
