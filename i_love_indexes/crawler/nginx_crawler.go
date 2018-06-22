package crawler

import (
	"regexp"
	"strings"

	"github.com/GopherGala/i_love_indexes/elasticsearch"
)

type NginxCrawler struct {
	BaseCrawler
}

var (
	spaceRegexp = regexp.MustCompile(`\s\s+`)
)

func (crawler *NginxCrawler) Crawl() error {
	errs := make(chan error)
	go func() {
		defer close(errs)
		defer crawler.End()
		doc := crawler.Doc

		// Run through each row and extract data
		pre := doc.Find("pre").First()
		as := pre.Find("a").Nodes
		if len(as) == 0 {
			return
		}
		entries := strings.Split(pre.Text(), "\n")
		var err error
		for i, entry := range entries {
			entry = strings.TrimSpace(entry)
			if len(entry) == 0 {
				continue
			}
			item := &elasticsearch.IndexItem{}
			item.Path = as[i].Attr[0].Val
			if strings.Contains(entry, "../") {
				continue
			}
			fields := strings.Split(strings.TrimSpace(spaceRegexp.ReplaceAllString(entry[51:], "\t")), "\t")
			item.LastModifiedAt, err = ApacheParseDate(fields[0])
			if err != nil {
				errs <- err
			}
			item.Size = mustInt64(fields[1])
			crawler.itemsToIndex <- item
		}
	}()
	return <-errs
}
