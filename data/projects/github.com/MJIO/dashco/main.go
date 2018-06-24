package main

import (
	"flag"
	"fmt"
	"os"
	"text/template"
)

const version = "0.0.1"

var helpMsg = `Usage: dashco URL

The following arguments can be used:
	--help		Show this help context.
	--version	Show the version number.
`

func showHelpMsg() {
	tmpl, err := template.New("").Parse(helpMsg)

	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(os.Stdout, nil); err != nil {
		panic(err)
	}

	os.Exit(2)
}

func showVersionMsg() {
	fmt.Printf("dashco %s\n", version)
	os.Exit(2)
}

func main() {
	var help, version bool

	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "help", false, "")
	flag.BoolVar(&version, "v", false, "")
	flag.BoolVar(&version, "version", false, "")
	flag.Parse()

	if help || len(os.Args[1:]) < 1 {
		showHelpMsg()
	}

	if version {
		showVersionMsg()
	}

	widgets, err := loadJsonData(os.Args[1])

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, widget := range widgets {
		switch widget.Type {
		case "line":
			fmt.Println(widget.Name, widget.Value)
		case "text":
			fmt.Println(widget.Name, widget.Value)
		default:
			fmt.Println("Default", widget.Type)
		}
	}
}
