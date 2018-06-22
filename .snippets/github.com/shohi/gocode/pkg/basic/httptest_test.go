package basic

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPTestServer(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	log.Printf("server addr: %v, server: %v", server.URL, server)
}

func TestHTTPRecorderReuse(t *testing.T) {
	w := httptest.NewRecorder()

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, r.URL.String())
	}

	r, _ := http.NewRequest("GET", "http://localhost/123", nil)
	handler(w, r)

	log.Println(w.Code)
	data, _ := ioutil.ReadAll(w.Body)

	log.Println(string(data))

	r, _ = http.NewRequest("GET", "http://localhost/321", nil)
	handler(w, r)

	log.Println(w.Code)
	data, _ = ioutil.ReadAll(w.Body)

	log.Println(string(data))

}
