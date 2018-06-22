package router

import (
	"fmt"
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/zhangjunfang/softwareScaffolding/codeGenerator/vender/lib"
)

func Write() {
	f, err := os.OpenFile(lib.Src()+lib.HeadToUpper(lib.GetDBName())+"/routers/router.go", os.O_RDWR, 0766)
	if err != nil {
		logs.Error(err.Error())
	}
	f.WriteString(Package() + Import() + Init())
	f.Close()
	lib.GoFmt(lib.Src() + lib.HeadToUpper(lib.GetDBName()) + "/routers/router.go")
}

func Package() string {
	return "package routers\n"
}

func Import() string {
	return fmt.Sprintf(`
		import (
			"%s/controllers"
			"github.com/astaxie/beego"
		)
	`, lib.URL+lib.HeadToUpper(lib.GetDBName()))
}

func Init() string {
	var str string
	for _, v := range lib.GetAllTables() {
		u := lib.HeadToUpper(v.TableName)
		l := lib.OutPerfix(v.TableName)
		str += fmt.Sprintf(`
		beego.Router("/%s", &controllers.%sController{})
		beego.Router("/%s/create", &controllers.%sController{}, "*:Create")
		beego.Router("/%s/update", &controllers.%sController{}, "*:Update")
		beego.Router("/%s/view", &controllers.%sController{}, "*:View")
		beego.Router("/%s/delete", &controllers.%sController{}, "*:Del")
		`, l, u, l, u, l, u, l, u, l, u)
	}
	return fmt.Sprintf(`
	func init() {
		beego.Router("/", &controllers.MainController{})
		%s
		}
	`, str)
}
