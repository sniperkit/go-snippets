// Package parser parses text to identify special elements
package parser

import (
	"fmt"
	"golang.org/x/net/html"
	"mvdan.cc/xurls"
	"net/http"
	"regexp"
)

var mentionsRegex = regexp.MustCompile(`@(\w+)`)
var emoticonRegex = regexp.MustCompile(`\(([0-9A-Za-z]+)\)`)

const MAX_EMOTICON_SIZE = 15

type Link struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

func Mentions(message string) []string {
	found := mentionsRegex.FindAllStringSubmatch(message, -1)
	mentions := []string{}
	for _, item := range found {
		mentions = append(mentions, item[1])
	}
	return mentions
}

func Emoticons(message string) []string {
	found := emoticonRegex.FindAllStringSubmatch(message, -1)
	emoticons := []string{}
	for _, item := range found {
		if len(item[1]) <= MAX_EMOTICON_SIZE+2 {
			emoticons = append(emoticons, item[1])
		}
	}
	return emoticons
}

func Links(message string) ([]Link, error) {
	urls := xurls.Strict().FindAllString(message, -1)
	links := []Link{}
	for _, url := range urls {
		title, err := getTitle(url)
		if err != nil {
			return nil, err
		}
		link := Link{URL: url, Title: title}
		links = append(links, link)
	}
	return links, nil
}

// Bellow functions taken mostly verbatim from https://siongui.github.io/2016/05/10/go-get-html-title-via-net-html/
func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func traverse(n *html.Node) (string, error) {
	if isTitleElement(n) {
		return n.FirstChild.Data, nil
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, err := traverse(c)
		if err == nil {
			return result, nil
		}
	}

	return "", fmt.Errorf("never found <title> element")
}

func getTitle(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}
	return traverse(doc)
}
