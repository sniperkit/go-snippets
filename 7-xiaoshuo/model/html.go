package model

type Html struct {
	Id       int `orm:"auto"`
	Data     string
	CreateAt int64
	Url      *Url `orm:"rel(fk)"`
}
