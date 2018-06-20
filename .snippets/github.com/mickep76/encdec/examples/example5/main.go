package main

import (
	"fmt"
	"log"

	"github.com/mickep76/encdec"
	_ "github.com/mickep76/encdec/yaml"
)

type Message struct {
	Name, Text string
}

type Messages []*Message

func main() {
	in := Messages{
		&Message{Name: "Ed", Text: "Knock knock."},
		&Message{Name: "Sam", Text: "Who's there?"},
		&Message{Name: "Ed", Text: "Go fmt."},
		&Message{Name: "Sam", Text: "Go fmt who?"},
		&Message{Name: "Ed", Text: "Go fmt yourself!"},
	}

	b, err := encdec.ToBytes("yaml", in)
	if err != nil {
		log.Fatal(err)
	}

	out := Messages{}
	if err := encdec.FromBytes("yaml", b, &out); err != nil {
		log.Fatal(err)
	}

	for _, m := range out {
		fmt.Printf("%s: %s\n", m.Name, m.Text)
	}
}
