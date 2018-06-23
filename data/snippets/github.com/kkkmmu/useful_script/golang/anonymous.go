package main

import (
	"log"
)

type Formatter interface {
	Lower(string) string
	Upper(string) string
}

type Html struct {
	document string
	Formatter
}

type Text struct {
	content string
	Formatter
}

func format(f Formatter) {
	log.Println(f.Lower)
	log.Println(f.Lower("Hello"))
	log.Println(f.Upper("Hello"))
}

func main() {
	h := &Html{document: "<p>Hello</p>"}
	t := &Text{content: "Hello world"}

	format(h)
	format(t)
}
