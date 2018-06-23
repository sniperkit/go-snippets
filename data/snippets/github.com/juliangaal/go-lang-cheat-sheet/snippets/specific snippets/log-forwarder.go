package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Struct handles user logins
type Login struct {
	user string
	time time.Time
}

// struct handles pricing data requests by user
type PriceDataRequest struct {
	user           string
	time           time.Time
	priceDataRange string
}

// struct handles
type FileImport struct {
	status string
	err    error
	time   time.Time
}

// struct handles go server messages
type Server struct {
	err              error
	dbStatus         string
	connectionStatus string
}

func main() {

	logs := []interface{}{
		Login{"Julian", time.Now()},
		Server{nil, "Online", "Available"},
		FileImport{"Success", nil, time.Now()},
		PriceDataRequest{"Julian", time.Now(), "03.17.-06.17."},
	}

	// logs = append(logs, Login{"Julian", time.Now()})
	//logs
	url := "127.0.0.1:5000"

	for _, l := range logs {
		b, err := json.Marshal(l)
		if err != nil {
			fmt.Println("Error parsing json")
		}

		if _, err := http.NewRequest("POST", url, bytes.NewBuffer(b)); err != nil {
			log.Fatal(err)
		}
	}
}
