package main

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	// db
	user string = "tire"
	pass string = "Raed1rei2a"
	dbname string = "tire"
	address string = "10.0.35.6"
	port int = 3306
	protocol string = "tcp"
	maxIdle int = 30
	maxConn int = 30
	o orm.Ormer
)

type Fbc_brand struct {
	Id				int			`orm:"pk;size(8);column(id);auto"`
	Key				string		`orm:"size(255)"`
	Value			string		`orm:"size(255)"`
}

func init() {
	// parseTime https://github.com/go-sql-driver/mysql#parsetime
	base := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8&parseTime=true", user, pass, protocol, address, port, dbname)
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", base, maxIdle, maxConn)
	// Timezone http://beego.me/docs/mvc/model/orm.md#timezone-config
	orm.DefaultTimeLoc, _ = time.LoadLocation("Asia/Novosibirsk")

	orm.RegisterModel(
	new(Fbc_brand),
	)
}

func main() {
	fmt.Println("mysql test script...\n")

	o = orm.NewOrm()
	o.Using("default")

	key := "LEXUS"
	table := "fbc_brand"
	exist := o.QueryTable(table).Filter("key", key).Exist()

	if(exist) {
		fmt.Printf("запись %s в таблице \"%s\" присутствует\n", key, table)
	} else {
		fmt.Printf("запись %s не найдена\n", key)
	}
}
