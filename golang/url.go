package main

import (
	"fmt"
	"net/url"
)

var Url = "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage?ticket=xxx&uuid=xxx&lang=xxx&scan=xxx"

func main() {
	u, err := url.Parse(Url)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", u)
}
