/*
* @Author: souravray
* @Date:   2015-01-24 10:11:13
* @Last Modified by:   souravray
* @Last Modified time: 2015-01-24 10:57:44
 */

package main

import (
	"encoding/gob"
	"github.com/gophergala/tinyembassy/site/models"
	"github.com/gophergala/tinyembassy/site/router"
	"github.com/gorilla/mux"
	// "labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"os"
)

func main() {
	// seting up router
	rtr := mux.NewRouter()
	router.Routes(rtr)
	http.Handle("/", rtr)

	// //registr Advertiser object to be serialized in seesion object
	gob.Register(&models.Advertiser{})

	// start server here
	log.Println("Listening...")

	//http.ListenAndServe(":8080", nil) os.Getenv("PORT")
	http.ListenAndServe(GetPort(), nil)
}

func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "8080"
	}
	return ":" + port
}
