package crawlers

import (
	"../model"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/orm"
	"math"
	"strings"
	"sync"
	"time"
)

var o orm.Ormer

func init() {
	o = model.GetOrm()

	go startSource()
	go startQuery()
}

//运行源表数据搜索
func startSource() {
	var sources []model.Source
	ch := make(chan interface{}, 20)
	o.QueryTable(&model.Source{}).Filter("state", 1).All(&sources)

	for _, m := range sources {
		ch <- m.Id
		c := &source{
			model: m,
		}
		c.setCh(ch)
		go c.start()

		ch <- 1
		go startChanParseHtml(m, ch)
	}
}

func startChanParseHtml(m model.Source, ch <-chan interface{}) {
	defer func() {
		<-ch
	}()
	var htmls []model.Html
	_, err := o.QueryTable(new(model.Html)).Filter("Url__source_id", m.Id).All(&htmls)
	if err != nil {
		panic(err)
	}

	parseChan := make(chan bool, int(math.Min(100, float64(len(htmls)))))
	var wait sync.WaitGroup
	for _, html := range htmls {
		wait.Add(1)
		go startParseHtml(m, html, parseChan, &wait)
	}
	wait.Wait()
}

//解析页面数据分析进程
func startParseHtml(m model.Source, html model.Html, ch <-chan bool, wait *sync.WaitGroup) {
	defer func() {
		<-ch
		wait.Done()
	}()
	var tmpSelectors []model.Selector
	selectorsMap := map[int][]model.Selector{}
	selectorQs := o.QueryTable(new(model.Selector)).Filter("source_id", m.Id).Filter("state", model.SwitchOn)
	_, err := selectorQs.OrderBy("type_id", "sorting").All(&tmpSelectors)
	if err != nil {
		panic(err)
	}
	for _, selector := range tmpSelectors {
		selectorsMap[selector.TypeId] = append(selectorsMap[selector.TypeId], selector)
	}

	var r *strings.Reader
	var dom *goquery.Document
	var selection *goquery.Selection
	var types model.Types
	r = strings.NewReader(html.Data)
	for _, selectors := range selectorsMap {
		dom, err = goquery.NewDocumentFromReader(r)
		if err != nil {
			panic(err)
		}
		selection = dom.Find("body")
		for _, selector := range selectors {
			selection = selection.Find(selector.Selector)
			if selector.Eq != -1 {
				selection = selection.Eq(selector.Eq)
			}
		}
		if selection.Length() > 0 {
			o.QueryTable(new(model.Types)).Filter("id", selectors[0].TypeId).One(&types)
			selection.Each(func(i int, selection *goquery.Selection) {
				text := strings.TrimSpace(selection.Text())
				href, _ := selection.Attr("href")
				var params []interface{}
				params = append(params, m, types, text, href)
				if types.TableName != "" {
					model.Call(types.TableName, params...)
				}
				model.Call("url", params...)
				parsed := model.Parsed{
					HtmlId:     html.Id,
					SelectorId: selectors[len(selectors)-1].Id,
					CreateAt:   time.Now().Unix(),
				}
				_, err = model.GetOrm().InsertOrUpdate(&parsed, "html_id, select_id", "create_at = excluded.create_at")
				if err != nil {
					panic(err)
				}
				model.Log().Debug(selection.Text())
			})
		}
	}
}

//循环获取URL的内容
func startQuery() {
	var urls []model.Url
	chanUrls := make(chan model.Url, 10)
	defer close(chanUrls)
	qs := o.QueryTable(&model.Url{}).Filter("state", 0).Limit(1000)
	for range time.Tick(time.Second) {
		urls = urls[len(urls):]
		qs.All(&urls)
		for _, m := range urls {
			chanUrls <- m
			c := new(url)
			c.model = m
			go c.start(chanUrls)
		}
	}
}
