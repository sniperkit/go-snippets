package main

import (
	"bytes"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/guotie/gogb2312"
	"log"
	//"github.com/guotie/gogb2312"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

//<a target="_blank" href="http://news.sina.com.cn/c/nd/2017-06-11/doc-ifyfzfyz3230662.shtml">东航悉尼飞上海航班发动机损伤 已返航人机安全</a>
//var newre = regexp.MustCompile(`\<a[[:word:][:space:]_\"\'=]*href=\"+(?P<link>[[:word:]:\.\/-_\?$%^&]+)\"`)
//var newre = regexp.MustCompile(`\<a[[:word:][:space:]_\"\'=]*href=\"+(?P<link>[[:word:]:\.\/-_\?$%^&]+)\"[[:space:]]*\>(?P<head>[\p{Han}[:space:][:word:]#%&@]+)\<`)

var newre = regexp.MustCompile(`\<a[[:word:][:space:]_\"\'=]*href=\"+(?P<link>[[:word:]:\.\/-_\?$%^&]+)\"[[:space:]]*\>(?P<head>.*)<`)

type News struct {
	Link string
	Head string
}

var Site = map[string]string{
	"http://www.163.com":                         "div.yaowen_news>div>ul>li",
	"http://sports.163.com/yc":                   "div.topnews>h2&&&div.topnews>ul>li",
	"http://www.sina.com.cn":                     "ul.news_top>li&&&div.top_newslist>ul>li",
	"http://sports.sina.com.cn/g/premierleague/": "ul.match_news_list>li",
	"https://soccer.hupu.com/england/":           "ul.england-bignews>li>h2&&&ul.england-bignews>li>h4",
}

func GetNewsFromSite(site, rule string) []*News {
	sp := NewSpider("News")

	result := make([]*News, 0, 1)
	rules := strings.Split(rule, "&&&")

	var raw []string
	for _, r := range rules {
		log.Println(site, r)
		news, err := sp.SpideHTML(site, r)
		if err != nil {
			log.Panic(err.Error())
		}
		raw = append(raw, news...)
	}

	for _, s := range raw {
		s := strings.Replace(s, " target=\"_blank\"", "", -1)
		log.Println(s)
		match := newre.FindStringSubmatch(s)
		if len(match) == 0 {
			continue
		}
		result = append(result, &News{Link: string(match[1]), Head: string(match[2])})
	}

	return result
}

func main() {
	news := make([]*News, 0, 1)
	for s, r := range Site {
		news = append(news, GetNewsFromSite(s, r)...)
	}

	for _, n := range news {
		fmt.Println(n.Head, n.Link)
	}
}

type Spider struct {
	name   string
	client *http.Client
}

func NewSpider(name string) *Spider {
	return &Spider{
		name: name,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (s *Spider) SpideContent(page, rule string) ([]string, error) {

	resp, err := s.client.Get(page)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ctype := resp.Header.Get("Content-Type")
	if strings.Contains(strings.ToLower(ctype), "gbk") {
		body, err, _, _ = gogb2312.ConvertGB2312(body)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	reader := bytes.NewReader(body)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("Cannot create Docuement by Response")
	}

	var res = make([]string, 0)
	doc.Find(rule).Each(func(ix int, sl *goquery.Selection) {
		res = append(res, sl.Text())
	})
	return res, nil
}

func (s *Spider) SpideAttribute(page, rule, attr string) ([]string, error) {
	resp, err := s.client.Get(page)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ctype := resp.Header.Get("Content-Type")
	if strings.Contains(strings.ToLower(ctype), "gbk") {
		body, err, _, _ = gogb2312.ConvertGB2312(body)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	reader := bytes.NewReader(body)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("Cannot create Docuement by Response")
	}

	var res = make([]string, 0)
	doc.Find(rule).Each(func(ix int, sl *goquery.Selection) {
		attr, ok := sl.Attr(attr)
		if ok {
			res = append(res, attr)
		}
	})
	return res, nil
}

func (s *Spider) SpideHTML(page, rule string) ([]string, error) {
	resp, err := s.client.Get(page)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ctype := resp.Header.Get("Content-Type")
	if strings.Contains(strings.ToLower(ctype), "gbk") {
		body, err, _, _ = gogb2312.ConvertGB2312(body)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	reader := bytes.NewReader(body)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("Cannot create Docuement by Response")
	}

	var res = make([]string, 0)
	doc.Find(rule).Each(func(ix int, sl *goquery.Selection) {
		content, _ := sl.Html()
		res = append(res, content)
	})
	return res, nil
}
