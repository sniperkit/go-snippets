package main

import (
	"flag"
	"sync"

	"log"

	"github.com/gianarb/sherlock/backend/restapi"
)

func main() {
	log.Println("Sherlock version ciao")
	noUi := flag.Bool("no-ui", false, "a bool")
	flag.Parse()
	restConf := &restapi.Config{
		NoUi: noUi,
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func(rconf *restapi.Config) {
		closeConn := restapi.StartRestApi(8080, rconf)
		defer closeConn()
	}(restConf)
	wg.Wait()
}
