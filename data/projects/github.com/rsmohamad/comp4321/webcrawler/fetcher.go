package webcrawler

import (
	"github.com/rsmohamad/comp4321/models"
	"github.com/rsmohamad/comp4321/stopword"
	"net/url"
	"regexp"
	"strings"
	"time"

	"fmt"
	"github.com/surgebase/porter2"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"mvdan.cc/xurls"
)

// Get url from html token
func getUrl(t html.Token) string {
	// Iterate over all attributes
	for _, a := range t.Attr {
		if a.Key == "href" {
			return a.Val
		}
	}
	return ""
}

func handleMetaRedirect(t html.Token) (url string, redirect bool) {
	redirect = false
	for _, a := range t.Attr {
		if a.Key == "http-equiv" && a.Val == "refresh" {
			redirect = true
		}
		if a.Key == "content" && redirect {
			url = xurls.Strict().FindString(a.Val)
		}
	}
	return
}

// Fix relative URL to absolute URL
func toAbsoluteUrl(links []string, base string) (rv []string) {
	baseUrl, _ := url.Parse(base)
	for _, link := range links {
		uri, err := url.Parse(link)

		if err != nil {
			continue
		}

		if !uri.IsAbs() {
			uri = baseUrl.ResolveReference(uri)
		}

		// Skip PDFs
		if strings.HasSuffix(strings.ToLower(uri.Path), ".pdf") {
			continue
		}

		if strings.HasPrefix(uri.String(), "http") {
			newUrl := "http://" + uri.Host + uri.Path
			rv = append(rv, newUrl)
		}
	}
	return
}

// Clean and tokenize string
func tokenizeString(s string) (rv []string) {
	regex := regexp.MustCompile("[^a-zA-Z0-9 ]")
	s = regex.ReplaceAllString(s, " ")
	regex = regexp.MustCompile("[^\\s]+")
	words := regex.FindAllString(s, -1)
	for _, word := range words {
		cleaned := strings.ToLower(word)
		cleaned = strings.TrimSpace(cleaned)
		cleaned = porter2.Stem(cleaned)
		if !stopword.IsStopWord(cleaned) {
			rv = append(rv, cleaned)
		}
	}
	return
}

func Fetch(uri string) (page *models.Document) {
	words := make([]string, 0)
	var lastElement string
	page = &models.Document{Uri: uri}
	inBody := false

	// Make HTTP GET request
	res, err := fetchClient.Get(uri)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Return if HTTP request is not successful
	if res == nil || res.StatusCode != 200 {
		return nil
	}

	if res.Header.Get("Content-Type") != "text/html" {
		return nil
	}

	tm, _ := time.Parse(time.RFC1123, res.Header.Get("Last-Modified"))
	page.Modtime = tm.Unix()
	if page.Modtime < 0 {
		tm, _ = time.Parse(time.RFC1123, res.Header.Get("Date"))
		page.Modtime = tm.Unix()
	}
	page.Len = 0

	// Tokenize
	tokenizer := html.NewTokenizer(res.Body)
	defer res.Body.Close()

	// Loop through all html elements
	for {
		tokenType := tokenizer.Next()
		t := tokenizer.Token()
		page.Len += len(tokenizer.Raw())

		if tokenType == html.ErrorToken {
			break
		}

		switch tokenType {
		case html.StartTagToken:
			// Indicate when inside body tags
			if t.Data == "body" {
				inBody = true
			}
			// Title
			if t.Data == "title" {
				tokenizer.Next()
				page.Title = strings.TrimSpace(tokenizer.Token().Data)
			}
			// Links
			if t.Data == "a" && inBody {
				page.Links = append(page.Links, getUrl(t))
			}
			lastElement = t.Data
			break
		case html.SelfClosingTagToken:
			if t.Data == "meta" {
				link, redirect := handleMetaRedirect(t)
				if redirect {
					return Fetch(link)
				}
				break
			}
			break
		case html.TextToken:
			// Skip if text is empty, not in between body tags or between script tags
			trimmed := strings.TrimSpace(t.Data)
			if trimmed != "" && inBody && lastElement != "script" && lastElement != "style" {
				words = append(words, trimmed)
			}
		}
	}

	// Makes sure all bytes are read
	io.Copy(ioutil.Discard, res.Body)

	// Clean data
	page.Titles = models.CountTfandIdx(tokenizeString(page.Title))
	page.Words = models.CountTfandIdx(tokenizeString(strings.Join(words, " ")))
	page.MaxTf = models.CountMaxTf(page.Words)
	page.TitleMaxTf = models.CountMaxTf(page.Titles)
	page.Links = toAbsoluteUrl(page.Links, uri)
	return
}
