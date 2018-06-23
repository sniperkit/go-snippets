package main

import (
	"bitodd/config"
	"bitodd/twitter"
	"bitodd/model"
	"bitodd/pages"
	"flag"
	"log"
	"net/http"
)

var (
	configFileLocation string
)

func init() {
	// Config file location
	flag.StringVar(&configFileLocation, "conf", "./bitodd.conf", "Config file location")

	flag.Parse()
}

func main() {

	// Load config file
	log.Println("Loading config from", configFileLocation)
	err := config.Load(configFileLocation)
	if err != nil {
		log.Fatalln("Configuration error:", err)
	}

	// Start websocket hub
	go model.Hub.Run()

    // Start twitter streaming
	go twitter.OpenStream()

	// Static resources
	router := pages.GetRouter()
	router.Handle("/{id:.+}", http.FileServer(http.Dir("resources")))

	// Routes
	http.Handle("/", router)

	log.Println("Starting web server on port", config.GetConfig().Port)
	if err := http.ListenAndServe(":"+config.GetConfig().Port, nil); err != nil {
		log.Fatalln("Error starting web server:", err)
	}
}
