package main

import (
	"goworkshop/web"
	"goworkshop/persistence"
)

func main() {
	dataStore, err := persistence.InitDB()
	if err != nil {
		panic(err)
	}

	server := web.RestServer{
		Port:  8080,
		Store: &dataStore,
	}

	server.StartServer()
}
