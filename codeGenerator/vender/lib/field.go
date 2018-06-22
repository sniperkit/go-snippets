package lib

import "github.com/astaxie/beego/orm"

type Field struct {
	Name    string `orm:"column(Field)"`
	Type    string `orm:"column(Type)"`
	Null    string `orm:"column(Null)"`
	Key     string `orm:"column(Key)"`
	Default string `orm:"column(Default)"`
	Extra   string `orm:"column(Extra)"`
}

// GetAllFields ...
// @return []slice
// @desc Gets all the field names in the table
func GetAllFields(tablename string) []Field {
	var fields []Field
	o := orm.NewOrm()
	o.Raw("desc " + tablename).QueryRows(&fields)
	return fields
}
