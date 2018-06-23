package basic

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestHTTPGet(t *testing.T) {

	//
	c := &http.Client{}

	// success one
	resp, err := c.Get("https://www.douban.com")
	if err == nil {
		defer resp.Body.Close()
		io.Copy(ioutil.Discard, resp.Body)
	}
	// log.Println(resp, err)

	// failure one
	resp, err = c.Get("https://localhost:12345")
	log.Println(resp, err)

	log.Println(resp == nil)
}

func TestCreateHTTPRequest(t *testing.T) {
	req, err := http.NewRequest("PUT", "localhost:8080", nil)
	if err != nil {
		log.Println(err)
	}
	log.Println(req.Body == nil)

	// ErrorCase
	// leading space error, ref: https://github.com/golang/go/issues/24246
	req, err = http.NewRequest("GEt", " http:/localhost:8080", nil)
	log.Printf("err: %v", err)
}
