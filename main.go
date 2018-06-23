package main

import (
	"github.com/LeadPipeSoftware/medkit/cmd/medkit"
)

var version string
var date string
var commit string

func main() {
	medkit.Version = version
	medkit.Date = date
	medkit.Commit = commit

	medkit.Execute()
}
