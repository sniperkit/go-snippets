package main

import (
	"os"
	"strings"
	"text/template"
)

func main() {
	// first we create a FuncMap to register the function.
	funcMap := template.FuncMap{
		// The name "title" is what the function will be called in the template text
		"title": strings.Title,
	}

	const templateText = `
Input: {{printf "%q" .}}
Output 0: {{title .}}
Output 1: {{title . | printf "%q"}}
Output 2: {{printf "%q" . | title}}
`

	// create a template, add the funcMap, and parse the text
	tmpl, err := template.New("titletest").Funcs(funcMap).Parse(templateText)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, "the go programming language")
	if err != nil {
		panic(err)
	}
}
