package crawlers

import (
	"../model"
	"net/http"
	"time"
)

type Icrawler interface {
	start()
	done()
	getUrl() string
	setCh(<-chan interface{})
	getResponse() *http.Response
}

type crawler struct {
	model    interface{}
	repeat   int
	response *http.Response
	ch       <-chan interface{}
	client   *http.Client
}

func (t *crawler) getResponse() *http.Response {
	return t.response
}

func (t *crawler) setCh(ch <-chan interface{}) {
	t.ch = ch
}

func (t crawler) start() {
	if t.model == nil {
		model.Log().Error("model is not nil")
		return
	}
}

func (t *crawler) done() {
	<-t.ch
}

func (t *crawler) initClient() {
	if t.client == nil {
		t.client = new(http.Client)
		t.client.Timeout = 3 * time.Second
	}
}

func (t *crawler) afterQuery(body string) {
}

func (t *crawler) getUrl() string {
	return ""
}

func (t *crawler) query(url string, flag bool) {
	if flag {
		htmlExists := model.GetOrm().QueryTable(new(model.Html)).Filter("Url__url", url).Exist()
		if htmlExists {
			return
		}
	}

	t.initClient()
	var response *http.Response
	var err error
	for {
		startTime := time.Now()
		model.Log().Trace("query %s start", url)
		response, err = t.client.Get(url)
		model.Log().Trace("query %s end, time:%s", url, time.Since(startTime))
		if err == nil && response.StatusCode == http.StatusOK {
			t.response = response
			return
		}
		if t.repeat != -1 && t.repeat <= 0 {
			break
		}
		t.repeat--
	}
}
