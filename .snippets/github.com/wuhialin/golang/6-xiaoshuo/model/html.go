package model

import (
	"../common"
	"database/sql"
	"time"
)

type Html struct {
	Id       int64
	Data     string
	UrlId    int64 `mapstructure:"url_id"`
	CreateAt int64 `mapstructure:"create_at"`
}

func (t *Html) Insert(params ...interface{}) (sql.Result, error) {
	lenParam := len(params)
	fields := []string{"data", "url_id", "create_at"}
	t.CreateAt = time.Now().Unix()
	if lenParam > 0 {
		if lenParam < 3 {
			params = append(params, t.CreateAt)
		}
	} else {
		params = append(params, t.Data, t.UrlId, t.CreateAt)
	}
	return common.BatchInsert("html", fields, params...)
}
