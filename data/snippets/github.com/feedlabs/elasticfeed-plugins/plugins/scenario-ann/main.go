package main

import (
	"github.com/feedlabs/elasticfeed-plugin/scenario/ann"
	"github.com/feedlabs/elasticfeed/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPipeline(new(ann.Scenario))
	server.Serve()
}
