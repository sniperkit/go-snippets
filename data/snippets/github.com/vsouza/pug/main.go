package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/vsouza/pug/api"
	"github.com/vsouza/pug/config"
)

func main() {
	mode := flag.String("env", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: mode [-env development] ")
		os.Exit(1)
	}
	flag.Parse()

	var err error
	config, err := config.NewConfig(*mode)
	if err != nil {
		log.Fatal(err)
	}

	api.StartService(config)
	log.Println("service started")
}
