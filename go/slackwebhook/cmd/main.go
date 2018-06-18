package main

import (
	"log"
	"net/http"
	"os"

	"github.com/wzulfikar/lab/go/slackwebhook"
)

const listenAddr = "localhost:9090"

func main() {
	handler := slackwebhook.NewHandler(listenAddr, os.Getenv("CONFIG_DIR"))
	http.HandleFunc("/", handler)

	log.Println("Listening at port 9090 ")
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
