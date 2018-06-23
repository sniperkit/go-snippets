package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

const tpl = `Snippets of Go code for Sublime Text
====================================

{{ range . }}
{{ .TabTrigger }}:
` + "```go" + `
{{ .Content }}
` + "```" + `
{{ end }}
`

type snippet struct {
	XMLName xml.Name `xml:"snippet"`

	Content    string `xml:"content"`
	TabTrigger string `xml:"tabTrigger"`
}

func main() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	var snippets []*snippet
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sublime-snippet") {
			snippets = append(snippets, extractCode(f.Name()))
		}
	}

	w, err := os.Create("README.md")
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err := template.New("doc").Parse(tpl)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, snippets)
	if err != nil {
		log.Fatal(err)
	}

}

func extractCode(fn string) *snippet {
	r, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	var data *snippet
	if err := xml.NewDecoder(r).Decode(&data); err != nil {
		log.Fatal(err)
	}

	return data
}
