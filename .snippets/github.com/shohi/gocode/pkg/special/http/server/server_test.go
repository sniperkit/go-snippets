package server

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttptestServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleRequest))
	log.Println(server.URL)

	var resp *http.Response
	var err error
	// resp, err := http.Get(server.URL + "/redirect")
	// log.Println(resp, err)

	resp, err = http.Get(server.URL + "/get")

	log.Println("===================response=================")
	log.Println("response code: ", resp.StatusCode)
	log.Println("response error: ", err)
	for k, v := range resp.Header {
		log.Println("response key: ", k, " value: ", v)
	}

}
