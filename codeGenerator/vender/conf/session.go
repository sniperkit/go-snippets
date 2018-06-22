package conf

import (
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/zhangjunfang/softwareScaffolding/codeGenerator/vender/lib"
)

// AddSessionOn ...
// @desc Open project session function
func AddSessionOn() {
	f, err := os.OpenFile(lib.Src()+lib.HeadToUpper(lib.GetDBName())+"/conf/app.conf", os.O_WRONLY, 0766)
	if err != nil {
		logs.Error(err.Error())
		os.Exit(-1)
	}
	// Find the offset at the end of the file
	n, _ := f.Seek(0, os.SEEK_END)
	//Start writing the content from the offset at the end
	_, err = f.WriteAt([]byte("sessionOn = true"), n)
}
