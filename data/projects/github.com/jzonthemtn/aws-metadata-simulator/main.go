package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
)

type tomlConfig struct {
	Port     int
	Metadata metadata
}

type metadata struct {
	Port           int
	AmiID          string
	LocalHostname  string
	ProductCodes   string
	ReservationID  string
	PublicHostname string
	PublicIPV4     string
}

var m map[string]string

func main() {

	configFile := "metadata.toml"

	if len(os.Args) == 2 {
		configFile = os.Args[1]
	}

	fmt.Printf("Starting EC2 metadata simluator from config %s\n", configFile)

	var config tomlConfig
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		fmt.Println(err)
		return
	}

	m = make(map[string]string)
	m["ami-id"] = config.Metadata.AmiID
	m["local-hostname"] = config.Metadata.LocalHostname
	m["reservation-id"] = config.Metadata.ReservationID
	m["product-codes"] = config.Metadata.ProductCodes
	m["public-hostname"] = config.Metadata.PublicHostname
	m["public-ipv4"] = config.Metadata.PublicIPV4

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/latest/meta-data/{category}", handle).Methods("GET")

	host := fmt.Sprintf(":%d", config.Port)
	fmt.Printf("Listening on: %s\n", host)
	log.Fatal(http.ListenAndServe(host, router))

}

func handle(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	i, ok := m[vars["category"]]

	if ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, i)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}
