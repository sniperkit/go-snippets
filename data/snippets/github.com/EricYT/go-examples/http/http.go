package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type Test struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	fmt.Println("test init ")

	test := Test{Key: "foo", Value: "bar"}
	testJson, err := json.Marshal(test)
	if err != nil {
		fmt.Println("init message error ", err)
		return
	}
	fmt.Println("init encode message ", string(testJson))

	/* send http request */
	t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
	}

	client := http.Client{}
	client.Transport = t
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/", bytes.NewBuffer(testJson))

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("client do request error ", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("Receive response from server:", string(body))

	resTest := &Test{}
	err = json.Unmarshal(body, resTest)
	if err != nil {
		fmt.Println("Unmarshal json error ", err)
		return
	}

	fmt.Printf("Receive response struct:%+v\n", resTest)
}
