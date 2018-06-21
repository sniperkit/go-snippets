package model

type Parsed struct {
	Id                 int `orm:"auto"`
	HtmlId, SelectorId int
	CreateAt           int64
}
