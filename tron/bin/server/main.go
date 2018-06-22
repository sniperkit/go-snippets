package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/golang/glog"

	_ "github.com/gophergala/tron"
)

var (
	port int
)

func init() {
	flag.IntVar(&port, "port", 8080, "port to bind to")
}

func main() {
	flag.Parse()

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		glog.Fatalf("%v", err)
	}
}
