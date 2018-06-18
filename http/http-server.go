package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Response struct {
	Key   string `json:"key"`
	Value string `json:"value`
}

func handler(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println("handle message from client:", string(reqBody))

	w.Header().Set("Content-type", "application/json; charset=utf-8")
	res := Response{Key: "foo", Value: "world"}
	resJson, err := json.Marshal(res)
	if err != nil {
		log.Println("json ecnoding error ", err)
		return
	}
	w.Write(resJson)
	return
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
