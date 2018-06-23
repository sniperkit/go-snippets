package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/pravj/hackman/models"
	_ "github.com/pravj/hackman/routers"
)

func main() {
	beego.SessionOn = true
	beego.Run()
}

func init() {
	dbPath := beego.AppConfig.String("mysqluser") + ":" + beego.AppConfig.String("mysqlpass") + "@/" + beego.AppConfig.String("mysqldb")

	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", dbPath)

	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		beego.Info(err)
	}
}
