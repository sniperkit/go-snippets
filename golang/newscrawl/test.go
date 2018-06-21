package main

import (
	"github.com/guotie/gogb2312"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("http://www.toutiao.com/api/pc/focus/")
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	body, err, _, _ = gogb2312.ConvertGB2312(body)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer resp.Body.Close()

	log.Println(string(body))
}
