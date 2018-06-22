package crawler

import (
	"log"

	"github.com/GopherGala/i_love_indexes/elasticsearch"
	"github.com/Scalingo/go-workers"
	"gopkg.in/errgo.v1"
)

func CrawlWorker(msg *workers.Msg) {
	params, err := msg.Args().Array()
	if err != nil {
		log.Println("param is not a string:", err)
		return
	}
	id := params[0].(string)
	path := params[1].(string)

	indexOf, err := elasticsearch.FindIndexOf(id)
	if err != nil {
		log.Println("index of not found", id)
		return
	}

	crawler, err := NewCrawler(indexOf, path)
	if err != nil {
		log.Println("failed to get crawler", errgo.Details(err))
		return
	}
	err = crawler.Crawl()
	if err != nil {
		log.Println("failed to crawl", errgo.Details(err))
		return
	}
}
