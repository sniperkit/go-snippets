package main

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
	"log"
)

type Html struct {
	Id       int `orm:"auto"`
	Data     string
	CreateAt int64
	Url      *Url `orm:"rel(fk)"`
}

type Url struct {
	Id              int `orm:"auto"`
	Url             string
	SourceId, State int
	CreateAt        int64
	Htmls           []*Html `orm:"reverse(many)"`
}

func main() {
	orm.Debug = true
	orm.RegisterDriver("postgres", orm.DRPostgres)
	source := "postgres://postgres:123456@localhost/postgres?sslmode=disable"
	maxIdle := 100
	maxConn := 1000
	alias := "default"
	driver := "postgres"

	orm.RegisterDataBase(alias, driver, source, maxIdle, maxConn)
	orm.RegisterModel(new(Html), new(Url))

	o := orm.NewOrm()
	var htmls []Html
	_, err := o.QueryTable(new(Html)).Filter("Url__source_id", 1).All(&htmls)
	if err != nil {
		log.Fatalln(err)
	}
	for _, html := range htmls {
		log.Println(html.Url)
	}
}
