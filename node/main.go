package main

import (
	"fmt"
	"node/application"
)

func main() {

	application.GetAppInstance().Run()
	fmt.Println("程序时从主协程这里退出的哈")
}
