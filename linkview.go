package main

import (
	"fmt"
	"os"
)

func main() {
	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "linkview: %s\n", err)
		return
	}

	var file *os.File
	if (stat.Mode() & os.ModeNamedPipe) == 0 {
		if len(os.Args) < 2 {
			fmt.Fprintf(os.Stderr, "linkview: %s\n", "expected at least one argument")
			return
		}
		file, err = os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "linkview: %s\n", err)
			return
		}
	} else {
		file = os.Stdin
	}

	links, err := FindLinks(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "linkview: %s\n", err)
		return
	}

	terminal := NewTerminal(&links)

	err = terminal.Start()
	if err != nil {
		panic(err)
	}
	defer terminal.Close()

	terminal.Render()
	terminal.SetSize()
	terminal.Render()

	for {
		go PollEvent()
		quit, err := terminal.HandleEvent(<-EventChan)
		if err != nil {
			fmt.Fprintf(os.Stderr, "linkview: %s\n", err)
		}

		if quit {
			return
		}
	}
}
