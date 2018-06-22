package lib

import (
	"os"
	"strings"

	"github.com/astaxie/beego/logs"
)

// HeadToUpper ...
// @return string
// @desc 以下划线作切割对首字母大写
func HeadToUpper(oldstring string) string {
	arr := strings.Split(oldstring, "_")
	var newstring string
	for _, v := range arr {
		s := string(v)
		newstring += strings.Replace(s, string(s[0]), strings.ToUpper(string(s[0])), 1)
	}
	return newstring
}

// OutPerfix ...
// @return string
// @desc 以下划线作切割去掉前缀
func OutPerfix(oldstring string) string {
	arr := strings.Split(oldstring, "_")
	var newstring string
	if len(arr) == 2 {
		newstring = strings.ToLower(arr[1])
	} else if len(arr) == 1 {
		newstring = strings.ToLower(oldstring)
	} else {
		logs.Error("下划线数目大于1,无法处理")
		os.Exit(-1)
	}
	return newstring
}
