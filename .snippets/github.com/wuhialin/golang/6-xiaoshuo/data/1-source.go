package data

import (
	"../common"
	"log"
	"time"
)

func init() {
	fields := []string{"name", "domain", "url", "main_action", "create_at"}
	var params []interface{}

	params = append(params, "61读书网", "m.61xsw.com", "http://m.61xsw.com", "/list.html", time.Now().Unix())

	_, err := common.BatchInsert("source", fields, params...)
	if err != nil {
		log.Fatalln(err)
	}
}
