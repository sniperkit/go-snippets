package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

var emailRegex *regexp.Regexp

func init() {
	emailRegex = regexp.MustCompile(`^((\"[0-9A-Za-z@.+-=]+\")|([0-9A-Za-z.+-=]+))@\w*(\w+\.)+\w{2,4}$`)
}

func email(q string) bool {
	return emailRegex.MatchString(q)
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(email(scanner.Text()))
	}
}
