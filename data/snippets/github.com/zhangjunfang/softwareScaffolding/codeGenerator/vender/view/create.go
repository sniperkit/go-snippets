package view

import (
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/zhangjunfang/softwareScaffolding/codeGenerator/vender/lib"
)

func Create() {
	appname := lib.Src() + lib.HeadToUpper(lib.GetDBName())
	for _, v := range lib.GetAllTables() {
		s := string(v.TableName)
		path := appname + "/views/" + lib.OutPerfix(s)
		err := os.MkdirAll(path, 0766)
		if err != nil {
			logs.Error(err.Error())
			os.Exit(-1)
		}
		f, err := os.Create(path + "/index.html")
		f.Close()
		f, err = os.Create(path + "/create.html")
		f.Close()
		f, err = os.Create(path + "/update.html")
		f.Close()
		f, err = os.Create(path + "/view.html")
		f.Close()
		f, err = os.Create(path + "/_form.html")
		f.Close()
		if err != nil {
			logs.Error(err.Error())
			os.Exit(-1)
		}
	}
	err := os.Mkdir(appname+"/views/public", 0766)
	f, err := os.Create(appname + "/views/public/tpl.html")
	f.Close()
	if err != nil {
		logs.Error(err.Error())
		os.Exit(-1)
	}
}
