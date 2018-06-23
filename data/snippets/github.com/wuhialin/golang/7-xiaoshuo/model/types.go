package model

type Types struct {
	Id                    int `orm:"auto"`
	Name, Code, TableName string
	CreateAt              int64
}
