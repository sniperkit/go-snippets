package main

import (
	_ "core"
	"log"
	"net/http"
)

func main() {
	log.Println("go-testy Hello World!")

	// TODO: implement the HTTP listener in the way that Docker prefers.
	/*
		This implementation is an adaptation from the way the Google App Engine does it,
		and I've continued to use it for convience' sake, which is the following:
		1. Place application code into a 'core' directory.
		2. Load the 'core' directory as a by-product package.
		3. Load the HTTP listener knowing that 'core' provides all the routes, etc.
	*/

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Failure starting web process: " + err.Error())
	}

}
