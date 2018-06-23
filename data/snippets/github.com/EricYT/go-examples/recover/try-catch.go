package main

import (
	"fmt"
	"reflect"
)

func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()

	fun()
}

func main() {
	Try(func() {
		panic("crash main")
	}, func(err interface{}) {
		fmt.Printf("main error:%+v type:%+v\n", err, reflect.TypeOf(err))
		err_, _ := err.(string)
		fmt.Printf("error:%s\n", err_)
	})
}
