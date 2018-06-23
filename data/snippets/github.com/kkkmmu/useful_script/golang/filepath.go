package main

import (
	"fmt"
	"path/filepath"
)

var Url = "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage.jpg"

func main() {
	fmt.Println(filepath.Base(Url))
}
