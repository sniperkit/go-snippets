package main

import (
	"github.com/feedlabs/elasticfeed-plugin/sensor/weather"
	"github.com/feedlabs/elasticfeed/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPipeline(new(weather.Sensor))
	server.Serve()
}
