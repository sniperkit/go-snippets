package http

import (
	"log"
	"net/url"
	"testing"
)

func TestUrlPathEscape(t *testing.T) {
	log.Println(url.PathEscape("lang:>50"))
	log.Println(url.PathEscape("https://www.amazon.com"))
}

func TestUrlParse(t *testing.T) {
	urlStr := "http://localhost:9090/hello/newyorker?season=summer"
	log.Println(url.Parse(urlStr))

	urlStr = "http://ip/?action=save"

	myURL, _ := url.Parse(urlStr)
	log.Println(url.ParseQuery(myURL.RawQuery))
}

func TestUrlResolveReference(t *testing.T) {
	urlStr := "http://localhost:9090/hello/newyorker?season=summer"
	base, _ := url.Parse(urlStr)

	log.Println(base)

	urlStr1 := "a/b/c/d.ts"
	url1, err := url.Parse(urlStr1)
	log.Println(err)

	log.Println(base.ResolveReference(url1))
}

func TestURLQuery(t *testing.T) {
	urlStr := "http://localhost:9090/hello/newyorker?season=summer&season=spring&show=tony&nokey"
	base, _ := url.Parse(urlStr)

	for key, value := range base.Query() {
		log.Printf("key ==> %s, value ==> %v", key, value)
	}

	log.Printf("raw query: %v, parsed query: %v", base.RawQuery, base.Query())
	log.Printf("Scheme: %v, Host: %v, Port: %v, Path: %v", base.Scheme, base.Host, base.Port(), base.Path)
}

func TestURLPathPrefix(t *testing.T) {
	urlStr := "http://ip/stream1/segment1"
	base, _ := url.Parse(urlStr)

	log.Println(base.Path)
}
