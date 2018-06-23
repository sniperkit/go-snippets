package http

import (
	"bytes"
	"log"
	"net/http"
	"testing"
)

func TestParseRequest(t *testing.T) {
	url := "http://localhost?season=summer&season=spring&show=tony&nokey&srcurl=http://localhost:8082"

	req, _ := http.NewRequest("GET", url, nil)
	values := req.URL.Query()

	key := "season"
	log.Printf("key: %v, value: %v", key, values.Get(key))

	log.Printf("key: %v, value: %v", "srcurl", values.Get("srcurl"))

	// case sensitive!!!
	key = "SHOW"
	log.Printf("key: %v, value: %v", key, values.Get(key))
}

func TestGetRequest(t *testing.T) {
	resp, err := http.NewRequest("GET", "https://www.douban.com", nil)
	log.Println(resp, err)

	var buf *bytes.Buffer
	if buf == nil {
		http.NewRequest("GET", "https://www.douban.com", nil)
		// Wrong, interface cast will create some wired result
		// http.NewRequest("GET", "https://www.douban.com", buf)
	} else {
		resp, err = http.NewRequest("GET", "https://www.douban.com", buf)
	}
	log.Println(resp, err)
}
