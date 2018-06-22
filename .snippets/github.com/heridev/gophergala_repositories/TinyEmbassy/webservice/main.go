/*
* @Author: souravray
* @Date:   2015-01-24 11:24:21
* @Last Modified by:   souravray
* @Last Modified time: 2015-01-24 11:33:58
 */

package main

import (
	// "encoding/gob"
	// "github.com/gophergala/tinyembassy/webservice/models"
	"github.com/gophergala/tinyembassy/webservice/router"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	// seting up router
	rtr := mux.NewRouter()
	router.Routes(rtr)
	http.Handle("/", rtr)

	// //registr Audience object to be serialized in seesion object
	// gob.Register(&models.Audience{})

	// start server here
	log.Println("Listening...")

	//http.ListenAndServe(":7979", nil) os.Getenv("PORT")
	http.ListenAndServe(GetPort(), nil)
}

func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "7979"
	}
	return ":" + port
}
