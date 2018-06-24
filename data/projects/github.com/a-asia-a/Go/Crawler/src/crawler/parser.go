package main

import (
	"log"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"sync"
)

func isSelected(url string) bool {
	return DefaultDomainSelector.Match([]byte(url))
}

func parseLinks(resp *http.Response, wg sync.WaitGroup) []string {
	if !isMIMEOk(resp) {
		return nil
	}
	var result []string
	tokenizer := html.NewTokenizer(resp.Body)

	checkAndAppend := func(uri string) {
		if idx := strings.LastIndex(uri, "#"); idx != -1 {
			uri = uri[:idx]
		}
		if idx := strings.Index(uri, "?"); idx != -1 {
			uri = uri[:idx]
		}
		if len(uri) < 4 {
			return
		}
		switch uri[len(uri) - 4:] {
		case
			".pdf", ".doc", ".zip", ".png", ".jpg", ".xls":
			return
		}
		uri = makeAbsolute(uri, resp.Request.URL)
		if len(uri) == 0 || !isSelected(uri) {
			return
		}
		uri = strings.TrimRight(uri, "/")
		result = append(result, uri)
	}

	loop:
	for {
		tt := tokenizer.Next()
		switch {
		case tt == html.ErrorToken:
			break loop
		case tt == html.StartTagToken || tt == html.SelfClosingTagToken:
			t := tokenizer.Token()
			if t.DataAtom.String() != "a" {
				continue
			}
			for _, attr := range t.Attr {
				if attr.Key == "href" {
					checkAndAppend(attr.Val)

				}
			}
		}
	}
	return result
}

func makeAbsolute(link string, req *url.URL) string {
	linkURI, err := url.Parse(link)
	if err != nil {
		log.Println(err)
		return ""
	}
	return req.ResolveReference(linkURI).String()
}

func isMIMEOk(resp *http.Response) bool {
	mediatype, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		log.Println(err)
		return false
	}

	switch mediatype {
	// only try to find links in HTML, or perhaps XML documents
	case
		"text/html",
		"application/atom+xml",
		"text/xml",
		"text/plain",
		"image/svg+xml":
		return true
	}

	return false
}
