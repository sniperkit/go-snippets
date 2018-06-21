package crawler

import (
	"../../model"
	"github.com/gocolly/colly"
	"github.com/jinzhu/gorm"
	"log"
	"path"
	"strings"
)

var db *gorm.DB
var c *colly.Collector

func init() {
	db = model.DB()
	c = colly.NewCollector()
	c.CacheDir = path.Join("log", "colly")
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 20})
	onError()
	onHtmlUrl()
	onHtmlCustom()
	run()
}

func onError() {
	c.OnError(func(_ *colly.Response, e error) {
		log.Println(e)
	})
}

func onHtmlUrl() {
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := strings.TrimSpace(e.Attr("href"))
		if link != "" {
			c.Visit(e.Request.AbsoluteURL(e.Attr("href")))
		}
	})
}

func onHtmlCustom() {
	c.OnHTML("body", func(e *colly.HTMLElement) {
		bookItem(e)
		bookPage(e)
	})
}

func run() {
	var sources []model.Source
	db.Find(&sources)
	for _, source := range sources {
		c.AllowedDomains = append(c.AllowedDomains, source.Domain)
		c.Visit(source.Home)
	}
	c.Wait()
}
