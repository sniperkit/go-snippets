package model

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
	"log"
	"time"
)

var o orm.Ormer

func init() {
	orm.RegisterDriver("postgres", orm.DRPostgres)
	source := "postgres://postgres:123456@localhost/postgres?sslmode=disable"
	maxIdle := 100
	maxConn := 1000
	alias := "default"
	driver := "postgres"

	orm.RegisterDataBase(alias, driver, source, maxIdle, maxConn)
	orm.RegisterModel(new(Source), new(Types), new(Tag), new(Task))
	orm.RegisterModel(new(Selector), new(Url), new(Html), new(Parsed))
	orm.RegisterModel(new(Category), new(Author), new(Book), new(BookTag))
	orm.RegisterModel(new(BookCategory), new(BookPages))

	o = orm.NewOrm()

	//initBootstrapData()
}

func initBootstrapData() {
	//初始化数据
	var err error
	var selector []Selector

	createAt := time.Now().Unix()

	source := &Source{
		CreateAt:   createAt,
		Name:       "61小说网",
		Domain:     "m.61xsw.com",
		Url:        "http://m.61xsw.com",
		MainAction: "/list.html",
		State:      SwitchOn,
	}
	_, err = o.InsertOrUpdate(source, "domain", "create_at = excluded.create_at")
	if err != nil {
		log.Fatalln(err)
	}

	types := &Types{
		Name:     "列表",
		Code:     TYPE_LIST,
		CreateAt: createAt,
	}
	_, err = o.InsertOrUpdate(types, "code", "create_at = excluded.create_at")
	if err != nil {
		log.Fatalln(err)
	}

	selector = selector[len(selector):]
	selector = append(selector, Selector{
		SourceId: source.Id,
		TypeId:   types.Id,
		Selector: "div.book-all-list div.bd",
		Eq:       0,
		Sorting:  0,
		State:    SwitchOn,
		CreateAt: time.Now().Unix(),
		UpdateAt: time.Now().Unix(),
	})
	selector = append(selector, Selector{
		SourceId: source.Id,
		TypeId:   types.Id,
		Selector: "a.name",
		Eq:       -1,
		Sorting:  1,
		State:    SwitchOn,
		CreateAt: time.Now().Unix(),
		UpdateAt: time.Now().Unix(),
	})
	_, err = o.InsertMulti(len(selector), &selector)

	types = &Types{
		Name:     "翻页",
		Code:     TYPE_PAGE,
		CreateAt: createAt,
	}
	_, err = o.InsertOrUpdate(types, "code", "create_at = excluded.create_at")
	if err != nil {
		log.Fatalln(err)
	}
	selector = selector[len(selector):]
	selector = append(selector, Selector{
		SourceId: source.Id,
		TypeId:   types.Id,
		Selector: "div.page",
		Eq:       -1,
		Sorting:  0,
		State:    SwitchOn,
		CreateAt: time.Now().Unix(),
		UpdateAt: time.Now().Unix(),
	})
	selector = append(selector, Selector{
		SourceId: source.Id,
		TypeId:   types.Id,
		Selector: "a",
		Eq:       -1,
		Sorting:  1,
		State:    SwitchOn,
		CreateAt: time.Now().Unix(),
		UpdateAt: time.Now().Unix(),
	})
	_, err = o.InsertMulti(len(selector), &selector)
	if err != nil {
		log.Fatalln(err)
	}

	types = &Types{
		Name:     "节点",
		Code:     TYPE_ITEM,
		CreateAt: createAt,
	}
	_, err = o.InsertOrUpdate(types, "code", "create_at = excluded.create_at")
	if err != nil {
		log.Fatalln(err)
	}

	types = &Types{
		Name:      "分类",
		Code:      TypeCategory,
		TableName: TypeCategory,
		CreateAt:  createAt,
	}
	_, err = o.InsertOrUpdate(types, "code", "create_at = excluded.create_at")
	if err != nil {
		log.Fatalln(err)
	}

	selector = selector[len(selector):]
	selector = append(selector, Selector{
		SourceId: source.Id,
		TypeId:   types.Id,
		Selector: "div.book-all-list div.hd div.filter",
		Eq:       0,
		Sorting:  1,
		CreateAt: createAt,
		State:    SwitchOn,
		UpdateAt: createAt,
	})
	selector = append(selector, Selector{
		SourceId: source.Id,
		TypeId:   types.Id,
		Selector: "a",
		Eq:       -1,
		Sorting:  2,
		State:    SwitchOn,
		CreateAt: createAt,
		UpdateAt: createAt,
	})
	_, err = o.InsertMulti(len(selector), &selector)
}

func GetOrm() orm.Ormer {
	return o
}
