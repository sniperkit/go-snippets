package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

//Upload file to remote
func main() {

}

func PostFile(filename, url string) error {
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)

	fwriter, err := writer.CreateFormFile("Hello world", filename)
	if err != nil {
		log.Println(err)
		return err
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = io.Copy(fwriter, file)
	if err != nil {
		log.Println(err)
		return err
	}

	ctype := writer.FormDataContentType()
	writer.Close()

	resp, err := http.Post(url, ctype, buffer)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(string(body))
}
