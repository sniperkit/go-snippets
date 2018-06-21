package user

import (
	"github.com/henrylee2cn/faygo"
	"faygo-demo2/models"
	"net/http"
	"fmt"
)

type Index struct {
	Page  uint `param:"<in:query>"`
	Limit uint `param:"<in:query>"`
}

func (t *Index) Serve(ctx *faygo.Context) (err error) {
	var users []models.Employee
	m := new(models.Employee)
	db := m.DB()
	var count uint
	db.Model(m).Count(&count)
	ctx.W.Header().Add("x-total-count", fmt.Sprintf("%d", count))
	if t.Limit == 0 {
		t.Limit = 100
	}
	db = db.Limit(t.Limit)
	if t.Page > 1 {
		db = db.Offset(t.Limit * (t.Page - 1))
	}
	db.Find(&users)
	ctx.JSON(http.StatusOK, users, true)
	return
}
