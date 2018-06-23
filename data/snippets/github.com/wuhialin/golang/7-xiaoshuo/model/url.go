package model

type Url struct {
	Id              int `orm:"auto"`
	Url             string
	SourceId, State int
	CreateAt        int64
	Html            []*Html `orm:"reverse(many)"`
}
