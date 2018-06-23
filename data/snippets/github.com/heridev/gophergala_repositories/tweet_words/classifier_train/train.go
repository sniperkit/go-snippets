package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var input, positive, negative *os.File

func init() {
	var err error
	input, err = os.OpenFile("training.txt", os.O_RDONLY, 0777)
	if err != nil {
		panic(err)
	}
	positive, err = os.OpenFile("training-1.txt", os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	negative, err = os.OpenFile("training-0.txt", os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
}

func main() {
	defer input.Close()
	defer positive.Close()
	defer negative.Close()
	scanner := bufio.NewScanner(input)
	reader := bufio.NewReader(os.Stdin)
	var text string
	for scanner.Scan() {
		text = scanner.Text() + "\n"
		fmt.Print(text)
		fmt.Print("Choose (P)ositive/(N)egative/(Q)uit/(I)gnore [p]: ")
		choice, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		choice = strings.TrimSpace(strings.ToLower(choice)) + "p"
		switch choice[0] {
		case 'n':
			_, err := negative.WriteString(text)
			if err != nil {
				panic(err)
			}
		case 'q':
			return
		case 'i':
			continue
		default:
			_, err := positive.WriteString(text)
			if err != nil {
				panic(err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

}
