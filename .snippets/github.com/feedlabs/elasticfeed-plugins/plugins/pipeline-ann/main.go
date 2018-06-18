package main

import (
	"github.com/feedlabs/elasticfeed-plugin/pipeline/ann"
	"github.com/feedlabs/elasticfeed/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPipeline(new(ann.Pipeline))
	server.Serve()
}
