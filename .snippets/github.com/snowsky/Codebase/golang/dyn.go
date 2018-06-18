package main

import (
	"log"

	"github.com/nesv/go-dynect/dynect"
)

func main() {
	client := dynect.NewClient("thinkingphones")
	client.Verbose(true)
	err := client.Login("fuze-api-ro", "7#zD$7ovB$jU")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := client.Logout()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Make a request to the API, to get a list of all, managed DNS zones
	var response dynect.ZonesResponse
	if err := client.Do("GET", "Zone", nil, &response); err != nil {
		log.Println(err)
	}
	log.Println(response)
	// for _, zone := range response.Data {
	// 	log.Println("Zone", zone)
	// }
}
