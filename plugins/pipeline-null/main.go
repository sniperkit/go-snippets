package main

import (
	"github.com/feedlabs/elasticfeed-plugin/pipeline/null"
	"github.com/feedlabs/elasticfeed/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPipeline(new(null.Pipeline))
	server.Serve()
}
