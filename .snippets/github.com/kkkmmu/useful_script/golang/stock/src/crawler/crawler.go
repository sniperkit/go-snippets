package crawler

import (
	"bytes"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/guotie/gogb2312"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func Crawl(page, rule string) ([]string, error) {
	var (
		res = make([]string, 0) //for leaf
		wg  sync.WaitGroup
		mu  sync.Mutex
	)

	resp, err := http.Get(page)
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
	//Need convert the encode to UTF-8
	body, err, _, _ = gogb2312.ConvertGB2312(body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	reader := bytes.NewReader(body)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("Cannot create Docuement by Response")
	}

	doc.Find(rule).Each(func(ix int, sl *goquery.Selection) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			res = append(res, sl.Text())
			mu.Unlock()
		}()
	})
	wg.Wait()
	return res, nil

}
