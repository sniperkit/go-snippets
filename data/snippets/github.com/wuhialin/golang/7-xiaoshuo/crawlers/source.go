package crawlers

import (
	"../model"
	"io/ioutil"
	"time"
)

type source struct {
	crawler
	model model.Source
}

func (t *source) getUrl() string {
	return t.model.Url + t.model.MainAction
}

func (t *source) afterQuery() {
	bytes, _ := ioutil.ReadAll(t.response.Body)
	body := string(bytes)
	if body == "" {
		return
	}
	createAt := time.Now().Unix()
	urlModel := &model.Url{
		Url:      t.response.Request.URL.String(),
		SourceId: t.model.Id,
		State:    1,
		CreateAt: createAt,
	}
	model.GetOrm().InsertOrUpdate(urlModel, "url", "create_at = excluded.create_at")

	htmlModel := &model.Html{
		Data:     body,
		Url:      urlModel,
		CreateAt: createAt,
	}
	colArgs := "create_at = excluded.create_at, data = excluded.data"
	model.GetOrm().InsertOrUpdate(htmlModel, "url_id", colArgs)
}

func (t *source) start() {
	t.query(t.getUrl(), true)
	if t.response == nil {
		return
	}
	t.afterQuery()
}
