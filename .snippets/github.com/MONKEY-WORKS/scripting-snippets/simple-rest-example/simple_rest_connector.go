/*
Copyright 2017 MONKEY WORKS GmbH

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

1. 	Redistributions of source code must retain the above copyright notice,
	this list of conditions and the following disclaimer.

2. 	Redistributions in binary form must reproduce the above copyright notice,
	this list of conditions and the following disclaimer in the documentation
	and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY
EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT,
INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

/*
 *
 * Simple REST connector to consume data provided by a web service
 *  Author: Jan Hoyer
 *  Date: 18.07.2017
 *
 *
 */
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//The endpoint of our web service we want to connect to
//Change this to your endpoint
const URL string = "http://localhost:8080/SPSSimulator/rest/spssimulator"

//The interval at which the connector is updated
//Time is defined in milliseconds
const UPDATE_INTERVAL time.Duration = 1000

// Global channel to consume received data
var Data = make(chan *Response)

// Global channel to control periodic task
var Stop = make(chan bool)

// Data structure of simple data item provided by REST interface
// Adapt this to your data structure
type DataItem struct {
	Name      string  `json:name`
	Value     float32 `json:value`
	Timestamp string  `json:timestamp`
}

// The repsonse of REST interface
// Adapt this to your data structure
type Response struct {
	DataItems []DataItem `json:"dataItems`
}

// starts a thread to receive data periodically and consumes messages from the global channels
func main() {
	go getData()
	fmt.Println("Subscription started")
	for {
		select {
		case val := <-Data:
			processData(val)
		case <-Stop:
			fmt.Println("Stop")
			return
		}
	}
}

// processing of received data
func processData(data *Response) {
	if data == nil {
		fmt.Println("Error")
		return
	}
	for _, value := range data.DataItems {
		fmt.Println("DataItem ", value.Name, " Value ", value.Value, " Timestamp ", value.Timestamp)
		//
		// Add your code here ...
		//
	}

}

// starts every UPDATE_INTERVAL milliseconds a HTTP Request to get data from web service
// the response is forwarded to the global channel if available
func getData() {

	for range time.Tick(time.Millisecond * UPDATE_INTERVAL) {
		// Build the request
		req, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			Data <- nil
			continue
		}

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			Data <- nil
			continue
		}

		if resp.StatusCode != 200 { // OK
			Data <- nil
			continue
		}

		var response Response
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&response)
		Data <- &response
		resp.Body.Close()
	}
	Stop <- true
}
