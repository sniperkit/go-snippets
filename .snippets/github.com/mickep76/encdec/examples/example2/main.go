package main

import (
	"bufio"
	"log"
	"os"

	"github.com/mickep76/encdec"
	_ "github.com/mickep76/encdec/json"
)

type Message struct {
	Name, Text string
}

type Messages []*Message

func main() {
	msgs := Messages{
		&Message{Name: "Ed", Text: "Knock knock."},
		&Message{Name: "Sam", Text: "Who's there?"},
		&Message{Name: "Ed", Text: "Go fmt."},
		&Message{Name: "Sam", Text: "Go fmt who?"},
		&Message{Name: "Ed", Text: "Go fmt yourself!"},
	}

	w := bufio.NewWriter(os.Stdout)
	enc := encdec.NewEncoder("json", w)
	for _, m := range msgs {
		if err := enc.Encode(m); err != nil {
			log.Fatal(err)
		}
	}
	if err := w.Flush(); err != nil {
		log.Fatal(err)
	}
}
