package main

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

	"github.com/wzulfikar/lab/go/graphqlboiler/templates"
)

func main() {
	data := struct {
		TypeName, ResolverName, RootResolver string
	}{"Person", "personResolver", "RootResolver"}

	tpl, err := parseTpl(templates.Resolver, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tpl)
}

func parseTpl(tplString string, data interface{}) (string, error) {
	t, err := template.New("").Parse(tplString)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
