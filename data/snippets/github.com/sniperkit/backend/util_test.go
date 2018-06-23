package backend

import (
	"testing"
)

type urlTestCase struct {
	url            string
	wantEncodedUrl string
}

func TestCleanURL(t *testing.T) {
	urls := []urlTestCase{}
	url := urlTestCase{}

	//url.url = "https://www.google.com/"
	//url.wantEncodedUrl = "www.google.com"
	//urls = append(urls, url)

	//	url.url = `https%3A%252F%252Fwww.google.com%252F`
	//	url.wantEncodedUrl = "www.google.com"
	//	urls = append(urls, url)

	url.url = `https%3A%2F%2Fwww.google.com%2F`
	url.wantEncodedUrl = "www.google.com"
	urls = append(urls, url)

	for _, url := range urls {
		t.Log(url)

		if gotEncodedUrl, err := CleanURL(url.url); err != nil {
			t.Fatal(err)
		} else if url.wantEncodedUrl != gotEncodedUrl {
			t.Error(url.url, ": ", "wanted ", url.wantEncodedUrl, " but got ", gotEncodedUrl)
		}

	}
}
