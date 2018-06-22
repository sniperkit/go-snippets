package basic

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestIOUtilRead(t *testing.T) {
	data := []byte{}

	_, err := ioutil.ReadAll(bytes.NewReader(data))
	log.Println(err)
}

func TestResponseRecorder(t *testing.T) {
	// data := []byte{}
	var data []byte
	log.Println(data)

	data = []byte{}
	log.Println(data)

	handler := func(w http.ResponseWriter, r *http.Request) {
		n, err := w.Write(data)
		log.Println(n, err)
		// io.WriteString(w, "<html><body>Hello World!</body></html>")
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	log.Println("response body", resp.Body)
	resp.Body = nil
	body, _ := ioutil.ReadAll(resp.Body)

	log.Println(resp.StatusCode)
	log.Println(resp.Header.Get("Content-Type"))
	log.Println(string(body))

	// Output:
	// 200
	// text/html; charset=utf-8
	// <html><body>Hello World!</body></html>
}

func TestReadDir(t *testing.T) {
	dirpath := "../.."
	fileList, _ := ioutil.ReadDir(dirpath)
	abspath, _ := filepath.Abs(dirpath)
	log.Println(abspath)
	for _, f := range fileList {
		log.Println(filepath.Join(dirpath, f.Name()))
	}
}

func TestReadFileInDir(t *testing.T) {
	dirpath := "."
	fileList, _ := ioutil.ReadDir(dirpath)
	for _, f := range fileList {
		log.Println(filepath.Join(dirpath, f.Name()))
	}
}

func TestRemoveDir(t *testing.T) {
	dirpath := "non-exist/non-exist"
	err := os.RemoveAll(dirpath)
	log.Println(err)
}
