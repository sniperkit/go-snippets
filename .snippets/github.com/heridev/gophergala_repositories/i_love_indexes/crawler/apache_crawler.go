package crawler

import (
	"regexp"
	"strings"

	"github.com/GopherGala/i_love_indexes/elasticsearch"
	"github.com/PuerkitoBio/goquery"
)

type ApacheCrawler struct {
	BaseCrawler
}

var (
	parentDirRegexp = regexp.MustCompile(`.*Parent Directory.*`)
)

func (crawler *ApacheCrawler) Crawl() error {
	errs := make(chan error)
	go func() {
		defer close(errs)
		defer crawler.End()

		doc := crawler.Doc

		headers := crawler.crawlHeaders()

		// Run through each row and extract data
		doc.Find("tr").Each(func(i int, s *goquery.Selection) {
			tds := s.Find("td")

			// Row is empty
			if text := strings.TrimSpace(tds.Text()); text == "" || parentDirRegexp.MatchString(text) {
				return
			}

			// Row has incorrect structure
			if tds.Size() != len(headers) {
				return
			}

			// Row has correct structure
			item := &elasticsearch.IndexItem{}
			tds.Each(func(i int, s *goquery.Selection) {
				// Ignore the img field
				if headers[i] == "img" {
					return
				}

				text := strings.TrimSpace(s.Text())
				if text == "" {
					return
				}

				if headers[i] == "Name" {
					item.Path, _ = s.Find("a").First().Attr("href")
				} else if headers[i] == "Size" {
					item.Size = ParseSize(text)
				} else if headers[i] == "LastModifiedAt" {
					date, err := ApacheParseDate(text)
					if err != nil {
						errs <- err
						return
					}
					item.LastModifiedAt = date
				}
			})
			crawler.itemsToIndex <- item
		})
	}()

	return <-errs
}

func (crawler *ApacheCrawler) crawlHeaders() []string {
	var fields []string
	// Find table header and determine available fields
	crawler.Doc.Find("tr th").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())

		if s.Children().Is("img") {
			fields = append(fields, "img")
		} else if text != "" {
			switch text {
			case "Last modified":
				text = "LastModifiedAt"
			}
			fields = append(fields, text)
		}
	})
	return fields
}
