package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func init() {
	var err error
	source := "host=localhost user=postgres dbname=postgres sslmode=disable password=123456"
	db, err = gorm.Open("postgres", source)
	if err != nil {
		panic(err)
	}

	db.SingularTable(true)

	db.DB().SetMaxOpenConns(1000)
	db.DB().SetMaxIdleConns(10)

	initAutoMigrate()
	//initAutoData()
}

func initAutoMigrate() {
	db.AutoMigrate(&Source{}, &Url{}, &HtmlSelector{}, &Html{}, &HtmlSelectorGroup{})
	db.AutoMigrate(&HtmlParse{})
}

func initAutoData() {
	source := Source{
		Name:   "61读书网",
		Domain: "61xsw.com",
		Home:   "http://m.61xsw.com",
	}
	db.Save(&source)

	sourceId := uint(source.ID)
	url := Url{
		Url:      "/list.html",
		SourceId: sourceId,
	}
	db.Save(&url)

	group := HtmlSelectorGroup{
		Name:      "列表页",
		SourceId:  sourceId,
		TableName: "url",
	}
	db.Save(&group)

	selector := HtmlSelector{
		Selector: "div.book-all-list div.bd a.name",
		Eq:       -1,
		SourceId: sourceId,
		GroupId:  group.ID,
	}
	db.Save(&selector)

	group = HtmlSelectorGroup{
		Name:      "列表翻页",
		SourceId:  sourceId,
		TableName: "url",
	}
	db.Save(&group)
	selector = HtmlSelector{
		Selector: "div.page a",
		Eq:       -1,
		SourceId: sourceId,
		GroupId:  group.ID,
	}
	db.Save(&selector)

	group = HtmlSelectorGroup{
		Name:      "开始阅读",
		SourceId:  sourceId,
		TableName: "url",
	}
	db.Save(&group)
	selector = HtmlSelector{
		Selector: "div.detail a.read.start",
		Eq:       -1,
		SourceId: sourceId,
		GroupId:  group.ID,
	}
	db.Save(&selector)

	group = HtmlSelectorGroup{
		Name:      "章节",
		SourceId:  sourceId,
		TableName: "url",
	}
	db.Save(&group)
	selector = HtmlSelector{
		Selector: "#chapterlist a",
		Eq:       -1,
		SourceId: sourceId,
		GroupId:  group.ID,
	}
	db.Save(&selector)
}

func DB() *gorm.DB {
	return db.Scopes(stateWhere)
}

type Model struct {
	gorm.Model
	State uint8 `gorm:"not null;default:1"`
}

func stateWhere(db *gorm.DB) *gorm.DB {
	return db.Where("state = ?", 1)
}
