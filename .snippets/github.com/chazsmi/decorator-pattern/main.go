package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chazsmi/decorator-pattern/client"
)

func main() {
	cli := client.Decorate(
		http.DefaultClient,
		client.Auth("verysecuretoken"),
		client.LoadBalancing(
			client.RoundRobin(
				0,
				"http://www.sainsburys.co.uk",
				"http://www.sainsburys.co.uk",
			),
		),
		client.Log(log.New(os.Stdout, "client: ", log.LstdFlags)),
		client.BackOff(5, time.Second),
	)

	req, err := http.NewRequest("GET", "http://www.sainsburys.co.uk", nil)
	if err != nil {
		fmt.Printf("%+v", err)
	}
	resp, err := cli.Do(req)
	fmt.Printf("%+v", resp)
}
