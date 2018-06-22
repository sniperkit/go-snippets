package main

import (
	"flag"

	"github.com/dynamicgo/aliyunlog"
	"github.com/dynamicgo/config"
	"github.com/dynamicgo/slf4go"
	"github.com/inwecrypto/neo-insight/insight"
	_ "github.com/lib/pq"
)

var logger = slf4go.Get("neo-indexer")
var configpath = flag.String("conf", "./insight.json", "neo indexer config file path")

func main() {

	flag.Parse()

	neocnf, err := config.NewFromFile(*configpath)

	if err != nil {
		logger.ErrorF("load neo config err , %s", err)
		return
	}

	factory, err := aliyunlog.NewAliyunBackend(neocnf)

	if err != nil {
		logger.ErrorF("create aliyun log backend err , %s", err)
		return
	}

	slf4go.Backend(factory)

	insight.OpenLogger()

	server, err := insight.NewServer(neocnf)

	if err != nil {
		logger.ErrorF("load neo config err , %s", err)
		return
	}

	server.Run()
}
