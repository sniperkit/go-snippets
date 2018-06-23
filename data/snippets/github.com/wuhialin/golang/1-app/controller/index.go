package controller

import (
	"../common"
	"github.com/astaxie/beego/orm"
	"math"
	"net/http"
	"strconv"
)

type Home struct {
	http.Handler
}

func (t *Home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	o := common.RenderOptions()
	o.Layout = "layout"

	r.ParseForm()
	pageSize, _ := strconv.Atoi(r.Form.Get("pageSize"))
	page, _ := strconv.Atoi(r.Form.Get("page"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	query, _ := orm.NewQueryBuilder("mysql")
	query.From("community c").LeftJoin("branch b").On("c.branch_id = b.id")
	query.LeftJoin("region r").On("r.id = c.region_id")
	countQuery := "SELECT COUNT(*) " + query.String()
	query.Limit(10).Offset(page * pageSize)
	queryStr := "SELECT c.name, c.domain, b.name branch_name, r.name region_name " + query.String()
	row, err := common.DB.Query(countQuery)
	rows, err := common.DB.Query(queryStr)
	if err != nil {
		common.Render().Text(w, http.StatusInternalServerError, err.Error())
	}
	list := common.FetchAll(rows)
	count, _ := strconv.Atoi(common.FetchCol(row))
	maxPage := int(math.Ceil(float64(count / pageSize)))

	pages := map[int]int{}
	for i := page; i <= maxPage; i++ {
		pages[i] = i
	}
	data := map[string]interface{}{}
	data["pages"] = pages
	data["list"] = list
	common.RenderHTML(o).HTML(w, http.StatusOK, "index", data)
}
