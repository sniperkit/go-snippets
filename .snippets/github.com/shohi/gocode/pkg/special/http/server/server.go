package server

import (
	"log"
	"net/http"
	"strings"
)

func init() {
	log.Println("hello")
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		log.Println("key ===> ", k, " value ===> ", v)
	}
	w.WriteHeader(http.StatusOK)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("aws", "aws-authenticatioin")
	w.Header().Set("Location", "https://www.baidu.com")
	w.WriteHeader(http.StatusFound)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	log.Println("URL ==> ", r.URL, " host ==> ", r.Host)

	if strings.HasPrefix(uri, "/redirect") {
		handleRedirect(w, r)
	} else {
		handleGet(w, r)
	}
}
