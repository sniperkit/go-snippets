package main

import (
	"github.com/bblfsh/go-driver/driver/golang"
	"gopkg.in/bblfsh/sdk.v2/sdk/driver"
)

func main() {
	driver.NativeMain(golang.NewDriver())
}
