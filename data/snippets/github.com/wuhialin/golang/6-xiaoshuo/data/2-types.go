package data

import (
	"../common"
	"time"
)

func init() {
	table := "types"
	fields := []string{"name", "code", "create_at"}

	var params []interface{}

	params = append(params, "列表", common.TYPE_LIST, time.Now().Unix())
	params = append(params, "节点", common.TYPE_ITEM, time.Now().Unix())
	params = append(params, "翻页", common.TYPE_PAGE, time.Now().Unix())

	common.BatchInsert(table, fields, params...)
}
