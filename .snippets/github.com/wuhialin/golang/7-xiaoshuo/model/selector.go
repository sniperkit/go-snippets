package model

type Selector struct {
	Id                 int `orm:"auto"`
	SourceId, TypeId   int
	Eq                 int `orm:"default(-1)"`
	Sorting            int `orm:"default(0)"`
	Pid                int `orm:"default(0)"`
	State              int `orm:"default(1)"`
	Selector           string
	CreateAt, UpdateAt int64
}
