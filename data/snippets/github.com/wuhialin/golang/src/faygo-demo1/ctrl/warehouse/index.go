package warehouse

import (
	"faygo-demo1/ctrl"
	"faygo-demo1/model"
	"github.com/henrylee2cn/faygo"
	"net/http"
	"sync"
	"fmt"
	"strings"
)

type Index struct {
	Name string `param:"<in:query>"`
	Code string `param:"<in:query>"`

	ctrl.Controller
}

func (t *Index) Serve(ctx *faygo.Context) error {
	m := model.Warehouse{}
	db := m.SearchList(ctx.QueryParamAll())
	wait := new(sync.WaitGroup)
	wait.Add(2)
	go func() {
		db.Count(&t.Total)
		wait.Done()
	}()
	var warehouses []model.Warehouse
	go func() {
		db.Limit(t.GetPageSize()).Offset((t.GetPage() - 1) * t.GetPageSize()).Find(&warehouses)
		wait.Done()
	}()
	wait.Wait()
	tempName := "index"
	if strings.TrimSpace(ctx.R.Header.Get("X-PJAX")) == "true" {
		tempName = "_list"
	}
	ctx.Render(http.StatusOK, fmt.Sprintf("tpl/warehouse/%s.html", tempName), faygo.Map{
		"ctrl":       t,
		"warehouses": warehouses,
		"title":"仓库列表页",
	})
	return nil
}
