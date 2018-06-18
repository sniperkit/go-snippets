package golangSnippets

import (
	"log"
	"net/http"

	"golang.org/x/net/proxy"
)

//SocksClient implements socks5 proxy it returns an *http.Client which can be used
//to access hidden web, inputs are host of socks5 proxy server & its port
func SocksClient(host string, port string) *http.Client {
	// create a socks5 dialer
	dialer, err := proxy.SOCKS5("tcp", host+":"+port, nil, proxy.Direct)
	if err != nil {
		log.Fatal("Can't connect to proxy: ", err)
	}

	httpTransport := &http.Transport{}
	// set our socks5 as the dialer
	httpTransport.Dial = dialer.Dial
	return &http.Client{Transport: httpTransport}
}
