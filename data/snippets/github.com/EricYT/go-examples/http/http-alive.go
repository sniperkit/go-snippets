package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func PrintLocalDial(network, addr string) (net.Conn, error) {
	dial := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	conn, err := dial.Dial(network, addr)
	if err != nil {
		return conn, err
	}

	fmt.Println("connect done, use", conn.LocalAddr().String())

	return conn, err
}

func doGet(client *http.Client, url string, id int) {
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	buf, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%d: %s -- %v\n", id, string(buf), err)
	if err := resp.Body.Close(); err != nil {
		fmt.Println(err)
	}
}

const URL = "http://localhost:8080/"

func main() {
	client := &http.Client{
		Transport: &http.Transport{
			Dial:                PrintLocalDial,
			MaxIdleConnsPerHost: 4,
		},
	}
	for {
		go doGet(client, URL, 1)
		go doGet(client, URL, 2)
		go doGet(client, URL, 3)
		go doGet(client, URL, 4)
		time.Sleep(5 * time.Millisecond)
	}
}
