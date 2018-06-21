/*
* @Author: wuhailin
* @Date:   2018-01-25 10:16:49
* @Last Modified by:   wuhailin
* @Last Modified time: 2018-01-25 11:43:04
 */
package warehouse

import (
	"../../model"
	"fmt"
	"github.com/henrylee2cn/faygo"
	"net/http"
	"sync"
	"time"
)

type Index struct {
	Name string `param:"<in:query><desc:名称查询>"`
	Code string `param:"<in:query><desc:编码查询>"`
}

func (i *Index) Serve(ctx *faygo.Context) error {
	wait := new(sync.WaitGroup)
	start := time.Now()
	var warehouses []model.Warehouse
	var count uint
	db := model.DB().Model(new(model.Warehouse))
	if i.Name != "" {
		db = db.Where("name LIKE ?", i.Name+"%")
	}
	if i.Code != "" {
		db = db.Where("code = ?", i.Code)
	}
	wait.Add(2)
	go func() {
		defer wait.Done()
		db.Limit(10).Find(&warehouses)
	}()
	go func() {
		defer wait.Done()
		db.Count(&count)
	}()
	ctx.HeaderParam("HTTP_X_PJAX")
	wait.Wait()
	return ctx.Render(http.StatusOK, `tpl/warehouse/index.html`, faygo.Map{
		"runtime":    fmt.Sprint(time.Since(start)),
		"warehouses": warehouses,
		"count":      count,
		"this":       i,
	})
}
