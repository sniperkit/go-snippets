package main

import (
	"encoding/json"
	"fmt"
	"log"
)

var my_json string = `{"name": "Chiheb","surname": {"first": "NeXus","second": {"foo": "nexus", "bar": "neXus"}}}`

type Response struct {
	Name    string          `json:"name"`
	Surname ResponseSurname `json:"surname"`
}

type ResponseSurname struct {
	First  string         `json:"first"`
	Second ResponseSecond `json:"second"`
}

type ResponseSecond struct {
	Foo  string `json:"foo"`
	Baar string `json:"bar"`
}

func main() {
	res := Response{}
	err := json.Unmarshal([]byte(my_json), &res)
	if err != nil {
		log.Fatal("Unmarshall failed !\n", err)
	}
	fmt.Println(res)
}
