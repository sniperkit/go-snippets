package model

type Source struct {
	Id                            int `orm:"auto"`
	CreateAt                      int64
	Name, Domain, Url, MainAction string
	State                         int
}
