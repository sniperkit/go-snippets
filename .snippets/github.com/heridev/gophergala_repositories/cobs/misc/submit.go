package main

import (
	"encoding/json"
	"fmt"

	"github.com/iron-io/iron_go/mq"
)

type Image struct {
	Id     string `json:"id"`
	Tag    string `json:"tag"`
	Source string `json:"source"`
}

func main() {
	machine := "x86_64"
	queueName := "builder-" + machine
	queue := mq.New(queueName)

	image := Image{Id: "moo", Tag: "tag", Source: "source"}
	b, err := json.Marshal(image)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b)
	id, err := queue.PushString(string(b))
	fmt.Println(id)

}
