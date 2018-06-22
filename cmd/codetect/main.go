package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/vmarkovtsev/CodeNeuron.v1"
)

func main() {
	textBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Reading from stdin: %v", err)
	}
	if len(textBytes) == 0 {
		log.Fatal("empty input")
	}
	text := string(textBytes)
	session, err := codetect.OpenSession()
	if err != nil {
		log.Fatalf("Failed to create a Tensorflow session: %v", err)
	}
	defer session.Close()
	boundaries, err := codetect.Run(text, session)
	if len(boundaries) == 0 {
		fmt.Print(text)
		return
	}
	pos := 0
	for i, c := range text {
		fmt.Print(string(c))
		if pos < len(boundaries) && boundaries[pos].PositionInRunes == i {
			if boundaries[pos].Start {
				fmt.Print("<code>")
			} else {
				fmt.Print("</code>")
			}
			pos++
		}
	}
}
