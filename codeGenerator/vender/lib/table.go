package lib

import "github.com/astaxie/beego/orm"
import "fmt"

type Table struct {
	TableName string `orm:"column(table_name)"`
}

// GetAllTables ...
// @return []slice
// @desc Gets all the table names in the database
func GetAllTables() []Table {
	var tablenames []Table
	o := orm.NewOrm()
	o.Raw(fmt.Sprintf("select distinct table_name from information_schema.columns where table_schema='%s'", GetDBName())).QueryRows(&tablenames)
	return tablenames
}
