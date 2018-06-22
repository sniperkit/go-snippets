package main

import (
	"golang.org/x/mobile/app"

	_ "golang.org/x/mobile/bind/java"
	_ "golang.org/x/mobile/example/libhello/hi/go_hi"
)

func main() {
	app.Run(app.Callbacks{})
}

