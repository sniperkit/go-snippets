package main

import (
	"github.com/rsmohamad/comp4321/controllers"
	"log"
	"net/http"
)

func main() {
	controllers.LoadHome()
	controllers.LoadSearch()
	controllers.LoadHistory()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
