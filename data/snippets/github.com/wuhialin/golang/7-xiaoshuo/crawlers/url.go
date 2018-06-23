package crawlers

import (
	"../model"
	"io/ioutil"
	"time"
)

type url struct {
	crawler
	model model.Url
}

func (t *url) getUrl() string {
	return t.model.Url
}

func (t *url) afterQuery() {
	bytes, _ := ioutil.ReadAll(t.response.Body)
	body := string(bytes)
	if body == "" {
		return
	}
	createAt := time.Now().Unix()
	t.model.State = 1
	_, err := model.GetOrm().Update(&t.model)
	if err != nil {
		panic(err)
	}

	htmlModel := &model.Html{
		Data:     body,
		Url:      &t.model,
		CreateAt: createAt,
	}
	colArgs := "create_at = excluded.create_at, data = excluded.data"
	model.GetOrm().InsertOrUpdate(htmlModel, "url_id", colArgs)
}

func (t *url) start(chanUrl <-chan model.Url) {
	defer func() {
		<-chanUrl
	}()
	t.query(t.getUrl(), true)
	if t.response == nil {
		return
	}
	t.afterQuery()
}
