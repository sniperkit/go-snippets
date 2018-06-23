package main

// Panopticon, Copyright © 2015, Huck Ridge Software LLC
// All rights reserved.

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// vim:sw=4:ts=4

var user, server string
var client = http.Client{}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: sensor your_email_address [server_url]")
	}
	user = os.Args[1]

	server = "plexiform-leaf-835.appspot.com"
	if len(os.Args) >= 3 {
		server = os.Args[2]
	}

	sendEntry()
	for _ = range time.Tick(15 * time.Second) {
		sendEntry()
	}
}

func sendEntry() {
	e, err := MakeEntry()
	if err != nil {
		log.Fatalf("MakeEntry: %v", err)
	}
	if e.WasIdle && e.Idle > CONSIDERED_IDLE {
		log.Printf("WasIdle & idletime = %v; not sending an update", e.Idle)
		return
	}

	jsonSampleEntry, err := json.Marshal(&e)
	if err != nil {
		log.Fatalf("Couldn't marshal sampleEntry: %v", err)
	}
	req, err := http.NewRequest("PUT",
		"http://"+server+"/api/v1/add_entry",
		bytes.NewBuffer(jsonSampleEntry))
	if err != nil {
		log.Fatalf("Failed to create req: %v", err)
	}
	req.Header.Set("X-Panopticon-Token", user)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error Do'ing the request: %v", err)
	}

	if resp.StatusCode != 200 {
		log.Printf("resp.StatusCode expected to be 200, not %v", resp.StatusCode)
	} else {
		if _, ok := resp.Header["Location"]; !ok {
			log.Fatalf("No Location header in the response")
		}
		log.Printf("Sent %v", e)
		// log.Printf("Response: %v", resp)
	}
}
