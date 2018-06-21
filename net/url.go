package net

import (
	"fmt"
	"strings"
)

//
func Url2Part(url string) (srvId, uid string, ok bool) {
	ss := strings.Split(url, "#")
	if len(ss) != 2 {
		return
	}
	return ss[0], ss[1], true
}

//
func GenUrl(srvId, uid string) string {
	return fmt.Sprintf("%v#%v", srvId, uid)
}
