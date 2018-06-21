package crawler

import (
	"../common"
	"../common/database"
	"../model"
	"github.com/goinggo/mapstructure"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type crawler struct {
	ch     <-chan int
	wait   *sync.WaitGroup
	repeat int
	client *http.Client
}

func (t crawler) start() {
	t.checkChannel()
	t.done()
}

func (t *crawler) checkChannel() {
	if t.ch == nil {
		panic("crawler.ch is not nil")
	}
	if t.wait == nil {
		panic("crawler.wait is not nil")
	}
}

func (t *crawler) initQuery() {
	if t.client == nil {
		t.client = new(http.Client)
		t.client.Timeout = common.HTTP_QUERY_TIMEOUT * time.Second
	}
}

func (t *crawler) query(url string, refresh bool) (response *http.Response) {
	t.initQuery()
	var err error
	repeat := t.repeat
	for {
		response, err = t.client.Get(url)
		if err == nil {
			return
		}
		repeat--
		if repeat <= 0 {
			break
		}
	}
	return
}

func (t *crawler) queryBody(url string, refresh bool) (body string) {
	r := t.query(url, refresh)
	if r != nil && r.StatusCode == http.StatusOK {
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		body = string(bytes)
	}
	return
}

func (t *crawler) done() {
	if t.ch != nil {
		<-t.ch
	}
	if t.wait != nil {
		t.wait.Done()
	}
}

func init() {
	go runSourceCrawler()
	go runTaskCrawler()
}

func runSourceCrawler() {
	var wait sync.WaitGroup
	s := database.Select{}
	s.From("source")
	rows, err := database.QueryAll(s)
	if err != nil {
		panic(err)
	}
	ch := make(chan int, 10)
	for k, row := range rows {
		var dataType model.Source
		if err = mapstructure.Decode(row, &dataType); err != nil {
			println(err)
			continue
		}
		ch <- k
		obj := source{}
		obj.ch = ch
		obj.wait = &wait
		wait.Add(1)
		go obj.start(dataType)
	}
	wait.Wait()
}

func runTaskCrawler() {
	for range time.Tick(time.Second) {
		//println(time.Now().Unix())
	}
}
