package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func doCycle() {
	fmt.Println("---===---\nCycle start")
	r1, _ := http.Get("http://localhost:9001")
	sensorData, _ := ioutil.ReadAll(r1.Body)
	fmt.Println("Sensor values", string(sensorData))
	r2, _ := http.Post("http://localhost:9002", "application/json", bytes.NewBuffer(sensorData))
	workerData, _ := ioutil.ReadAll(r2.Body)
	fmt.Println("Cycle result:", string(workerData))
}

func main() {
	for range time.Tick(time.Second * 1) {
		doCycle()
	}
}
