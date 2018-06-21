package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
)

func main() {
	document, err := goquery.NewDocument("http://m.61xsw.com/list.html")
	if err != nil {
		log.Fatalln(err)
	}
	document.Find("div.book-all-list div.bd a.name").Each(func(i int, s *goquery.Selection) {
		log.Println(s.Text())
	})
}
