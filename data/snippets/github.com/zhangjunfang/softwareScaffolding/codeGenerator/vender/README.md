# Copy the code and run it
 
```
package main
 
import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zhangjunfang/softwareScaffolding/codeGenerator/vender/conf"
	"github.com/zhangjunfang/softwareScaffolding/codeGenerator/vender/controller"
	"github.com/zhangjunfang/softwareScaffolding/codeGenerator/vender/lib"
	"github.com/zhangjunfang/softwareScaffolding/codeGenerator/vender/model"
	"github.com/zhangjunfang/softwareScaffolding/codeGenerator/vender/router"
	"github.com/zhangjunfang/softwareScaffolding/codeGenerator/vender/utils"
	"github.com/zhangjunfang/softwareScaffolding/codeGenerator/vender/view"
)

func init() {
	orm.RegisterDataBase("default", "mysql", lib.GetDBUrl())
}

func main() {
	if lib.CreateApp() {
		router.Write()
		conf.AddSessionOn()
		conf.Newmain()
		conf.Write()
		controller.Create()
		controller.Write()
		model.Create()
		model.Write()
		view.Create()
		view.Write()
		utils.Create()
		utils.Write()
	} else {
		logs.Error("Project creation fails")
	}
}
```
##  Code run and parameter format
go run github.com/zhangjunfang/softwareScaffolding/crud.go -conn="root:20170725@tcp(127.0.0.1:3306)/onlinebbs"  -URL="github.com"
1.parameter specification：
  1.1 conn：Golang links to the URL of the mysql database
  1.2 URL :  git repository url 