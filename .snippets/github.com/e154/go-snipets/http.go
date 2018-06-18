package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"time"
	"net"
)

var (
	domain string = "dev.tire.loc/images/"
	filename string = "logo-w.svg"
)

var timeout = time.Duration(2 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

func saveFile(uri, dir, file string) {
	transport := http.Transport{
		Dial: dialTimeout,
	}

	client := http.Client{
		Transport: &transport,
	}

	resp, err := client.Get(uri)
	resp.Header.Set("User-Agent", "Golang Spider Bot v. 3.0")

	os.MkdirAll(dir, 0777)
	out, err := os.Create(dir+file)
	defer out.Close()
	if(err != nil) {
		fmt.Printf("error 1 %s\n", err)
		os.Exit(1)
	}
	n, err := io.Copy(out, resp.Body)
	defer resp.Body.Close()
	if(err != nil) {
		fmt.Printf("error 2 %s\n", err)
		os.Exit(1)
	}
	fmt.Println(n)

}

func main() {

	uri := fmt.Sprintf("http://%s/%s", domain, filename)
	saveFile(uri, "dir/", filename)

}
