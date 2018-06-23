package conf

import (
	"fmt"
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/zhangjunfang/softwareScaffolding/codeGenerator/vender/lib"
)

func Write() {
	f, err := os.OpenFile(lib.Src()+lib.HeadToUpper(lib.GetDBName())+"/controllers/default.go", os.O_RDWR, 0766)
	if err != nil {
		logs.Error(err.Error())
		os.Exit(-1)
	}
	f.WriteString(fmt.Sprintf(`
	package controllers

	import (
		"%s/utils"
		"github.com/astaxie/beego"
	)

	type MainController struct {
		beego.Controller
	}

	func (c *MainController) Get() {
		// Data
		var tpl utils.Tpl
		tpl.Title = "OceanCURD"
		tpl.Success = c.GetSession("success")
		tpl.Danger = c.GetSession("danger")
		c.DelSession("success")
		c.DelSession("danger")
		c.Data["Tpl"] = tpl
		c.Data["Website"] = "beego.me"
		c.Data["Email"] = "zhangjunfang0505@163.com"
		c.Layout = "public/tpl.html"
		c.TplName = "index.tpl"
	}
	`), lib.URL+lib.HeadToUpper(lib.GetDBName()))
	f.Close()
	lib.GoFmt(lib.Src() + lib.HeadToUpper(lib.GetDBName()) + "/controllers/default.go")
}
