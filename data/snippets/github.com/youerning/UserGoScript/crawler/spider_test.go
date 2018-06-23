package videos

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"log"
	"strings"
)

func main() {
	c := colly.NewCollector()
	testUrl := "http://www.80s.tw/movie/22021"
	c.SetDebugger(&debug.LogDebugger{})
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.84 Safari/537.36"

	c.OnHTML("#block3", func(e *colly.HTMLElement) {
		imgSrc := e.ChildAttr("img", "src")
		title := e.ChildText("h1")
		// desc := e.ChildText("div.info span:first-child")
		// esc := e.DOM.Find("div.info span").Get(1)
		desc := childOfText(e.DOM, 1, "div.info span")
		alias := childOfText(e.DOM, 2, "div.info span")
		alias = strings.TrimSpace(strings.Split(alias, "：")[1])
		// content := e.ChildText("div#movie_content")

		videoSrc := e.ChildAttr("input[name='list[]']", "value")
		thunderLink := e.ChildAttr("span.dlname span a", "href")
		fmt.Println("Detail in url: ", e.Request.URL)
		fmt.Println("image url: ", imgSrc)
		fmt.Printf("Title: %s\nDesc: %s\nAlias: %s\nThunderLink: %q\nvideoSrc: %q\n",
			title, desc, alias, thunderLink, videoSrc)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL: ", r.Request.URL, "faild with response: ", r.Body, "\nERROR: ", err)

	})

	c.Visit(testUrl)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func childAttr(ele *goquery.Selection, goquerySelector, attrName string) string {
	if attr, ok := ele.Find(goquerySelector).Attr(attrName); ok {
		return strings.TrimSpace(attr)
	}
	return ""
}

func childText(ele *goquery.Selection, goquerySelector string) string {
	return strings.TrimSpace(ele.Find(goquerySelector).Text())
}

func childOfText(ele *goquery.Selection, index int, goquerySelector string) string {
	ret := ""
	ele.Find(goquerySelector).Each(func(i int, s *goquery.Selection) {
		if i == index {
			ret = strings.TrimSpace(s.Text())
		}
	})
	return ret
}
