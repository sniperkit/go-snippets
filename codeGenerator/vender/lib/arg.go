package lib

import (
	"os"
	"regexp"
	"strings"

	"github.com/astaxie/beego/logs"
)

// GetDBName ...
// @reutrn string
// @desc 验证无误返回
func GetDBName() string {
	// 参数个数判断
	if len(os.Args) >= 2 {
		logs.Error("参数输入不当")
		os.Exit(1)
	}
	// 参数正则验证
	r, _ := regexp.Compile(`-conn=([\w]+):([\S]+)@tcp\(([\d]+).([\d]+).([\d]+).([\d]+):([\d]+)\)\/([\w]+)$`)
	b := r.MatchString(os.Args[1])
	if !b {
		logs.Error("数据库URL输入不当")
		os.Exit(1)
	}
	return strings.Split(os.Args[1], "/")[1]
}

// GetDBUrl ...
// @reutrn string
// @desc 切去-conn=
func GetDBUrl() string {
	GetDBName()
	return strings.Split(os.Args[1], "=")[1]
}
