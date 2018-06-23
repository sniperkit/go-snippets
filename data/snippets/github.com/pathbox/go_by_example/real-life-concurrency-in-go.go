package main

import (
	"fmt"
	"net/http"
	"time"
)

var urls = []string{
	"http://www.hao123.com",
	"http://www.baidu.com",
	"https://ruby-china.org/",
}

type HttpResponse struct {
	url      string
	response *http.Response
	err      error
}

func asyncHttpGets(urls []string) []*HttpResponse {
	ch := make(chan *HttpResponse)
	responses := []*HttpResponse{}
	client := http.Client{}
	for _, url := range urls {
		go func(url string) {
			fmt.Printf("Fetching %s", url)
			resp, err := client.Get(url)
			ch <- &HttpResponse{url, resp, err}
			if err != nil && resp != nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
			}
		}(url)
	}

	for {
		select {
		case r := <-ch:
			fmt.Printf("%s was fetched\n", r.url)
			if r.err != nil {
				fmt.Println("with an error", r.err)
			}
			responses = append(responses, r)
			if len(responses) == len(urls) {
				return responses
			}
		default:
			fmt.Printf(".")
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func main() {
	results := asyncHttpGets(urls)
	for _, result := range results {
		if result != nil && result.response != nil {
			fmt.Printf("%s status: %s\n", result.url,
				result.response.Status)
		}
	}
}
