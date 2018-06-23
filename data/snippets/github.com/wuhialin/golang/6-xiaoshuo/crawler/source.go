package crawler

import (
	m "../model"
)

type source struct {
	crawler
	model *m.Source
}

func (t *source) afterQueryBody(body string, url string) {
	if body != "" {
		u := &m.Url{
			Url:      url,
			SourceId: t.model.Id,
		}
		id, err := u.Insert()
		if err != nil {
			return
		}
		h := &m.Html{
			Data:  body,
			UrlId: id,
		}
		h.Insert()
	}
}

func (t *source) queryBody(url string, refresh bool) (body string) {
	body = t.crawler.queryBody(url, refresh)
	t.afterQueryBody(body, url)
	return
}

func (t *source) start(m m.Source) {
	t.model = &m
	t.queryBody(m.Url+m.MainAction, false)
	t.crawler.start()
}
