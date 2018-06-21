package component

import (
	"../model"
	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	netUrl "net/url"
	"strings"
	"sync"
	"time"
)

var db *gorm.DB
var client *http.Client

func init() {
	db = model.DB()
	db.LogMode(false)
	client = new(http.Client)
	client.Timeout = 3 * time.Second

	crawler()
	crawlerParseUrl()
}

func crawler() {
	var urls []model.Url
	var lastId uint
	var urlChan chan bool
	var lenChan int
	var err error
	maxChan := 3000
	limit := 10000
	where := "id > ?"
	wait := new(sync.WaitGroup)
	for {
		err = db.Debug().Preload("Source").Limit(limit).Find(&urls, where, lastId).Error
		if err != nil {
			log.Fatalln(err)
		}
		if len(urls) == 0 {
			break
		}
		lenChan = int(math.Max(float64(len(urls)), float64(maxChan)))
		urlChan = make(chan bool, lenChan)
		for _, url := range urls {
			urlChan <- true
			wait.Add(1)
			lastId = url.ID
			go crawlerUrl(&url, urlChan, wait)
		}
	}
	wait.Wait()
}

func crawlerUrl(url *model.Url, ch <-chan bool, wait *sync.WaitGroup) {
	defer func() {
		<-ch
		wait.Done()
	}()
	var count uint8
	db.Debug().Table("html").Where("url_id = ?", url.ID).Count(&count)
	if count > 0 { //已存在记录
		return
	}
	var err error
	var response *http.Response
	repeat := 5
	for {
		response, err = client.Get(url.Source.Home + url.Url)
		if err != nil {
			log.Println(err)
		}
		if response != nil && response.StatusCode == http.StatusOK {
			break
		}
		repeat--
		if repeat <= 0 {
			break
		}
	}
	if response == nil || response.StatusCode != http.StatusOK {
		return
	}
	bytes, _ := ioutil.ReadAll(response.Body)
	htmlString := string(bytes)
	html := model.Html{
		UrlId: url.ID,
		Data:  htmlString,
	}

	db.Debug().Save(&html)
	url.State = 0
	db.Debug().Save(url)
}

func crawlerParseUrl() {
	var htmlArr []model.Html
	var err error
	var lastId uint
	var htmlChan chan bool
	maxChan := 3000
	limit := 10000
	where := "id > ?"
	wait := new(sync.WaitGroup)
	for {
		err = db.Debug().Preload("Url").Limit(limit).Find(&htmlArr, where, lastId).Error
		if err != nil {
			log.Fatalln(err)
		}
		if len(htmlArr) == 0 {
			break
		}
		htmlChan = make(chan bool, int(math.Min(float64(len(htmlArr)), float64(maxChan))))
		for _, html := range htmlArr {
			htmlChan <- true
			wait.Add(1)
			go parseUrl(&html, htmlChan, wait)
		}
		lastId = htmlArr[len(htmlArr)-1].ID
	}
	wait.Wait()
}

func parseUrl(html *model.Html, ch <-chan bool, wait *sync.WaitGroup) {
	defer func() {
		<-ch
		wait.Done()
	}()
	var source model.Source
	var groups []model.HtmlSelectorGroup
	var selectors []model.HtmlSelector
	var lastId uint
	var err error
	var reader *strings.Reader
	var document *goquery.Document
	var selection *goquery.Selection
	base, err := netUrl.Parse(html.Url.Source.Home + html.Url.Url)
	if err != nil {
		log.Fatalln(err)
	}
	db.First(&source, html.Url.SourceId)
	if &source == nil {
		return
	}
	where := "id > ? AND table_name = ?"
	for {
		groups = groups[len(groups):]
		err = db.Limit(1000).Find(&groups, where, lastId, "url").Error
		if err != nil {
			log.Fatalln(err)
		}
		if len(groups) == 0 {
			break
		}
		lastId = groups[len(groups)-1].ID

		for _, group := range groups {
			selectors = selectors[len(selectors):]
			err = db.Order("sorting").Find(&selectors, "group_id = ?", group.ID).Error
			if err != nil {
				log.Fatalln(err)
			}
			if len(selectors) == 0 {
				continue
			}
			reader = strings.NewReader(html.Data)
			document, err = goquery.NewDocumentFromReader(reader)
			if err != nil {
				log.Println(err)
				continue
			}
			selection = document.Find("html")
			for _, selector := range selectors {
				selection = selection.Find(selector.Selector)
				if selector.Eq != -1 {
					selection = selection.Eq(int(selector.Eq))
				}
			}
			selection.Each(func(i int, selection *goquery.Selection) {
				href, _ := selection.Attr("href")
				u, err := netUrl.Parse(strings.TrimSpace(href))
				if err != nil {
					log.Fatalln(err)
				}
				if href != "" {
					newUrl := model.Url{
						SourceId: source.ID,
						Url:      base.ResolveReference(u).String(),
					}
					db.Debug().Save(&newUrl)
				}
			})
		}
	}
}
