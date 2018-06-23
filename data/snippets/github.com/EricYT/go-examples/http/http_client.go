package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

type RoundTripper interface {
	RoundTrip(req, res interface{}) error
}

type Client struct {
	http.Client

	t *http.Transport
	u *url.URL
}

func ParseURL(s string) (*url.URL, error) {
	var err error
	var u *url.URL

	if s != "" {
		u, err = url.Parse(s)
		if err != nil {
			return nil, err
		}
	}

	return u, nil
}

func NewClient(s string) (*Client, error) {
	u, err := ParseURL(s)
	if err != nil {
		return nil, err
	}

	c := Client{
		u: u,
	}
	c.t = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			// Enough ?
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
	}

	c.Client.Transport = c.t
	c.u = c.URL()

	return &c, nil
}

func (c *Client) URL() *url.URL {
	urlCopy := *c.u
	return &urlCopy
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	var resc = make(chan *http.Response, 1)
	var errc = make(chan error, 1)

	// Perform request from separate routine
	go func() {
		res, err := c.Client.Do(req)
		if err != nil {
			errc <- err
		} else {
			resc <- res
		}
	}()

	// Wait for request completion
	select {
	// Request timeout
	case <-time.After(time.Second * 30):
		fmt.Printf("http client do timeout.Request is %+v\n", req)
		return nil, errors.New("timeout")
	case err := <-errc:
		return nil, err
	case res := <-resc:
		return res, nil
	}
}

func (c *Client) RoundTripRAWData(reqBody []byte) error {
	var err error
	log.Println("-------- RoundTrip reqBodyjson:", string(reqBody))

	reqContentReader := bytes.NewBuffer(reqBody)
	req, err := http.NewRequest("POST", c.u.String(), reqContentReader)
	if err != nil {
		panic(err)
	}
	//	req.Close = true
	//req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("Transfer-Encoding", "chunked")

	tstart := time.Now()
	res, err := c.do(req)
	tstop := time.Now()
	fmt.Printf("Http request time:%6dms \n", tstop.Sub(tstart)/time.Millisecond)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	for key, value := range res.Header {
		fmt.Println("Header:Key:", key, " Value:", value)
	}

	fmt.Printf("http request code:%d", res.StatusCode)
	switch res.StatusCode {
	case http.StatusOK:
		// OK
	case http.StatusInternalServerError:
		fmt.Printf("http client 500")
		panic("xxxxxxxxxxxxxxxxxxxxxxxxxx")
		// Error, 500
	default:
		return errors.New(res.Status)
	}
	fmt.Printf("http request code1111:%d\n", res.StatusCode)

	// There is no body,pass it.
	if res.ContentLength > 0 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println("IO read body error:", err)
			panic("read body error")
		}
		fmt.Println("response body:", body)
	} else {
		fmt.Printf("Body is empty")
	}

	return nil
}

func main() {
	fmt.Println("request client")

	client, err := NewClient("http://localhost:9090")
	if err != nil {
		panic("create http client error")
	}

	client.RoundTripRAWData([]byte("hello,world"))

}
