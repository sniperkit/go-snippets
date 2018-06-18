package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Stat struct {
	Wins  int
	Loses int
}

func (s Stat) String() string {
	return fmt.Sprintf("Wins - %d, Loses - %d", s.Wins, s.Loses)
}

func round(a Shape, stat *Stat) {
	b := nextShape()
	var r rune
	if a == b {
		r = '='
	} else {
		if a.Beats(b) {
			r = '>'
			stat.Wins++
		} else {
			r = '<'
			stat.Loses++
		}
	}
	fmt.Printf("%s %c %s\n", a, r, b)
}

func main() {

	var stat Stat

	r := bufio.NewReader(os.Stdin)

	for {
		line, _, err := r.ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		command := strings.ToLower(string(line))

		switch command {
		case "q", "quit":
			return
		case "r", "rock":
			round(Rock, &stat)
		case "p", "paper":
			round(Paper, &stat)
		case "s", "scissors":
			round(Scissors, &stat)
		case "stat":
			fmt.Println(stat)
		default:
			fmt.Println("unknown command:", command)
		}
	}
}
