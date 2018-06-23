package model

import (
	"../common"
	"time"
)

type Url struct {
	Id       int64
	Url      string
	SourceId int64 `mapstructure:"source_id"`
	State    int
	CreateAt int64 `mapstructure:"create_at"`
}

func (t *Url) Insert(params ...interface{}) (id int64, err error) {
	lenParam := len(params)
	fields := []string{"url", "source_id", "create_at"}
	t.CreateAt = time.Now().Unix()
	if lenParam > 0 {
		if lenParam < 3 {
			params = append(params, t.CreateAt)
		}
	} else {
		params = append(params, t.Url, t.SourceId, t.CreateAt)
	}
	id, err = common.Insert("url", fields, params...)
	if err == nil {
		t.Id = id
	}
	return
}
