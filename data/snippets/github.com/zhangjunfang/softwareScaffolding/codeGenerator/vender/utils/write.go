package utils

import (
	"fmt"
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/zhangjunfang/softwareScaffolding/codeGenerator/vender/lib"
)

func Write() {
	f1, err := os.OpenFile(lib.Src()+lib.HeadToUpper(lib.GetDBName())+"/utils/log.go", os.O_RDWR, 0766)
	if err != nil {
		logs.Error(err.Error())
	}
	_, err = f1.WriteString(fmt.Sprintf(`
	package utils

	import (
		"fmt"

		"github.com/astaxie/beego/logs"
	)

	/**
	*@param status string ["Success","Error","Warn","Info"]
	*@param msg
	*@param level []
	*@SaveFile /logs/[status].log
	*/

	func ToLog(status string, msg string, level int) {
		logs.SetLogger(logs.AdapterFile, fmt.Sprintf("{\"filename\":\"logs/%s.log\",\"level\":%d}", status, level))
		switch status {
		case "Error":
			logs.Error(msg)
		case "Success":
			logs.Notice(msg)
		case "Warn":
			logs.Warning(msg)
		case "Info":
			logs.Informational(msg)
		default:
			logs.Debug(msg)
		}
	}
	`))
	f1.Close()
	lib.GoFmt(lib.Src() + lib.HeadToUpper(lib.GetDBName()) + "/utils/tpl.go")

	f2, err := os.OpenFile(lib.Src()+lib.HeadToUpper(lib.GetDBName())+"/utils/tpl.go", os.O_RDWR, 0766)
	if err != nil {
		logs.Error(err.Error())
	}
	_, err = f2.WriteString(fmt.Sprintf(`
	package utils

	type Tpl struct {
		Attributes map[string]string
		Action     string
		Danger     interface{}
		Model      interface{}
		Success    interface{}
		Title      string
	}
	`))
	f2.Close()
	lib.GoFmt(lib.Src() + lib.HeadToUpper(lib.GetDBName()) + "/utils/tpl.go")
}
