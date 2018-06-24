package webcrawler

import (
	"fmt"
	"net/http"
	"time"
)

var tr = &http.Transport{
	MaxIdleConnsPerHost: 1024,
	TLSHandshakeTimeout: 0 * time.Second,
}
var fetchClient = http.Client{
	Timeout:   time.Second * 30,
	Transport: tr,
}

var robotClient = http.Client{Timeout: time.Second * 5}

var throttle = time.Tick(time.Millisecond * 100)

func initClients(aggressive bool) {
	if !aggressive {
		return
	}

	fmt.Println("Running crawler aggressively")

	tr = &http.Transport{
		MaxIdleConnsPerHost: 32767,
		TLSHandshakeTimeout: 0 * time.Second,
	}

	fetchClient = http.Client{
		Timeout:   time.Second * 5,
		Transport: tr,
	}

	throttle = time.Tick(time.Millisecond * 1)
}
