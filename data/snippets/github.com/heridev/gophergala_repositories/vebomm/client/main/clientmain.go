package main

import (
	"fmt"

	"github.com/phaikawl/wade/app"
	"github.com/phaikawl/wade/platform/clientside"

	"github.com/gophergala/vebomm/client"
)

func main() {
	ap := app.New(app.Config{
		BasePath: "/web",
	}, clientside.CreateBackend())

	err := ap.Start(client.AppMain{Application: ap})
	if err != nil {
		panic(fmt.Sprintf("Failed to load, error: %v.", err.Error()))
	}
}
