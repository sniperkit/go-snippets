package model

import (
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/zhangjunfang/softwareScaffolding/codeGenerator/vender/lib"
)

func Create() {
	for _, v := range lib.GetAllTables() {
		s := string(v.TableName)
		f, err := os.Create(lib.Src() + lib.HeadToUpper(lib.GetDBName()) + "/models/" + lib.HeadToUpper(s) + ".go")
		defer f.Close()
		if err != nil {
			logs.Error(err.Error())
		}
	}
}
