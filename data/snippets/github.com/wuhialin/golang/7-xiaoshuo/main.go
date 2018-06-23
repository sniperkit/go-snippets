package main

import (
	_ "./crawlers"
	_ "./model"
	"github.com/astaxie/beego/orm"
	"time"
)

func main() {
	orm.Debug = true
	for range time.Tick(time.Second) {
	}
}
