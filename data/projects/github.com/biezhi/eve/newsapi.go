package main

import (
	"github.com/biezhi/moe"
	"encoding/json"
	"log"
	"github.com/olekukonko/tablewriter"
	"os"
	"reflect"
	"sync"
	"net/http"
	"io/ioutil"
)

var (
	wg   sync.WaitGroup
	lock sync.Mutex
)

// BaseArticle article
type BaseArticle struct {
	Title   string `json:"title,omitempty"`
	URL     string `json:"url,omitempty"`
	Content string `json:"content,omitempty"`
}

// ShowNews show news
func ShowNews(q string) {
	tip := "loading all news..."
	url := "http://reader.one/api/all/hn,reddit,ph,dn,github,medium,lifehacker?limit=20"
	if q != "" {
		tip = "loading " + q + " ..."
		url = "http://reader.one/api/news/" + q + "?limit=20"
	}

	moe := moe.New(tip).Color(moe.Green).Start()
	body := GetRequestBody(url)

	var resp []BaseArticle
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Fatalln("load news fail")
		return
	}

	items := mapToString(resp, "URL")
	shortUrls := shortURLItems(items)

	moe.Stop()

	render(resp, shortUrls)
}

// render articles to console
func render(articles []BaseArticle, shortUrls map[string]*string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetColWidth(70)

	table.SetHeader([]string{"Title", "URL"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)

	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgBlackColor})

	table.SetColumnColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgRedColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlueColor})

	for _, article := range articles {
		URL := shortUrls[article.URL]
		Title := article.Title
		row := []string{Title, *URL}
		table.Append(row)
	}
	table.Render()
}

// GetShortURL get sina short url
func GetShortURL(url string) string {
	reqURL := "http://api.t.sina.com.cn/short_url/shorten.json?source=3271760578&url_long=" + url
	if url == "" {
		return ""
	}
	var urlArr []map[string]interface{}
	body := GetRequestBody(reqURL)
	if err := json.Unmarshal(body, &urlArr); err != nil {
		return ""
	}
	return urlArr[0]["url_short"].(string)
}

// GetRequestBody http get body
func GetRequestBody(reqURL string) []byte {
	resp, err := http.Get(reqURL)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

func mapToString(struts []BaseArticle, fieldName string) []string {
	result := make([]string, len(struts))
	for _, item := range struts {
		v := reflect.ValueOf(item)
		f := v.FieldByName(fieldName)
		result = append(result, f.String())
	}
	return result
}

func shortURLItems(longUrls []string) map[string]*string {
	shortUrls := make(map[string]*string, len(longUrls))
	for _, url := range longUrls {
		wg.Add(1)
		go func(url string) {
			lock.Lock()
			defer func() {
				lock.Unlock()
				wg.Done()
			}()
			URL := GetShortURL(url)
			shortUrls[url] = &URL
		}(url)
	}
	wg.Wait()
	return shortUrls
}
