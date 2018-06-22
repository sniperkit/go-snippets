package crawler

import (
	"strings"

	"github.com/GopherGala/i_love_indexes/elasticsearch"
	"github.com/PuerkitoBio/goquery"
)

type LighttpdCrawler struct {
	BaseCrawler
}

func (crawler *LighttpdCrawler) Crawl() error {
	errs := make(chan error)
	go func() {
		defer close(errs)
		defer crawler.End()
		doc := crawler.Doc

		// Run through each row and extract data
		doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
			tds := s.Find("td")

			item := &elasticsearch.IndexItem{}
			tds.Each(func(i int, s *goquery.Selection) {
				text := strings.TrimSpace(s.Text())
				if text == "" {
					return
				}

				class, _ := s.Attr("class")
				var err error
				switch class {
				case "n":
					item.Path, _ = s.Find("a").First().Attr("href")
					if item.Path == "../" {
						return
					}
				case "m":
					item.LastModifiedAt, err = LighttpdParseDate(text)
					if err != nil {
						errs <- err
					}
				case "s":
					item.Size = ParseSize(text)
				case "t":
					item.MIMEType = text
				}
			})
			crawler.itemsToIndex <- item
		})
	}()

	return <-errs
}
