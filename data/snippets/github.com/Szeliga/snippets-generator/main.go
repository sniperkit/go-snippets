package main

import (
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Snippets is a container for a list of snippets
type Snippets struct {
	Snippets []Snippet
}

// Snippet is a snippet
type Snippet struct {
	Name        string
	Description string
	Acronym     string
	Body        string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func split(input string) []string {
	return strings.Split(input, "\n")
}

func join(input []string) string {
	return strings.Join(input, "\n      ")
}

func main() {
	atomTemplate, err := ioutil.ReadFile("./atom_template.cson")
	check(err)
	s := Snippets{}
	data, err := ioutil.ReadFile("./snippets.yml")
	check(err)
	err = yaml.Unmarshal([]byte(data), &s)
	check(err)

	funcMap := template.FuncMap{
		"split": split,
		"join":  join,
		"trim":  strings.TrimSpace,
	}

	tmpl, err := template.New("text").Funcs(funcMap).Parse(string([]byte(atomTemplate)))
	check(err)
	err = tmpl.Execute(os.Stdout, s)
	check(err)
}
